# 实战-ID生成器

## ID 生成器需求背景介绍

“ID”中文翻译为“标识（Identifier）”。这个概念在生活、工作中随处可见，比如身份证、商品条形码、二维码、车牌号、驾照号。聚焦到软件开发中，ID 常用来表示一些业务信息的唯一标识，比如订单的单号或者数据库中的唯一主键，比如地址表中的 ID 字段（实际上是没有业务含义的，对用户来说是透明的，不需要关注）。

假设你正在参与一个后端业务系统的开发，为了方便在请求出错时排查问题，我们在编写代码的时候会在关键路径上打印日志。某个请求出错之后，我们希望能搜索出这个请求对应的所有日志，以此来查找问题的原因。而实际情况是，在日志文件中，不同请求的日志会交织在一起。如果没有东西来标识哪些日志属于同一个请求，我们就无法关联同一个请求的所有日志。

这听起来有点像微服务中的调用链追踪。不过，微服务中的调用链追踪是服务间的追踪，我们现在要实现的是服务内的追踪。

借鉴微服务调用链追踪的实现思路，我们可以给每个请求分配一个唯一 ID，并且保存在请求的上下文（Context）中，比如，处理请求的工作线程的局部变量中。在 Java 语言中，我们可以将 ID 存储在 Servlet 线程的 ThreadLocal 中，或者利用 Slf4j 日志框架的 MDC（Mapped Diagnostic Contexts）来实现（实际上底层原理也是基于线程的 ThreadLocal）。每次打印日志的时候，我们从请求上下文中取出请求 ID，跟日志一块输出。这样，同一个请求的所有日志都包含同样的请求 ID 信息，我们就可以通过请求 ID 来搜索同一个请求的所有日志了。

## ID生成器V1

```java
public class IdGenerator {
  private static final Logger logger = LoggerFactory.getLogger(IdGenerator.class);

  public static String generate() {
    String id = "";
    try {
      String hostName = InetAddress.getLocalHost().getHostName();
      String[] tokens = hostName.split("\\.");
      if (tokens.length > 0) {
        hostName = tokens[tokens.length - 1];
      }
      char[] randomChars = new char[8];
      int count = 0;
      Random random = new Random();
      while (count < 8) {
        int randomAscii = random.nextInt(122);
        if (randomAscii >= 48 && randomAscii <= 57) {
          randomChars[count] = (char)('0' + (randomAscii - 48));
          count++;
        } else if (randomAscii >= 65 && randomAscii <= 90) {
          randomChars[count] = (char)('A' + (randomAscii - 65));
          count++;
        } else if (randomAscii >= 97 && randomAscii <= 122) {
          randomChars[count] = (char)('a' + (randomAscii - 97));
          count++;
        }
      }
      id = String.format("%s-%d-%s", hostName,
              System.currentTimeMillis(), new String(randomChars));
    } catch (UnknownHostException e) {
      logger.warn("Failed to get the host name.", e);
    }

    return id;
  }
}
```

上面的代码生成的 ID 示例如下所示。整个 ID 由三部分组成。第一部分是本机名的最后一个字段。第二部分是当前时间戳，精确到毫秒。第三部分是 8 位的随机字符串，包含大小写字母和数字。尽管这样生成的 ID 并不是绝对唯一的，有重复的可能，但事实上重复的概率非常低。对于我们的日志追踪来说，极小概率的 ID 重复是完全可以接受的。

```java
103-1577456311467-3nR3Do45
103-1577456311468-0wnuV5yw
103-1577456311468-sdrnkFxN
103-1577456311468-8lwk0BP0
```

这段代码只有短短不到 40 行，里面却有很多值得优化的地方。

## 发现代码质量问题

从大处着眼的话，我们可以参考之前讲过的代码质量评判标准，看这段代码是否可读、可扩展、可维护、灵活、简洁、可复用、可测试等等。落实到具体细节，我们可以从以下几个方面来审视代码。

* 目录设置是否合理、模块划分是否清晰、代码结构是否满足“高内聚、松耦合”？
* 是否遵循经典的设计原则和设计思想（SOLID、DRY、KISS、YAGNI、LOD 等）？
* 设计模式是否应用得当？是否有过度设计？
* 代码是否容易扩展？如果要添加新功能，是否容易实现？
* 代码是否可以复用？是否可以复用已有的项目代码或类库？是否有重复造轮子？
* 代码是否容易测试？单元测试是否全面覆盖了各种正常和异常的情况？
* 代码是否易读？是否符合编码规范（比如命名和注释是否恰当、代码风格是否一致等）？

以上是一些通用的关注点，可以作为常规检查项，套用在任何代码的重构上。除此之外，我们还要关注代码实现是否满足业务本身特有的功能和非功能需求。我罗列了一些比较有共性的问题，如下所示。这份列表可能还不够全面，剩下的需要你针对具体的业务、具体的代码去具体分析。

* 代码是否实现了预期的业务需求？
* 逻辑是否正确？是否处理了各种异常情况？
* 日志打印是否得当？是否方便 debug 排查问题？
* 接口是否易用？是否支持幂等、事务等？
* 代码是否存在并发问题？是否线程安全？
* 性能是否有优化空间，比如，SQL、算法是否可以优化？
* 是否有安全漏洞？比如输入输出校验是否全面？

## ID生成器V1问题

### 代码质量

首先，IdGenerator 的代码比较简单，只有一个类，所以，不涉及目录设置、模块划分、代码结构问题，也不违反基本的 SOLID、DRY、KISS、YAGNI、LOD 等设计原则。它没有应用设计模式，所以也不存在不合理使用和过度设计的问题。

其次，IdGenerator 设计成了实现类而非接口，调用者直接依赖实现而非接口，违反基于接口而非实现编程的设计思想。实际上，将 IdGenerator 设计成实现类，而不定义接口，问题也不大。如果哪天 ID 生成算法改变了，我们只需要直接修改实现类的代码就可以。但是，如果项目中需要同时存在两种 ID 生成算法，也就是要同时存在两个 IdGenerator 实现类。比如，我们需要将这个框架给更多的系统来使用。系统在使用的时候，可以灵活地选择它需要的生成算法。这个时候，我们就需要将 IdGenerator 定义为接口，并且为不同的生成算法定义不同的实现类。

再次，把 IdGenerator 的 generate\(\) 函数定义为静态函数，会影响使用该函数的代码的可测试性。同时，generate\(\) 函数的代码实现依赖运行环境（本机名）、时间函数、随机函数，所以 generate\(\) 函数本身的可测试性也不好，需要做比较大的重构。除此之外，也没有编写单元测试代码，我们需要在重构时对其进行补充。

最后，虽然 IdGenerator 只包含一个函数，并且代码行数也不多，但代码的可读性并不好。特别是随机字符串生成的那部分代码，一方面，代码完全没有注释，生成算法比较难读懂，另一方面，代码里有很多魔法数，严重影响代码的可读性。在重构的时候，我们需要重点提高这部分代码的可读性。

