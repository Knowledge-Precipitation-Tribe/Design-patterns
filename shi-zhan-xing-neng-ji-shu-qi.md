# 实战-性能计数器

## 版本1

我们之前实现过一个[初代版本的性能计数器](she-ji-yuan-ze/ru-he-kai-fa-fei-ye-wu-xi-tong.md)，在版本 1 中，整个框架的代码被划分为下面这几个类。

* MetricsCollector：负责打点采集原始数据，包括记录每次接口请求的响应时间和请求时间戳，并调用 MetricsStorage 提供的接口来存储这些原始数据。
* MetricsStorage 和 RedisMetricsStorage：负责原始数据的存储和读取。
* Aggregator：是一个工具类，负责各种统计数据的计算，比如响应时间的最大值、最小值、平均值、百分位值、接口访问次数、tps。
* ConsoleReporter 和 EmailReporter：相当于一个上帝类（God Class），定时根据给定的时间区间，从数据库中取出数据，借助 Aggregator 类完成统计工作，并将统计结果输出到相应的终端，比如命令行、邮件。

MetricCollector、MetricsStorage、RedisMetricsStorage 的设计与实现比较简单，不是版本 2 重构的重点。今天，我们重点来看一下 Aggregator 和 ConsoleReporter、EmailReporter 这几个类。

## Aggregator 类存在的问题

Aggregator 类里面只有一个静态函数，有 50 行左右的代码量，负责各种统计数据的计算。当要添加新的统计功能的时候，我们需要修改 aggregate\(\) 函数代码。一旦越来越多的统计功能添加进来之后，这个函数的代码量会持续增加，可读性、可维护性就变差了。因此，我们需要在版本 2 中对其进行重构。

```java
public class Aggregator {
  public static RequestStat aggregate(List<RequestInfo> requestInfos, long durationInMillis) {
    double maxRespTime = Double.MIN_VALUE;
    double minRespTime = Double.MAX_VALUE;
    double avgRespTime = -1;
    double p999RespTime = -1;
    double p99RespTime = -1;
    double sumRespTime = 0;
    long count = 0;
    for (RequestInfo requestInfo : requestInfos) {
      ++count;
      double respTime = requestInfo.getResponseTime();
      if (maxRespTime < respTime) {
        maxRespTime = respTime;
      }
      if (minRespTime > respTime) {
        minRespTime = respTime;
      }
      sumRespTime += respTime;
    }
    if (count != 0) {
      avgRespTime = sumRespTime / count;
    }
    long tps = (long)(count / durationInMillis * 1000);
    Collections.sort(requestInfos, new Comparator<RequestInfo>() {
      @Override
      public int compare(RequestInfo o1, RequestInfo o2) {
        double diff = o1.getResponseTime() - o2.getResponseTime();
        if (diff < 0.0) {
          return -1;
        } else if (diff > 0.0) {
          return 1;
        } else {
          return 0;
        }
      }
    });
 
    if (count != 0) {
      int idx999 = (int)(count * 0.999);
      int idx99 = (int)(count * 0.99);
      p999RespTime = requestInfos.get(idx999).getResponseTime();
      p99RespTime = requestInfos.get(idx99).getResponseTime();
    }
    RequestStat requestStat = new RequestStat();
    requestStat.setMaxResponseTime(maxRespTime);
    requestStat.setMinResponseTime(minRespTime);
    requestStat.setAvgResponseTime(avgRespTime);
    requestStat.setP999ResponseTime(p999RespTime);
    requestStat.setP99ResponseTime(p99RespTime);
    requestStat.setCount(count);
    requestStat.setTps(tps);
    return requestStat;
  }
}

public class RequestStat {
  private double maxResponseTime;
  private double minResponseTime;
  private double avgResponseTime;
  private double p999ResponseTime;
  private double p99ResponseTime;
  private long count;
  private long tps;
  //...省略getter/setter方法...
}
```

## ConsoleReporter 和 EmailReporter 这两个类存在的问题

ConsoleReporter 和 EmailReporter 两个类中存在代码重复问题。在这两个类中，从数据库中取数据、做统计的逻辑都是相同的，可以抽取出来复用，否则就违反了 DRY 原则。

整个类负责的事情比较多，不相干的逻辑糅合在里面，职责不够单一。特别是显示部分的代码可能会比较复杂（比如 Email 的显示方式），最好能将这部分显示逻辑剥离出来，设计成一个独立的类。