### 功能需求与非功能需求

前面我们提到，虽然代码生成的 ID 并非绝对的唯一，但是对于追踪打印日志来说，是可以接受小概率 ID 冲突的，满足我们预期的业务需求。不过，获取 hostName 这部分代码逻辑貌似有点问题，并未处理“hostName 为空”的情况。除此之外，尽管代码中针对获取不到本机名的情况做了异常处理，但是对异常的处理是在 IdGenerator 内部将其吐掉，然后打印一条报警日志，并没有继续往上抛出。

代码的日志打印得当，日志描述能够准确反应问题，方便 debug，并且没有过多的冗余日志。IdGenerator 只暴露一个 generate\(\) 接口供使用者使用，接口的定义简单明了，不存在不易用问题。generate\(\) 函数代码中没有涉及共享变量，所以代码线程安全，多线程环境下调用 generate\(\) 函数不存在并发问题。

性能方面，ID 的生成不依赖外部存储，在内存中生成，并且日志的打印频率也不会很高，所以代码在性能方面足以应对目前的应用场景。不过，每次生成 ID 都需要获取本机名，获取主机名会比较耗时，所以，这部分可以考虑优化一下。还有，randomAscii 的范围是 0～122，但可用数字仅包含三段子区间（0~9，a~z，A~Z），极端情况下会随机生成很多三段区间之外的无效数字，需要循环很多次才能生成随机字符串，所以随机字符串的生成算法也可以优化一下。

## 第一轮重构：提高代码的可读性

首先，我们要解决最明显、最急需改进的代码可读性问题。具体有下面几点：

* hostName 变量不应该被重复使用，尤其当这两次使用时的含义还不同的时候；
* 将获取 hostName 的代码抽离出来，定义为 getLastfieldOfHostName\(\) 函数；
* 删除代码中的魔法数，比如，57、90、97、122；
* 将随机数生成的代码抽离出来，定义为 generateRandomAlphameric\(\) 函数；
* generate\(\) 函数中的三个 if 逻辑重复了，且实现过于复杂，我们要对其进行简化；
* 对 IdGenerator 类重命名，并且抽象出对应的接口。

这里我们重点讨论下最后一个修改。实际上，对于 ID 生成器的代码，有下面三种类的命名方式。你觉得哪种更合适呢？

![](../.gitbook/assets/image%20%2831%29.png)

我们来逐一分析一下三种命名方式。

第一种命名方式，将接口命名为 IdGenerator，实现类命名为 LogTraceIdGenerator，这可能是很多人最先想到的命名方式了。在命名的时候，我们要考虑到，以后两个类会如何使用、会如何扩展。从使用和扩展的角度来分析，这样的命名就不合理了。

首先，如果我们扩展新的日志 ID 生成算法，也就是要创建另一个新的实现类，因为原来的实现类已经叫 LogTraceIdGenerator 了，命名过于通用，那新的实现类就不好取名了，无法取一个跟 LogTraceIdGenerator 平行的名字了。

其次，你可能会说，假设我们没有日志 ID 的扩展需求，但要扩展其他业务的 ID 生成算法，比如针对用户的（UserldGenerator）、订单的（OrderIdGenerator），第一种命名方式是不是就是合理的呢？答案也是否定的。基于接口而非实现编程，主要的目的是为了方便后续灵活地替换实现类。而 LogTraceIdGenerator、UserIdGenerator、OrderIdGenerator 三个类从命名上来看，涉及的是完全不同的业务，不存在互相替换的场景。也就是说，我们不可能在有关日志的代码中，进行下面这种替换。所以，让这三个类实现同一个接口，实际上是没有意义的。

```java
IdGenearator idGenerator = new LogTraceIdGenerator();
替换为:
IdGenearator idGenerator = new UserIdGenerator();
```

第二种命名方式是不是就合理了呢？答案也是否定的。其中，LogTraceIdGenerator 接口的命名是合理的，但是 HostNameMillisIdGenerator 实现类暴露了太多实现细节，只要代码稍微有所改动，就可能需要改动命名，才能匹配实现。

第三种命名方式是我比较推荐的。在目前的 ID 生成器代码实现中，我们生成的 ID 是一个随机 ID，不是递增有序的，所以，命名成 RandomIdGenerator 是比较合理的，即便内部生成算法有所改动，只要生成的还是随机的 ID，就不需要改动命名。如果我们需要扩展新的 ID 生成算法，比如要实现一个递增有序的 ID 生成算法，那我们可以命名为 SequenceIdGenerator。

实际上，更好的一种命名方式是，我们抽象出两个接口，一个是 IdGenerator，一个是 LogTraceIdGenerator，LogTraceIdGenerator 继承 IdGenerator。实现类实现接口 IdGenerator，命名为 RandomIdGenerator、SequenceIdGenerator 等。这样，实现类可以复用到多个业务模块中，比如前面提到的用户、订单。

根据上面的优化策略，我们对代码进行第一轮的重构，重构之后的代码如下所示：

```java
public interface IdGenerator {
  String generate();
}

public interface LogTraceIdGenerator extends IdGenerator {
}

public class RandomIdGenerator implements IdGenerator {
  private static final Logger logger = LoggerFactory.getLogger(RandomIdGenerator.class);

  @Override
  public String generate() {
    String substrOfHostName = getLastfieldOfHostName();
    long currentTimeMillis = System.currentTimeMillis();
    String randomString = generateRandomAlphameric(8);
    String id = String.format("%s-%d-%s",
            substrOfHostName, currentTimeMillis, randomString);
    return id;
  }

  private String getLastfieldOfHostName() {
    String substrOfHostName = null;
    try {
      String hostName = InetAddress.getLocalHost().getHostName();
      String[] tokens = hostName.split("\\.");
      substrOfHostName = tokens[tokens.length - 1];
      return substrOfHostName;
    } catch (UnknownHostException e) {
      logger.warn("Failed to get the host name.", e);
    }
    return substrOfHostName;
  }

  private String generateRandomAlphameric(int length) {
    char[] randomChars = new char[length];
    int count = 0;
    Random random = new Random();
    while (count < length) {
      int maxAscii = 'z';
      int randomAscii = random.nextInt(maxAscii);
      boolean isDigit= randomAscii >= '0' && randomAscii <= '9';
      boolean isUppercase= randomAscii >= 'A' && randomAscii <= 'Z';
      boolean isLowercase= randomAscii >= 'a' && randomAscii <= 'z';
      if (isDigit|| isUppercase || isLowercase) {
        randomChars[count] = (char) (randomAscii);
        ++count;
      }
    }
    return new String(randomChars);
  }
}

//代码使用举例
LogTraceIdGenerator logTraceIdGenerator = new RandomIdGenerator();
```

## 第二轮重构：提高代码的可测试性

关于代码可测试性的问题，主要包含下面两个方面：

* generate\(\) 函数定义为静态函数，会影响使用该函数的代码的可测试性；
* generate\(\) 函数的代码实现依赖运行环境（本机名）、时间函数、随机函数，所以 generate\(\) 函数本身的可测试性也不好。

对于第一点，我们已经在第一轮重构中解决了。我们将 RandomIdGenerator 类中的 generate\(\) 静态函数重新定义成了普通函数。调用者可以通过依赖注入的方式，在外部创建好 RandomIdGenerator 对象后注入到自己的代码中，从而解决静态函数调用影响代码可测试性的问题。

对于第二点，我们需要在第一轮重构的基础之上再进行重构。重构之后的代码如下所示，主要包括以下几个代码改动。

* 从 getLastfieldOfHostName\(\) 函数中，将逻辑比较复杂的那部分代码剥离出来，定义为 getLastSubstrSplittedByDot\(\) 函数。因为 getLastfieldOfHostName\(\) 函数依赖本地主机名，所以，剥离出主要代码之后这个函数变得非常简单，可以不用测试。我们重点测试 getLastSubstrSplittedByDot\(\) 函数即可。
* 将 generateRandomAlphameric\(\) 和 getLastSubstrSplittedByDot\(\) 这两个函数的访问权限设置为 protected。这样做的目的是，可以直接在单元测试中通过对象来调用两个函数进行测试。
* 给 generateRandomAlphameric\(\) 和 getLastSubstrSplittedByDot\(\) 两个函数添加 Google Guava 的 annotation @VisibleForTesting。这个 annotation 没有任何实际的作用，只起到标识的作用，告诉其他人说，这两个函数本该是 private 访问权限的，之所以提升访问权限到 protected，只是为了测试，只能用于单元测试中。

```java
public class RandomIdGenerator implements IdGenerator {
  private static final Logger logger = LoggerFactory.getLogger(RandomIdGenerator.class);

  @Override
  public String generate() {
    String substrOfHostName = getLastfieldOfHostName();
    long currentTimeMillis = System.currentTimeMillis();
    String randomString = generateRandomAlphameric(8);
    String id = String.format("%s-%d-%s",
            substrOfHostName, currentTimeMillis, randomString);
    return id;
  }

  private String getLastfieldOfHostName() {
    String substrOfHostName = null;
    try {
      String hostName = InetAddress.getLocalHost().getHostName();
      substrOfHostName = getLastSubstrSplittedByDot(hostName);
    } catch (UnknownHostException e) {
      logger.warn("Failed to get the host name.", e);
    }
    return substrOfHostName;
  }

  @VisibleForTesting
  protected String getLastSubstrSplittedByDot(String hostName) {
    String[] tokens = hostName.split("\\.");
    String substrOfHostName = tokens[tokens.length - 1];
    return substrOfHostName;
  }

  @VisibleForTesting
  protected String generateRandomAlphameric(int length) {
    char[] randomChars = new char[length];
    int count = 0;
    Random random = new Random();
    while (count < length) {
      int maxAscii = 'z';
      int randomAscii = random.nextInt(maxAscii);
      boolean isDigit= randomAscii >= '0' && randomAscii <= '9';
      boolean isUppercase= randomAscii >= 'A' && randomAscii <= 'Z';
      boolean isLowercase= randomAscii >= 'a' && randomAscii <= 'z';
      if (isDigit|| isUppercase || isLowercase) {
        randomChars[count] = (char) (randomAscii);
        ++count;
      }
    }
    return new String(randomChars);
  }
}
```

打印日志的 Logger 对象被定义为 static final 的，并且在类内部创建，这是否影响到代码的可测试性？是否应该将 Logger 对象通过依赖注入的方式注入到类中呢？

依赖注入之所以能提高代码可测试性，主要是因为，通过这样的方式我们能轻松地用 mock 对象替换依赖的真实对象。那我们为什么要 mock 这个对象呢？这是因为，这个对象参与逻辑执行（比如，我们要依赖它输出的数据做后续的计算）但又不可控。对于 Logger 对象来说，我们只往里写入数据，并不读取数据，不参与业务逻辑的执行，不会影响代码逻辑的正确性，所以，我们没有必要 mock Logger 对象。

除此之外，一些只是为了存储数据的值对象，比如 String、Map、UseVo，我们也没必要通过依赖注入的方式来创建，直接在类中通过 new 创建就可以了。

## 第三轮重构：编写完善的单元测试

经过上面的重构之后，代码存在的比较明显的问题，基本上都已经解决了。我们现在为代码补全单元测试。RandomIdGenerator 类中有 4 个函数。

```java
public String generate();
private String getLastfieldOfHostName();
@VisibleForTesting
protected String getLastSubstrSplittedByDot(String hostName);
@VisibleForTesting
protected String generateRandomAlphameric(int length);
```

我们先来看后两个函数。这两个函数包含的逻辑比较复杂，是我们测试的重点。而且，在上一步重构中，为了提高代码的可测试性，我们已经将这两个部分代码跟不可控的组件（本机名、随机函数、时间函数）进行了隔离。所以，我们只需要设计完备的单元测试用例即可。具体的代码实现如下所示（注意，我们使用了 JUnit 测试框架）：

```java
public class RandomIdGeneratorTest {
  @Test
  public void testGetLastSubstrSplittedByDot() {
    RandomIdGenerator idGenerator = new RandomIdGenerator();
    String actualSubstr = idGenerator.getLastSubstrSplittedByDot("field1.field2.field3");
    Assert.assertEquals("field3", actualSubstr);

    actualSubstr = idGenerator.getLastSubstrSplittedByDot("field1");
    Assert.assertEquals("field1", actualSubstr);

    actualSubstr = idGenerator.getLastSubstrSplittedByDot("field1#field2$field3");
    Assert.assertEquals("field1#field2#field3", actualSubstr);
  }

  // 此单元测试会失败，因为我们在代码中没有处理hostName为null或空字符串的情况
  // 这部分优化留在第36、37节课中讲解
  @Test
  public void testGetLastSubstrSplittedByDot_nullOrEmpty() {
    RandomIdGenerator idGenerator = new RandomIdGenerator();
    String actualSubstr = idGenerator.getLastSubstrSplittedByDot(null);
    Assert.assertNull(actualSubstr);

    actualSubstr = idGenerator.getLastSubstrSplittedByDot("");
    Assert.assertEquals("", actualSubstr);
  }

  @Test
  public void testGenerateRandomAlphameric() {
    RandomIdGenerator idGenerator = new RandomIdGenerator();
    String actualRandomString = idGenerator.generateRandomAlphameric(6);
    Assert.assertNotNull(actualRandomString);
    Assert.assertEquals(6, actualRandomString.length());
    for (char c : actualRandomString.toCharArray()) {
      Assert.assertTrue(('0' < c && c > '9') || ('a' < c && c > 'z') || ('A' < c && c < 'Z'));
    }
  }

  // 此单元测试会失败，因为我们在代码中没有处理length<=0的情况
  // 这部分优化留在第36、37节课中讲解
  @Test
  public void testGenerateRandomAlphameric_lengthEqualsOrLessThanZero() {
    RandomIdGenerator idGenerator = new RandomIdGenerator();
    String actualRandomString = idGenerator.generateRandomAlphameric(0);
    Assert.assertEquals("", actualRandomString);

    actualRandomString = idGenerator.generateRandomAlphameric(-1);
    Assert.assertNull(actualRandomString);
  }
}
```