除此之外，因为代码中涉及线程操作，并且调用了 Aggregator 的静态函数，所以代码的可测试性也有待提高。

```java
public class ConsoleReporter {
  private MetricsStorage metricsStorage;
  private ScheduledExecutorService executor;

  public ConsoleReporter(MetricsStorage metricsStorage) {
    this.metricsStorage = metricsStorage;
    this.executor = Executors.newSingleThreadScheduledExecutor();
  }

  public void startRepeatedReport(long periodInSeconds, long durationInSeconds) {
    executor.scheduleAtFixedRate(new Runnable() {
      @Override
      public void run() {
        long durationInMillis = durationInSeconds * 1000;
        long endTimeInMillis = System.currentTimeMillis();
        long startTimeInMillis = endTimeInMillis - durationInMillis;
        Map<String, List<RequestInfo>> requestInfos =
                metricsStorage.getRequestInfos(startTimeInMillis, endTimeInMillis);
        Map<String, RequestStat> stats = new HashMap<>();
        for (Map.Entry<String, List<RequestInfo>> entry : requestInfos.entrySet()) {
          String apiName = entry.getKey();
          List<RequestInfo> requestInfosPerApi = entry.getValue();
          RequestStat requestStat = Aggregator.aggregate(requestInfosPerApi, durationInMillis);
          stats.put(apiName, requestStat);
        }
        System.out.println("Time Span: [" + startTimeInMillis + ", " + endTimeInMillis + "]");
        Gson gson = new Gson();
        System.out.println(gson.toJson(stats));
      }
    }, 0, periodInSeconds, TimeUnit.SECONDS);
  }

}

public class EmailReporter {
  private static final Long DAY_HOURS_IN_SECONDS = 86400L;

  private MetricsStorage metricsStorage;
  private EmailSender emailSender;
  private List<String> toAddresses = new ArrayList<>();

  public EmailReporter(MetricsStorage metricsStorage) {
    this(metricsStorage, new EmailSender(/*省略参数*/));
  }

  public EmailReporter(MetricsStorage metricsStorage, EmailSender emailSender) {
    this.metricsStorage = metricsStorage;
    this.emailSender = emailSender;
  }

  public void addToAddress(String address) {
    toAddresses.add(address);
  }

  public void startDailyReport() {
    Calendar calendar = Calendar.getInstance();
    calendar.add(Calendar.DATE, 1);
    calendar.set(Calendar.HOUR_OF_DAY, 0);
    calendar.set(Calendar.MINUTE, 0);
    calendar.set(Calendar.SECOND, 0);
    calendar.set(Calendar.MILLISECOND, 0);
    Date firstTime = calendar.getTime(); 
    Timer timer = new Timer();
    timer.schedule(new TimerTask() {
      @Override
      public void run() {
        long durationInMillis = DAY_HOURS_IN_SECONDS * 1000;
        long endTimeInMillis = System.currentTimeMillis();
        long startTimeInMillis = endTimeInMillis - durationInMillis;
        Map<String, List<RequestInfo>> requestInfos =
                metricsStorage.getRequestInfos(startTimeInMillis, endTimeInMillis);
        Map<String, RequestStat> stats = new HashMap<>();
        for (Map.Entry<String, List<RequestInfo>> entry : requestInfos.entrySet()) {
          String apiName = entry.getKey();
          List<RequestInfo> requestInfosPerApi = entry.getValue();
          RequestStat requestStat = Aggregator.aggregate(requestInfosPerApi, durationInMillis);
          stats.put(apiName, requestStat);
        }
        // TODO: 格式化为html格式，并且发送邮件
      }
    }, firstTime, DAY_HOURS_IN_SECONDS * 1000);
  }

}
```

## 版本1问题重构

Aggregator 类和 ConsoleReporter、EmailReporter 类主要负责统计显示的工作。之前我们提到，如果我们把统计显示所要完成的功能逻辑细分一下，主要包含下面 4 点：

* 根据给定的时间区间，从数据库中拉取数据；
* 根据原始数据，计算得到统计数据；
* 将统计数据显示到终端（命令行或邮件）；
* 定时触发以上三个过程的执行；

之前的划分方法是将所有的逻辑都放到 ConsoleReporter 和 EmailReporter 这两个上帝类中，而 Aggregator 只是一个包含静态方法的工具类。这样的划分方法存在前面提到的一些问题，我们需要对其进行重新划分。