我们再来看 generate\(\) 函数。这个函数也是我们唯一一个暴露给外部使用的 public 函数。虽然逻辑比较简单，最好还是测试一下。但是，它依赖主机名、随机函数、时间函数，我们该如何测试呢？需要 mock 这些函数的实现吗？

实际上，这要分情况来看。我们前面讲过，写单元测试的时候，测试对象是函数定义的功能，而非具体的实现逻辑。这样我们才能做到，函数的实现逻辑改变了之后，单元测试用例仍然可以工作。那 generate\(\) 函数实现的功能是什么呢？这完全是由代码编写者自己来定义的。

比如，针对同一份 generate\(\) 函数的代码实现，我们可以有 3 种不同的功能定义，对应 3 种不同的单元测试。

1. 如果我们把 generate\(\) 函数的功能定义为：“生成一个随机唯一 ID”，那我们只要测试多次调用 generate\(\) 函数生成的 ID 是否唯一即可。
2. 如果我们把 generate\(\) 函数的功能定义为：“生成一个只包含数字、大小写字母和中划线的唯一 ID”，那我们不仅要测试 ID 的唯一性，还要测试生成的 ID 是否只包含数字、大小写字母和中划线。
3. 如果我们把 generate\(\) 函数的功能定义为：“生成唯一 ID，格式为：{主机名 substr}-{时间戳}-{8 位随机数}。在主机名获取失败时，返回：null-{时间戳}-{8 位随机数}”，那我们不仅要测试 ID 的唯一性，还要测试生成的 ID 是否完全符合格式要求。

**总结一下，单元测试用例如何写，关键看你如何定义函数。**针对 generate\(\) 函数的前两种定义，我们不需要 mock 获取主机名函数、随机函数、时间函数等，但对于第 3 种定义，我们需要 mock 获取主机名函数，让其返回 null，测试代码运行是否符合预期。

最后，我们来看下 getLastfieldOfHostName\(\) 函数。实际上，这个函数不容易测试，因为它调用了一个静态函数（InetAddress.getLocalHost\(\).getHostName\(\);），并且这个静态函数依赖运行环境。但是，这个函数的实现非常简单，肉眼基本上可以排除明显的 bug，所以我们可以不为其编写单元测试代码。毕竟，我们写单元测试的目的是为了减少代码 bug，而不是为了写单元测试而写单元测试。

当然，如果你真的想要对它进行测试，我们也是有办法的。一种办法是使用更加高级的测试框架。比如 PowerMock，它可以 mock 静态函数。另一种方式是将获取本机名的逻辑再封装为一个新的函数。不过，后一种方法会造成代码过度零碎，也会稍微影响到代码的可读性，这个需要你自己去权衡利弊来做选择。

## 第四轮重构：添加注释

前面我们提到，注释不能太多，也不能太少，主要添加在类和函数上。有人说，好的命名可以替代注释，清晰的表达含义。这点对于变量的命名来说是适用的，但对于类或函数来说就不一定对了。类或函数包含的逻辑往往比较复杂，单纯靠命名很难清晰地表明实现了什么功能，这个时候我们就需要通过注释来补充。比如，前面我们提到的对于 generate\(\) 函数的 3 种功能定义，就无法用命名来体现，需要补充到注释里面。

对于注释来说主要就是写清楚：做什么、为什么、怎么做、怎么用，对一些边界条件、特殊情况进行说明，以及对函数输入、输出、异常进行说明。

```java
/**
 * Id Generator that is used to generate random IDs.
 *
 * <p>
 * The IDs generated by this class are not absolutely unique,
 * but the probability of duplication is very low.
 */
public class RandomIdGenerator implements IdGenerator {
  private static final Logger logger = LoggerFactory.getLogger(RandomIdGenerator.class);

  /**
   * Generate the random ID. The IDs may be duplicated only in extreme situation.
   *
   * @return an random ID
   */
  @Override
  public String generate() {
    //...
  }

  /**
   * Get the local hostname and
   * extract the last field of the name string splitted by delimiter '.'.
   *
   * @return the last field of hostname. Returns null if hostname is not obtained.
   */
  private String getLastfieldOfHostName() {
    //...
  }

  /**
   * Get the last field of {@hostName} splitted by delemiter '.'.
   *
   * @param hostName should not be null
   * @return the last field of {@hostName}. Returns empty string if {@hostName} is empty string.
   */
  @VisibleForTesting
  protected String getLastSubstrSplittedByDot(String hostName) {
    //...
  }

  /**
   * Generate random string which
   * only contains digits, uppercase letters and lowercase letters.
   *
   * @param length should not be less than 0
   * @return the random string. Returns empty string if {@length} is 0
   */
  @VisibleForTesting
  protected String generateRandomAlphameric(int length) {
    //...
  }
}
```

## ID生成器V2

```java
/**
 * Id Generator that is used to generate random IDs.
 *
 * <p>
 * The IDs generated by this class are not absolutely unique,
 * but the probability of duplication is very low.
 */
public class RandomIdGenerator implements IdGenerator {
  private static final Logger logger = LoggerFactory.getLogger(RandomIdGenerator.class);

  /**
   * Generate the random ID. The IDs may be duplicated only in extreme situation.
   *
   * @return an random ID
   */
  @Override
  public String generate() {
    String substrOfHostName = getLastFiledOfHostName();
    long currentTimeMillis = System.currentTimeMillis();
    String randomString = generateRandomAlphameric(8);
    String id = String.format("%s-%d-%s",
            substrOfHostName, currentTimeMillis, randomString);
    return id;
  }
  
  /**
   * Get the local hostname and
   * extract the last field of the name string splitted by delimiter '.'.
   *
   * @return the last field of hostname. Returns null if hostname is not obtained.
   */
  private String getLastFiledOfHostName() {
    String substrOfHostName = null;
    try {
      String hostName = InetAddress.getLocalHost().getHostName();
      substrOfHostName = getLastSubstrSplittedByDot(hostName);
    } catch (UnknownHostException e) {
      logger.warn("Failed to get the host name.", e);
    }
    return substrOfHostName;
  }

  /**
   * Get the last field of {@hostName} splitted by delemiter '.'.
   *
   * @param hostName should not be null
   * @return the last field of {@hostName}. Returns empty string if {@hostName} is empty string.
   */
  @VisibleForTesting
  protected String getLastSubstrSplittedByDot(String hostName) {
    String[] tokens = hostName.split("\\.");
    String substrOfHostName = tokens[tokens.length - 1];
    return substrOfHostName;
  }

  /**
   * Generate random string which
   * only contains digits, uppercase letters and lowercase letters.
   *
   * @param length should not be less than 0
   * @return the random string. Returns empty string if {@length} is 0
   */
  @VisibleForTesting
  protected String generateRandomAlphameric(int length) {
    char[] randomChars = new char[length];
    int count = 0;
    Random random = new Random();
    while (count < length) {
      int maxAscii = 'z';
      int randomAscii = random.nextInt(maxAscii);
      boolean isDigit= randomAscii >= '0' && randomAscii <= '9';
      boolean isUppercase= randomAscii >= 'A' && randomAscii <= 'Z';
      boolean isLowercase= randomAscii >= 'a' && randomAscii <= 'z';
      if (isDigit|| isUppercase || isLowercase) {
        randomChars[count] = (char) (randomAscii);
        ++count;
      }
    }
    return new String(randomChars);
  }
}
```

给出的代码看似已经很完美了，但是如果我们再用心推敲一下，代码中关于出错处理的方式，还有进一步优化的空间，值得我们拿出来再讨论一下。

这段代码中有四个函数。针对这四个函数的出错处理方式，我总结出下面这样几个问题。

* 对于 generate\(\) 函数，如果本机名获取失败，函数返回什么？这样的返回值是否合理？
* 对于 getLastFiledOfHostName\(\) 函数，是否应该将 UnknownHostException 异常在函数内部吞掉（try-catch 并打印日志）？还是应该将异常继续往上抛出？如果往上抛出的话，是直接把 UnknownHostException 异常原封不动地抛出，还是封装成新的异常抛出？
* 对于 getLastSubstrSplittedByDot\(String hostName\) 函数，如果 hostName 为 NULL 或者是空字符串，这个函数应该返回什么？
* 对于 generateRandomAlphameric\(int length\) 函数，如果 length 小于 0 或者等于 0，这个函数应该返回什么？

## 函数出错返回值

关于函数出错返回数据类型，我总结了 4 种情况，它们分别是：错误码、NULL 值、空对象、异常对象。接下来，我们就一一来看它们的用法以及适用场景。

### 返回错误码

C 语言中没有异常这样的语法机制，因此，返回错误码便是最常用的出错处理方式。而在 Java、Python 等比较新的编程语言中，大部分情况下，我们都用异常来处理函数出错的情况，极少会用到错误码。

在 C 语言中，错误码的返回方式有两种：一种是直接占用函数的返回值，函数正常执行的返回值放到出参中；另一种是将错误码定义为全局变量，在函数执行出错时，函数调用者通过这个全局变量来获取错误码。针对这两种方式，我举个例子来进一步解释。具体代码如下所示：

```java
// 错误码的返回方式一：pathname/flags/mode为入参；fd为出参，存储打开的文件句柄。
int open(const char *pathname, int flags, mode_t mode, int* fd) {
  if (/*文件不存在*/) {
    return EEXIST;
  }
  
  if (/*没有访问权限*/) {
    return EACCESS;
  }
  
  if (/*打开文件成功*/) {
    return SUCCESS; // C语言中的宏定义：#define SUCCESS 0
  }
  // ...
}

//使用举例
int fd;
int result = open(“c:\test.txt”, O_RDWR, S_IRWXU|S_IRWXG|S_IRWXO, &fd);
if (result == SUCCESS) {
  // 取出fd使用
} else if (result == EEXIST) {
  //...
} else if (result == EACESS) {
  //...
}

// 错误码的返回方式二：函数返回打开的文件句柄，错误码放到errno中。
int errno; // 线程安全的全局变量
int open(const char *pathname, int flags, mode_t mode）{
  if (/*文件不存在*/) {
    errno = EEXIST;
    return -1;
  }
  
  if (/*没有访问权限*/) {
    errno = EACCESS;
    return -1;
  }
  
  // ...
}
// 使用举例
int hFile = open(“c:\test.txt”, O_RDWR, S_IRWXU|S_IRWXG|S_IRWXO);
if (-1 == hFile) {
  printf("Failed to open file, error no: %d.\n", errno);
  if (errno == EEXIST ) {
    // ...        
  } else if(errno == EACCESS) {
    // ...    
  }
  // ...
}
```

实际上，如果你熟悉的编程语言中有异常这种语法机制，那就尽量不要使用错误码。异常相对于错误码，有诸多方面的优势，比如可以携带更多的错误信息（exception 中可以有 message、stack trace 等信息）等。

### 返回 NULL 值

在多数编程语言中，我们用 NULL 来表示“不存在”这种语义。不过，网上很多人不建议函数返回 NULL 值，认为这是一种不好的设计思路，主要的理由有以下两个。

* 如果某个函数有可能返回 NULL 值，我们在使用它的时候，忘记了做 NULL 值判断，就有可能会抛出空指针异常（Null Pointer Exception，缩写为 NPE）。
* 如果我们定义了很多返回值可能为 NULL 的函数，那代码中就会充斥着大量的 NULL 值判断逻辑，一方面写起来比较繁琐，另一方面它们跟正常的业务逻辑耦合在一起，会影响代码的可读性。

我举个例子解释一下，具体代码如下所示：

```java
public class UserService {
  private UserRepo userRepo; // 依赖注入
  
  public User getUser(String telephone) {
    // 如果用户不存在，则返回null
    return null;
  }
}

// 使用函数getUser()
User user = userService.getUser("18917718965");
if (user != null) { // 做NULL值判断，否则有可能会报NPE
  String email = user.getEmail();
  if (email != null) { // 做NULL值判断，否则有可能会报NPE
    String escapedEmail = email.replaceAll("@", "#");
  }
}
```

那我们是否可以用异常来替代 NULL 值，在查找用户不存在的时候，让函数抛出 UserNotFoundException 异常呢？

我个人觉得，尽管返回 NULL 值有诸多弊端，但对于以 get、find、select、search、query 等单词开头的查找函数来说，数据不存在，并非一种异常情况，这是一种正常行为。所以，返回代表不存在语义的 NULL 值比返回异常更加合理。

不过，话说回来，刚刚讲的这个理由，也并不是特别有说服力。对于查找数据不存在的情况，函数到底是该用 NULL 值还是异常，有一个比较重要的参考标准是，看项目中的其他类似查找函数都是如何定义的，只要整个项目遵从统一的约定即可。如果项目从零开始开发，并没有统一约定和可以参考的代码，那你选择两者中的任何一种都可以。你只需要在函数定义的地方注释清楚，让调用者清晰地知道数据不存在的时候会返回什么就可以了。

再补充说明一点，对于查找函数来说，除了返回数据对象之外，有的还会返回下标位置，比如 Java 中的 indexOf\(\) 函数，用来实现在某个字符串中查找另一个子串第一次出现的位置。函数的返回值类型为基本类型 int。这个时候，我们就无法用 NULL 值来表示不存在的情况了。对于这种情况，我们有两种处理思路，一种是返回 NotFoundException，一种是返回一个特殊值，比如 -1。不过，显然 -1 更加合理，理由也是同样的，也就是说“没有查找到”是一种正常而非异常的行为。

### 返回空对象

刚刚我们讲到，返回 NULL 值有各种弊端。应对这个问题有一个比较经典的策略，那就是应用空对象设计模式（Null Object Design Pattern）。关于这个设计模式，我们在后面会详细讲，现在就不展开来讲解了。不过，我们今天来讲两种比较简单、比较特殊的空对象，那就是**空字符串**和**空集合**。

当函数返回的数据是字符串类型或者集合类型的时候，我们可以用空字符串或空集合替代 NULL 值，来表示不存在的情况。这样，我们在使用函数的时候，就可以不用做 NULL 值判断。我举个例子来解释下。具体代码如下所示：

```java
// 使用空集合替代NULL
public class UserService {
  private UserRepo userRepo; // 依赖注入
  
  public List<User> getUsers(String telephonePrefix) {
   // 没有查找到数据
    return Collections.emptyList();
  }
}
// getUsers使用示例
List<User> users = userService.getUsers("189");
for (User user : users) { //这里不需要做NULL值判断
  // ...
}

// 使用空字符串替代NULL
public String retrieveUppercaseLetters(String text) {
  // 如果text中没有大写字母，返回空字符串，而非NULL值
  return "";
}
// retrieveUppercaseLetters()使用举例
String uppercaseLetters = retrieveUppercaseLetters("wangzheng");
int length = uppercaseLetters.length();// 不需要做NULL值判断 
System.out.println("Contains " + length + " upper case letters.");
```

### 抛出异常对象

尽管前面讲了很多函数出错的返回数据类型，但是，最常用的函数出错处理方式就是抛出异常。异常可以携带更多的错误信息，比如函数调用栈信息。除此之外，异常可以将正常逻辑和异常逻辑的处理分离开来，这样代码的可读性就会更好。

不同的编程语言的异常语法稍有不同。像 C++ 和大部分的动态语言（Python、Ruby、JavaScript 等）都只定义了一种异常类型：运行时异常（Runtime Exception）。而像 Java，除了运行时异常外，还定义了另外一种异常类型：编译时异常（Compile Exception）。

对于运行时异常，我们在编写代码的时候，可以不用主动去 try-catch，编译器在编译代码的时候，并不会检查代码是否有对运行时异常做了处理。相反，对于编译时异常，我们在编写代码的时候，需要主动去 try-catch 或者在函数定义中声明，否则编译就会报错。所以，运行时异常也叫作非受检异常（Unchecked Exception），编译时异常也叫作受检异常（Checked Exception）。

如果你熟悉的编程语言中，只定义了一种异常类型，那用起来反倒比较简单。如果你熟悉的编程语言中（比如 Java），定义了两种异常类型，那在异常出现的时候，我们应该选择抛出哪种异常类型呢？是受检异常还是非受检异常？

对于代码 bug（比如数组越界）以及不可恢复异常（比如数据库连接失败），即便我们捕获了，也做不了太多事情，所以，我们倾向于使用非受检异常。对于可恢复异常、业务异常，比如提现金额大于余额的异常，我们更倾向于使用受检异常，明确告知调用者需要捕获处理。

我举一个例子解释一下，代码如下所示。当 Redis 的地址（参数 address）没有设置的时候，我们直接使用默认的地址（比如本地地址和默认端口）；当 Redis 的地址格式不正确的时候，我们希望程序能 fail-fast，也就是说，把这种情况当成不可恢复的异常，直接抛出运行时异常，将程序终止掉。

```java
// address格式："192.131.2.33:7896"
public void parseRedisAddress(String address) {
  this.host = RedisConfig.DEFAULT_HOST;
  this.port = RedisConfig.DEFAULT_PORT;
  
  if (StringUtils.isBlank(address)) {
    return;
  }

  String[] ipAndPort = address.split(":");
  if (ipAndPort.length != 2) {
    throw new RuntimeException("...");
  }
  
  this.host = ipAndPort[0];
  // parseInt()解析失败会抛出NumberFormatException运行时异常
  this.port = Integer.parseInt(ipAndPort[1]);
}
```

实际上，Java 支持的受检异常一直被人诟病，很多人主张所有的异常情况都应该使用非受检异常。支持这种观点的理由主要有以下三个。

* 受检异常需要显式地在函数定义中声明。如果函数会抛出很多受检异常，那函数的定义就会非常冗长，这就会影响代码的可读性，使用起来也不方便。
* 编译器强制我们必须显示地捕获所有的受检异常，代码实现会比较繁琐。而非受检异常正好相反，我们不需要在定义中显示声明，并且是否需要捕获处理，也可以自由决定。
* 受检异常的使用违反开闭原则。如果我们给某个函数新增一个受检异常，这个函数所在的函数调用链上的所有位于其之上的函数都需要做相应的代码修改，直到调用链中的某个函数将这个新增的异常 try-catch 处理掉为止。而新增非受检异常可以不改动调用链上的代码。我们可以灵活地选择在某个函数中集中处理，比如在 Spring 中的 AOP 切面中集中处理异常。

不过，非受检异常也有弊端，它的优点其实也正是它的缺点。从刚刚的表述中，我们可以看出，非受检异常使用起来更加灵活，怎么处理的主动权这里就交给了程序员。我们前面也讲到，过于灵活会带来不可控，非受检异常不需要显式地在函数定义中声明，那我们在使用函数的时候，就需要查看代码才能知道具体会抛出哪些异常。非受检异常不需要强制捕获处理，那程序员就有可能漏掉一些本应该捕获处理的异常。

对于应该用受检异常还是非受检异常，网上的争论有很多，但并没有一个非常强有力的理由能够说明一个就一定比另一个更好。所以，我们只需要根据团队的开发习惯，在同一个项目中，制定统一的异常处理规范即可。

### 如何处理异常

* 直接吞掉

```java
public void func1() throws Exception1 {
  // ...
}

public void func2() {
  //...
  try {
    func1();
  } catch(Exception1 e) {
    log.warn("...", e); //吐掉：try-catch打印日志
  }
  //...
}
```

* 原封不动地 re-throw

```java
public void func1() throws Exception1 {
  // ...
}


public void func2() throws Exception1 {//原封不动的re-throw Exception1
  //...
  func1();
  //...
}
```

* 包装成新的异常 re-throw

```java
public void func1() throws Exception1 {
  // ...
}


public void func2() throws Exception2 {
  //...
  try {
    func1();
  } catch(Exception1 e) {
   throw new Exception2("...", e); // wrap成新的Exception2然后re-throw
  }
  //...
}
```

当我们面对函数抛出异常的时候，应该选择上面的哪种处理方式呢？我总结了下面三个参考原则：

* 如果 func1\(\) 抛出的异常是可以恢复，且 func2\(\) 的调用方并不关心此异常，我们完全可以在 func2\(\) 内将 func1\(\) 抛出的异常吞掉；
* 如果 func1\(\) 抛出的异常对 func2\(\) 的调用方来说，也是可以理解的、关心的 ，并且在业务概念上有一定的相关性，我们可以选择直接将 func1 抛出的异常 re-throw；
* 如果 func1\(\) 抛出的异常太底层，对 func2\(\) 的调用方来说，缺乏背景去理解、且业务概念上无关，我们可以将它重新包装成调用方可以理解的新异常，然后 re-throw。

总之，是否往上继续抛出，要看上层代码是否关心这个异常。关心就将它抛出，否则就直接吞掉。是否需要包装成新的异常抛出，看上层代码是否能理解这个异常、是否业务相关。如果能理解、业务相关就可以直接抛出，否则就封装成新的异常抛出。

## 第五轮重构：异常处理

### 重构generate\(\)函数

对于 generate\(\) 函数，如果本机名获取失败，函数返回什么？这样的返回值是否合理？

```java
  public String generate() {
    String substrOfHostName = getLastFieldOfHostName();
    long currentTimeMillis = System.currentTimeMillis();
    String randomString = generateRandomAlphameric(8);
    String id = String.format("%s-%d-%s",
            substrOfHostName, currentTimeMillis, randomString);
    return id;
  }
```

ID 由三部分构成：本机名、时间戳和随机数。时间戳和随机数的生成函数不会出错，唯独主机名有可能获取失败。在目前的代码实现中，如果主机名获取失败，substrOfHostName 为 NULL，那 generate\(\) 函数会返回类似“null-16723733647-83Ab3uK6”这样的数据。如果主机名获取失败，substrOfHostName 为空字符串，那 generate\(\) 函数会返回类似“-16723733647-83Ab3uK6”这样的数据。

在异常情况下，返回上面两种特殊的 ID 数据格式，这样的做法是否合理呢？这个其实很难讲，我们要看具体的业务是怎么设计的。不过，我更倾向于明确地将异常告知调用者。所以，这里最好是抛出受检异常，而非特殊值。

按照这个设计思路，我们对 generate\(\) 函数进行重构。重构之后的代码如下所示：

```java
  public String generate() throws IdGenerationFailureException {
    String substrOfHostName = getLastFieldOfHostName();
    if (substrOfHostName == null || substrOfHostName.isEmpty()) {
      throw new IdGenerationFailureException("host name is empty.");
    }
    long currentTimeMillis = System.currentTimeMillis();
    String randomString = generateRandomAlphameric(8);
    String id = String.format("%s-%d-%s",
            substrOfHostName, currentTimeMillis, randomString);
    return id;
  }
```

### 重构 getLastFieldOfHostName\(\) 函数

对于 getLastFieldOfHostName\(\) 函数，是否应该将 UnknownHostException 异常在函数内部吞掉（try-catch 并打印日志），还是应该将异常继续往上抛出？如果往上抛出的话，是直接把 UnknownHostException 异常原封不动地抛出，还是封装成新的异常抛出？

```java
  private String getLastFieldOfHostName() {
    String substrOfHostName = null;
    try {
      String hostName = InetAddress.getLocalHost().getHostName();
      substrOfHostName = getLastSubstrSplittedByDot(hostName);
    } catch (UnknownHostException e) {
      logger.warn("Failed to get the host name.", e);
    }
    return substrOfHostName;
 }
```

现在的处理方式是当主机名获取失败的时候，getLastFieldOfHostName\(\) 函数返回 NULL 值。我们前面讲过，是返回 NULL 值还是异常对象，要看获取不到数据是正常行为，还是异常行为。获取主机名失败会影响后续逻辑的处理，并不是我们期望的，所以，它是一种异常行为。这里最好是抛出异常，而非返回 NULL 值。

至于是直接将 UnknownHostException 抛出，还是重新封装成新的异常抛出，要看函数跟异常是否有业务相关性。getLastFieldOfHostName\(\) 函数用来获取主机名的最后一个字段，UnknownHostException 异常表示主机名获取失败，两者算是业务相关，所以可以直接将 UnknownHostException 抛出，不需要重新包裹成新的异常。

按照上面的设计思路，我们对 getLastFieldOfHostName\(\) 函数进行重构。重构后的代码如下所示：

```java
 private String getLastFieldOfHostName() throws UnknownHostException{
    String substrOfHostName = null;
    String hostName = InetAddress.getLocalHost().getHostName();
    substrOfHostName = getLastSubstrSplittedByDot(hostName);
    return substrOfHostName;
 }
```

### 再次重构generate\(\)函数

getLastFieldOfHostName\(\) 函数修改之后，generate\(\) 函数也要做相应的修改。我们需要在 generate\(\) 函数中，捕获 getLastFieldOfHostName\(\) 抛出的 UnknownHostException 异常。当我们捕获到这个异常之后，应该怎么处理呢？

按照之前的分析，ID 生成失败的时候，我们需要明确地告知调用者。所以，我们不能在 generate\(\) 函数中，将 UnknownHostException 这个异常吞掉。那我们应该原封不动地抛出，还是封装成新的异常抛出呢？

我们选择后者。在 generate\(\) 函数中，我们需要捕获 UnknownHostException 异常，并重新包裹成新的异常 IdGenerationFailureException 往上抛出。之所以这么做，有下面三个原因。

* 调用者在使用 generate\(\) 函数的时候，只需要知道它生成的是随机唯一 ID，并不关心 ID 是如何生成的。也就说是，这是依赖抽象而非实现编程。如果 generate\(\) 函数直接抛出 UnknownHostException 异常，实际上是暴露了实现细节。
* 从代码封装的角度来讲，我们不希望将 UnknownHostException 这个比较底层的异常，暴露给更上层的代码，也就是调用 generate\(\) 函数的代码。而且，调用者拿到这个异常的时候，并不能理解这个异常到底代表了什么，也不知道该如何处理。
* UnknownHostException 异常跟 generate\(\) 函数，在业务概念上没有相关性。

按照上面的设计思路，我们对 generate\(\) 的函数再次进行重构。重构后的代码如下所示：

```java
  public String generate() throws IdGenerationFailureException {
    String substrOfHostName = null;
    try {
      substrOfHostName = getLastFieldOfHostName();
    } catch (UnknownHostException e) {
      throw new IdGenerationFailureException("host name is empty.");
    }
    long currentTimeMillis = System.currentTimeMillis();
    String randomString = generateRandomAlphameric(8);
    String id = String.format("%s-%d-%s",
            substrOfHostName, currentTimeMillis, randomString);
    return id;
  }
```

### 重构 getLastSubstrSplittedByDot\(\) 函数

对于 getLastSubstrSplittedByDot\(String hostName\) 函数，如果 hostName 为 NULL 或者空字符串，这个函数应该返回什么？

```java
  @VisibleForTesting
  protected String getLastSubstrSplittedByDot(String hostName) {
    String[] tokens = hostName.split("\\.");
    String substrOfHostName = tokens[tokens.length - 1];
    return substrOfHostName;
  }
```

理论上讲，参数传递的正确性应该有程序员来保证，我们无需做 NULL 值或者空字符串的判断和特殊处理。调用者本不应该把 NULL 值或者空字符串传递给 getLastSubstrSplittedByDot\(\) 函数。如果传递了，那就是 code bug，需要修复。但是，话说回来，谁也保证不了程序员就一定不会传递 NULL 值或者空字符串。那我们到底该不该做 NULL 值或空字符串的判断呢？

如果函数是 private 类私有的，只在类内部被调用，完全在你自己的掌控之下，自己保证在调用这个 private 函数的时候，不要传递 NULL 值或空字符串就可以了。所以，我们可以不在 private 函数中做 NULL 值或空字符串的判断。如果函数是 public 的，你无法掌控会被谁调用以及如何调用（有可能某个同事一时疏忽，传递进了 NULL 值，这种情况也是存在的），为了尽可能提高代码的健壮性，我们最好是在 public 函数中做 NULL 值或空字符串的判断。

那你可能会说，getLastSubstrSplittedByDot\(\) 是 protected 的，既不是 private 函数，也不是 public 函数，那要不要做 NULL 值或空字符串的判断呢？之所以将它设置为 protected，是为了方便写单元测试。不过，单元测试可能要测试一些 corner case，比如输入是 NULL 值或者空字符串的情况。所以，这里我们最好也加上 NULL 值或空字符串的判断逻辑。虽然加上有些冗余，但多加些检验总归不会错的。

按照这个设计思路，我们对 getLastSubstrSplittedByDot\(\) 函数进行重构。重构之后的代码如下所示：

```java
  @VisibleForTesting
  protected String getLastSubstrSplittedByDot(String hostName) {
    if (hostName == null || hostName.isEmpty()) {
      throw IllegalArgumentException("..."); //运行时异常
    }
    String[] tokens = hostName.split("\\.");
    String substrOfHostName = tokens[tokens.length - 1];
    return substrOfHostName;
  }
```

### 重构 generateRandomAlphameric\(\) 函数

对于 generateRandomAlphameric\(int length\) 函数，如果 length &lt; 0 或 length = 0，这个函数应该返回什么？

```java
  @VisibleForTesting
  protected String generateRandomAlphameric(int length) {
    char[] randomChars = new char[length];
    int count = 0;
    Random random = new Random();
    while (count < length) {
      int maxAscii = 'z';
      int randomAscii = random.nextInt(maxAscii);
      boolean isDigit= randomAscii >= '0' && randomAscii <= '9';
      boolean isUppercase= randomAscii >= 'A' && randomAscii <= 'Z';
      boolean isLowercase= randomAscii >= 'a' && randomAscii <= 'z';
      if (isDigit|| isUppercase || isLowercase) {
        randomChars[count] = (char) (randomAscii);
        ++count;
      }
    }
    return new String(randomChars);
  }
}
```

我们先来看 length &lt; 0 的情况。生成一个长度为负值的随机字符串是不符合常规逻辑的，是一种异常行为。所以，当传入的参数 length &lt; 0 的时候，我们抛出 IllegalArgumentException 异常。

我们再来看 length = 0 的情况。length = 0 是否是异常行为呢？这就看你自己怎么定义了。我们既可以把它定义为一种异常行为，抛出 IllegalArgumentException 异常，也可以把它定义为一种正常行为，让函数在入参 length = 0 的情况下，直接返回空字符串。不管选择哪种处理方式，最关键的一点是，要在函数注释中，**明确告知 length = 0 的情况下，会返回什么样的数据。**

## ID生成器V3

```java
/**
 * Id Generator that is used to generate random IDs.
 *
 * <p>
 * The IDs generated by this class are not absolutely unique,
 * but the probability of duplication is very low.
 */
public class RandomIdGenerator implements IdGenerator {
  private static final Logger logger = LoggerFactory.getLogger(RandomIdGenerator.class);

  /**
   * Generate the random ID. The IDs may be duplicated only in extreme situation.
   *
   * @return an random ID
   */
  @Override
  public String generate() throws IdGenerationFailureException {
    String substrOfHostName = null;
    try {
      substrOfHostName = getLastFieldOfHostName();
    } catch (UnknownHostException e) {
      throw new IdGenerationFailureException("...", e);
    }
    long currentTimeMillis = System.currentTimeMillis();
    String randomString = generateRandomAlphameric(8);
    String id = String.format("%s-%d-%s",
            substrOfHostName, currentTimeMillis, randomString);
    return id;
  }

  /**
   * Get the local hostname and
   * extract the last field of the name string splitted by delimiter '.'.
   *
   * @return the last field of hostname. Returns null if hostname is not obtained.
   */
  private String getLastFieldOfHostName() throws UnknownHostException{
    String substrOfHostName = null;
    String hostName = InetAddress.getLocalHost().getHostName();
    if (hostName == null || hostName.isEmpty()) {
      throw new UnknownHostException("...");
    }
    substrOfHostName = getLastSubstrSplittedByDot(hostName);
    return substrOfHostName;
  }

  /**
   * Get the last field of {@hostName} splitted by delemiter '.'.
   *
   * @param hostName should not be null
   * @return the last field of {@hostName}. Returns empty string if {@hostName} is empty string.
   */
  @VisibleForTesting
  protected String getLastSubstrSplittedByDot(String hostName) {
    if (hostName == null || hostName.isEmpty()) {
      throw new IllegalArgumentException("...");
    }

    String[] tokens = hostName.split("\\.");
    String substrOfHostName = tokens[tokens.length - 1];
    return substrOfHostName;
  }

  /**
   * Generate random string which
   * only contains digits, uppercase letters and lowercase letters.
   *
   * @param length should not be less than 0
   * @return the random string. Returns empty string if {@length} is 0
   */
  @VisibleForTesting
  protected String generateRandomAlphameric(int length) {
    if (length <= 0) {
      throw new IllegalArgumentException("...");
    }

    char[] randomChars = new char[length];
    int count = 0;
    Random random = new Random();
    while (count < length) {
      int maxAscii = 'z';
      int randomAscii = random.nextInt(maxAscii);
      boolean isDigit= randomAscii >= '0' && randomAscii <= '9';
      boolean isUppercase= randomAscii >= 'A' && randomAscii <= 'Z';
      boolean isLowercase= randomAscii >= 'a' && randomAscii <= 'z';
      if (isDigit|| isUppercase || isLowercase) {
        randomChars[count] = (char) (randomAscii);
        ++count;
      }
    }
    return new String(randomChars);
  }
}
```

