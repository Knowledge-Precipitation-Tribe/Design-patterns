# 改善代码质量的20条编程规范

## 命名

大到项目名、模块名、包名、对外暴露的接口，小到类名、函数名、变量名、参数名，只要是做开发，我们就逃不过“起名字”这一关。命名的好坏，对于代码的可读性来说非常重要，甚至可以说是起决定性作用的。除此之外，命名能力也体现了一个程序员的基本编程素养。这也是我把“命名”放到第一个来讲解的原因。

取一个特别合适的名字是一件非常有挑战的事情，即便是对母语是英语的程序员来说，也是如此。而对于我们这些英语非母语的程序员来说，想要起一个能准确达意的名字，更是难上加难了。

实际上，命名这件事说难也不难，关键还是看你重不重视，愿不愿意花时间。对于影响范围比较大的命名，比如包名、接口、类名，我们一定要反复斟酌、推敲。实在想不到好名字的时候，可以去 GitHub 上用相关的关键词联想搜索一下，看看类似的代码是怎么命名的。

### 命名多长最合适

尽管长的命名可以包含更多的信息，更能准确直观地表达意图，但是，如果函数、变量的命名很长，那由它们组成的语句就会很长。在代码列长度有限制的情况下，就会经常出现一条语句被分割成两行的情况，这其实会影响代码可读性。

实际上，在足够表达其含义的情况下，命名当然是越短越好。但是，大部分情况下，短的命名都没有长的命名更能达意。所以，很多书籍或者文章都不推荐在命名时使用缩写。对于一些默认的、大家都比较熟知的词，我比较推荐用缩写。这样一方面能让命名短一些，另一方面又不影响阅读理解，比如，sec 表示 second、str 表示 string、num 表示 number、doc 表示 document。除此之外，对于作用域比较小的变量，我们可以使用相对短的命名，比如一些函数内的临时变量。相反，对于类名这种作用域比较大的，我更推荐用长的命名方式。

总之，命名的一个原则就是以能准确达意为目标。不过，对于代码的编写者来说，自己对代码的逻辑很清楚，总感觉用什么样的命名都可以达意，实际上，对于不熟悉你代码的同事来讲，可能就不这么认为了。所以，命名的时候，我们一定要学会换位思考，假设自己不熟悉这块代码，从代码阅读者的角度去考量命名是否足够直观。

### 利用上下文简化命名

我们先来看一个简单的例子。

```java
public class User {
  private String userName;
  private String userPassword;
  private String userAvatarUrl;
  //...
}
```

在 User 类这样一个上下文中，我们没有在成员变量的命名中重复添加“user”这样一个前缀单词，而是直接命名为 name、password、avatarUrl。在使用这些属性时候，我们能借助对象这样一个上下文，表意也足够明确。具体代码如下所示：

```java
User user = new User();
user.getName(); // 借助user对象这个上下文
```

除了类之外，函数参数也可以借助函数这个上下文来简化命名。关于这一点，我举了下面这个例子，你一看就能明白，我就不多啰嗦了。

```java
public void uploadUserAvatarImageToAliyun(String userAvatarImageUri);
//利用上下文简化为：
public void uploadUserAvatarImageToAliyun(String imageUri);
```

### 命名要可读、可搜索

什么是命名可读。先解释一下，我这里所说的“可读”，指的是不要用一些特别生僻、难发音的英文单词来命名。

命名可搜索，我们在 IDE 中编写代码的时候，经常会用“关键词联想”的方法来自动补全和搜索。比如，键入某个对象“.get”，希望 IDE 返回这个对象的所有 get 开头的方法。再比如，通过在 IDE 搜索框中输入“Array”，搜索 JDK 中数组相关的类。所以，我们在命名的时候，最好能符合整个项目的命名习惯。大家都用“selectXXX”表示查询，你就不要用“queryXXX”；大家都用“insertXXX”表示插入一条数据，你就要不用“addXXX”，统一规约是很重要的，能减少很多不必要的麻烦。

### 如何命名接口和抽象类

对于接口的命名，一般有两种比较常见的方式。一种是加前缀“I”，表示一个 Interface。比如 IUserService，对应的实现类命名为 UserService。另一种是不加前缀，比如 UserService，对应的实现类加后缀“Impl”，比如 UserServiceImpl。对于抽象类的命名，也有两种方式，一种是带上前缀“Abstract”，比如 AbstractConfiguration；另一种是不带前缀“Abstract”。实际上，对于接口和抽象类，选择哪种命名方式都是可以的，只要项目里能够统一就行。

### [Codelf](https://unbug.github.io/codelf/)

Codelf通过搜索在线开源平台Github, Bitbucket, Google Code, Codeplex, Sourceforge, Fedora Projec的项目源码，帮开发者从中找出已有的匹配关键字的变量名。这个搜索服务支持直接搜索中文。codeif支持中文查询，输入中文意思，codeif可以根据需要查询尽可能满足需要的结果，并展示与查询结果相关的支持各种编程语言的代码片段以及代码。

## 注释

命名很重要，注释跟命名同等重要。很多书籍认为，好的命名完全可以替代注释。如果需要注释，那说明命名不够好，需要在命名上下功夫，而不是添加注释。实际上，我个人觉得，这样的观点有点太过极端。命名再好，毕竟有长度限制，不可能足够详尽，而这个时候，注释就是一个很好的补充。

### 注释到底该写什么？

注释的目的就是让代码更容易看懂。只要符合这个要求的内容，你就可以将它写到注释里。总结一下，注释的内容主要包含这样三个方面：做什么、为什么、怎么做。我来举一个例子给你具体解释一下。

```java
/**
* (what) Bean factory to create beans. 
* 
* (why) The class likes Spring IOC framework, but is more lightweight. 
*
* (how) Create objects from different sources sequentially:
* user specified object > SPI > configuration > default object.
*/
public class BeansFactory {
  // ...
}
```

有些人认为，注释是要提供一些代码没有的额外信息，所以不要写“做什么、怎么做”，这两方面在代码中都可以体现出来，只需要写清楚“为什么”，表明代码的设计意图即可。我个人不是特别认可这样的观点，理由主要有下面 3 点。

* 注释比代码承载的信息更多

命名的主要目的是解释“做什么”。比如，void increaseWalletAvailableBalance\(BigDecimal amount\) 表明这个函数用来增加钱包的可用余额，boolean isValidatedPassword 表明这个变量用来标识是否是合法密码。函数和变量如果命名得好，确实可以不用再在注释中解释它是做什么的。但是，对于类来说，包含的信息比较多，一个简单的命名就不够全面详尽了。这个时候，在注释中写明“做什么”就合情合理了。

* 注释起到总结性作用、文档的作用

代码之下无秘密。阅读代码可以明确地知道代码是“怎么做”的，也就是知道代码是如何实现的，那注释中是不是就不用写“怎么做”了？实际上也可以写。在注释中，关于具体的代码实现思路，我们可以写一些总结性的说明、特殊情况的说明。这样能够让阅读代码的人通过注释就能大概了解代码的实现思路，阅读起来就会更加容易。实际上，对于有些比较复杂的类或者接口，我们可能还需要在注释中写清楚“如何用”，举一些简单的 quick start 的例子，让使用者在不阅读代码的情况下，快速地知道该如何使用。

* 一些总结性注释能让代码结构更清晰

对于逻辑比较复杂的代码或者比较长的函数，如果不好提炼、不好拆分成小的函数调用，那我们可以借助总结性的注释来让代码结构更清晰、更有条理。

```java
public boolean isValidPasword(String password) {
  // check if password is null or empty
  if (StringUtils.isBlank(password)) {
    return false;
  }

  // check if the length of password is between 4 and 64
  int length = password.length();
  if (length < 4 || length > 64) {
    return false;
  }
    
  // check if password contains only a~z,0~9,dot
  for (int i = 0; i < length; ++i) {
    char c = password.charAt(i);
    if (!((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '.')) {
      return false;
    }
  }
  return true;
}
```

### 注释是不是越多越好？

注释太多和太少都有问题。太多，有可能意味着代码写得不够可读，需要写很多注释来补充。除此之外，注释太多也会对代码本身的阅读起到干扰。而且，后期的维护成本也比较高，有时候代码改了，注释忘了同步修改，就会让代码阅读者更加迷惑。当然，如果代码中一行注释都没有，那只能说明这个程序员很懒，我们要适当督促一下，让他注意添加一些必要的注释。

按照我的经验来说，类和函数一定要写注释，而且要写得尽可能全面、详细，而函数内部的注释要相对少一些，一般都是靠好的命名、提炼函数、解释性变量、总结性注释来提高代码的可读性。

## 类、函数多大才合适？

总体上来讲，类或函数的代码行数不能太多，但也不能太少。类或函数的代码行数太多，一个类上千行，一个函数几百行，逻辑过于繁杂，阅读代码的时候，很容易就会看了后面忘了前面。相反，类或函数的代码行数太少，在代码总量相同的情况下，被分割成的类和函数就会相应增多，调用关系就会变得更复杂，阅读某个代码逻辑的时候，需要频繁地在 n 多类或者 n 多函数之间跳来跳去，阅读体验也不好。

对于函数代码行数的最大限制，网上有一种说法，那就是不要超过一个显示屏的垂直高度。比如，在我的电脑上，如果要让一个函数的代码完整地显示在 IDE 中，那最大代码行数不能超过 50。这个说法我觉得挺有道理的。因为超过一屏之后，在阅读代码的时候，为了串联前后的代码逻辑，就可能需要频繁地上下滚动屏幕，阅读体验不好不说，还容易出错。

对于类的代码行数的最大限制，这个就更难给出一个确切的值了。当一个类的代码读起来让你感觉头大了，实现某个功能时不知道该用哪个函数了，想用哪个函数翻半天都找不到了，只用到一个小功能要引入整个类（类中包含很多无关此功能实现的函数）的时候，这就说明类的行数过多了。

## 一行代码多长最合适？

在[Google Java Style Guide](https://google.github.io/styleguide/javaguide.html)文档中，一行代码最长限制为 100 个字符。不过，不同的编程语言、不同的规范、不同的项目团队，对此的限制可能都不相同。不管这个限制是多少，总体上来讲我们要遵循的一个原则是：一行代码最长不能超过 IDE 显示的宽度。需要滚动鼠标才能查看一行的全部代码，显然不利于代码的阅读。当然，这个限制也不能太小，太小会导致很多稍长点的语句被折成两行，也会影响到代码的整洁，不利于阅读。

## 善用空行分割单元块

对于比较长的函数，如果逻辑上可以分为几个独立的代码块，在不方便将这些独立的代码块抽取成小函数的情况下，为了让逻辑更加清晰，除了上一节课中提到的用总结性注释的方法之外，我们还可以使用空行来分割各个代码块。

除此之外，在类的成员变量与函数之间、静态成员变量与普通成员变量之间、各函数之间、甚至各成员变量之间，我们都可以通过添加空行的方式，让这些不同模块的代码之间，界限更加明确。写代码就类似写文章，善于应用空行，可以让代码的整体结构看起来更加有清晰、有条理。

## 四格缩进还是两格缩进？

“PHP 是世界上最好的编程语言？代码换行应该四格缩进还是两格缩进？”这应该是程序员争论得最多的两个话题了。据我所知，Java 语言倾向于两格缩进，PHP 语言倾向于四格缩进。至于到底应该是两格缩进还是四格缩进，我觉得这个取决于个人喜好。只要项目内部能够统一就行了。

当然，还有一个选择的标准，那就是跟业内推荐的风格统一、跟著名开源项目统一。当我们需要拷贝一些开源的代码到项目里的时候，能够让引入的代码跟我们项目本身的代码，保持风格统一。

除此之外，值得强调的是，不管是用两格缩进还是四格缩进，一定不要用 tab 键缩进。因为在不同的 IDE 下，tab 键的显示宽度不同，有的显示为四格缩进，有的显示为两格缩进。如果在同一个项目中，不同的同事使用不同的缩进方式（空格缩进或 tab 键缩进），有可能会导致有的代码显示为两格缩进、有的代码显示为四格缩进。

## 大括号是否要另起一行？

左大括号是否要另起一行呢？这个也有争论。具体代码示例如下所示：

```java
// PHP
class ClassName
{
    public function foo()
    {
        // method body
    }
}

// Java
public class ClassName {
  public void foo() {
    // method body
  }
}
```

我个人还是比较推荐，将括号放到跟语句同一行的风格。理由跟上面类似，节省代码行数。但是将大括号另起新的一行的方式，也有它的优势。这样的话，左右括号可以垂直对齐，哪些代码属于哪一个代码块，更一目了然。

不过，还是那句话，大括号跟上一条语句在同一行，还是另起新的一行，只要团队统一、业内统一、跟开源项目看齐就好了，没有绝对的优劣之分。

## 类中成员的排列顺序

在 Java 类文件中，先要书写类所属的包名，然后再罗列 import 引入的依赖类。在 Google 编码规范中，依赖类按照字母序从小到大排列。

在类中，成员变量排在函数的前面。成员变量之间或函数之间，都是按照“先静态（静态函数或静态成员变量）、后普通（非静态函数或非静态成员变量）”的方式来排列的。除此之外，成员变量之间或函数之间，还会按照作用域范围从大到小的顺序来排列，先写 public 成员变量或函数，然后是 protected 的，最后是 private 的。

不过，不同的编程语言中，类内部成员的排列顺序可能会有比较大的差别。比如 C++ 中，成员变量会习惯性放到函数后面。除此之外，函数之间的排列顺序，会按照刚刚我们提到的作用域的大小来排列。实际上，还有另外一种排列习惯，那就是把有调用关系的函数放到一块。比如，一个 public 函数调用了另外一个 private 函数，那就把这两者放到一块。

## 把代码分割成更小的单元块

大部分人阅读代码的习惯都是，先看整体再看细节。所以，我们要有模块化和抽象思维，善于将大块的复杂逻辑提炼成类或者函数，屏蔽掉细节，让阅读代码的人不至于迷失在细节中，这样能极大地提高代码的可读性。不过，只有代码逻辑比较复杂的时候，我们其实才建议提炼类或者函数。毕竟如果提炼出的函数只包含两三行代码，在阅读代码的时候，还得跳过去看一下，这样反倒增加了阅读成本。

这里我举一个例子来进一步解释一下。代码具体如下所示。重构前，在 invest\(\) 函数中，最开始的那段关于时间处理的代码，是不是很难看懂？重构之后，我们将这部分逻辑抽象成一个函数，并且命名为 isLastDayOfMonth，从名字就能清晰地了解它的功能，判断今天是不是当月的最后一天。这里，我们就是通过将复杂的逻辑代码提炼成函数，大大提高了代码的可读性。

```java
// 重构前的代码
public void invest(long userId, long financialProductId) {
  Calendar calendar = Calendar.getInstance();
  calendar.setTime(date);
  calendar.set(Calendar.DATE, (calendar.get(Calendar.DATE) + 1));
  if (calendar.get(Calendar.DAY_OF_MONTH) == 1) {
    return;
  }
  //...
}

// 重构后的代码：提炼函数之后逻辑更加清晰
public void invest(long userId, long financialProductId) {
  if (isLastDayOfMonth(new Date())) {
    return;
  }
  //...
}

public boolean isLastDayOfMonth(Date date) {
  Calendar calendar = Calendar.getInstance();
  calendar.setTime(date);
  calendar.set(Calendar.DATE, (calendar.get(Calendar.DATE) + 1));
  if (calendar.get(Calendar.DAY_OF_MONTH) == 1) {
   return true;
  }
  return false;
}
```

## 避免函数参数过多

我个人觉得，函数包含 3、4 个参数的时候还是能接受的，大于等于 5 个的时候，我们就觉得参数有点过多了，会影响到代码的可读性，使用起来也不方便。针对参数过多的情况，一般有 2 种处理方法。

* 考虑函数是否职责单一，是否能通过拆分成多个函数的方式来减少参数。示例代码如下所示：

```java
public User getUser(String username, String telephone, String email);

// 拆分成多个函数
public User getUserByUsername(String username);
public User getUserByTelephone(String telephone);
public User getUserByEmail(String email);
```

* 将函数的参数封装成对象。示例代码如下所示：

```java
public void postBlog(String title, String summary, String keywords, String content, String category, long authorId);

// 将参数封装成对象
public class Blog {
  private String title;
  private String summary;
  private String keywords;
  private Strint content;
  private String category;
  private long authorId;
}
public void postBlog(Blog blog);
```

除此之外，如果函数是对外暴露的远程接口，将参数封装成对象，还可以提高接口的兼容性。在往接口中添加新的参数的时候，老的远程接口调用者有可能就不需要修改代码来兼容新的接口了。

## 勿用函数参数来控制逻辑

不要在函数中使用布尔类型的标识参数来控制内部逻辑，true 的时候走这块逻辑，false 的时候走另一块逻辑。这明显违背了单一职责原则和接口隔离原则。我建议将其拆成两个函数，可读性上也要更好。我举个例子来说明一下。

```java
public void buyCourse(long userId, long courseId, boolean isVip);

// 将其拆分成两个函数
public void buyCourse(long userId, long courseId);
public void buyCourseForVip(long userId, long courseId);
```

不过，如果函数是 private 私有函数，影响范围有限，或者拆分之后的两个函数经常同时被调用，我们可以酌情考虑保留标识参数。示例代码如下所示：

```java
// 拆分成两个函数的调用方式
boolean isVip = false;
//...省略其他逻辑...
if (isVip) {
  buyCourseForVip(userId, courseId);
} else {
  buyCourse(userId, courseId);
}

// 保留标识参数的调用方式更加简洁
boolean isVip = false;
//...省略其他逻辑...
buyCourse(userId, courseId, isVip);
```

除了布尔类型作为标识参数来控制逻辑的情况外，还有一种“根据参数是否为 null”来控制逻辑的情况。针对这种情况，我们也应该将其拆分成多个函数。拆分之后的函数职责更明确，不容易用错。具体代码示例如下所示：

```java
public List<Transaction> selectTransactions(Long userId, Date startDate, Date endDate) {
  if (startDate != null && endDate != null) {
    // 查询两个时间区间的transactions
  }
  if (startDate != null && endDate == null) {
    // 查询startDate之后的所有transactions
  }
  if (startDate == null && endDate != null) {
    // 查询endDate之前的所有transactions
  }
  if (startDate == null && endDate == null) {
    // 查询所有的transactions
  }
}

// 拆分成多个public函数，更加清晰、易用
public List<Transaction> selectTransactionsBetween(Long userId, Date startDate, Date endDate) {
  return selectTransactions(userId, startDate, endDate);
}

public List<Transaction> selectTransactionsStartWith(Long userId, Date startDate) {
  return selectTransactions(userId, startDate, null);
}

public List<Transaction> selectTransactionsEndWith(Long userId, Date endDate) {
  return selectTransactions(userId, null, endDate);
}

public List<Transaction> selectAllTransactions(Long userId) {
  return selectTransactions(userId, null, null);
}

private List<Transaction> selectTransactions(Long userId, Date startDate, Date endDate) {
  // ...
}
```

## 函数设计要职责单一

我们在前面讲到单一职责原则的时候，针对的是类、模块这样的应用对象。实际上，对于函数的设计来说，更要满足单一职责原则。相对于类和模块，函数的粒度比较小，代码行数少，所以在应用单一职责原则的时候，没有像应用到类或者模块那样模棱两可，能多单一就多单一。

具体的代码示例如下所示：

```java
public boolean checkUserIfExisting(String telephone, String username, String email)  { 
  if (!StringUtils.isBlank(telephone)) {
    User user = userRepo.selectUserByTelephone(telephone);
    return user != null;
  }
  
  if (!StringUtils.isBlank(username)) {
    User user = userRepo.selectUserByUsername(username);
    return user != null;
  }
  
  if (!StringUtils.isBlank(email)) {
    User user = userRepo.selectUserByEmail(email);
    return user != null;
  }
  
  return false;
}

// 拆分成三个函数
public boolean checkUserIfExistingByTelephone(String telephone);
public boolean checkUserIfExistingByUsername(String username);
public boolean checkUserIfExistingByEmail(String email);
```

## 移除过深的嵌套层次

代码嵌套层次过深往往是因为 if-else、switch-case、for 循环过度嵌套导致的。我个人建议，嵌套最好不超过两层，超过两层之后就要思考一下是否可以减少嵌套。过深的嵌套本身理解起来就比较费劲，除此之外，嵌套过深很容易因为代码多次缩进，导致嵌套内部的语句超过一行的长度而折成两行，影响代码的整洁。

解决嵌套过深的方法也比较成熟，有下面 4 种常见的思路。

* 去掉多余的 if 或 else 语句。代码示例如下所示：

```java
// 示例一
public double caculateTotalAmount(List<Order> orders) {
  if (orders == null || orders.isEmpty()) {
    return 0.0;
  } else { // 此处的else可以去掉
    double amount = 0.0;
    for (Order order : orders) {
      if (order != null) {
        amount += (order.getCount() * order.getPrice());
      }
    }
    return amount;
  }
}

// 示例二
public List<String> matchStrings(List<String> strList,String substr) {
  List<String> matchedStrings = new ArrayList<>();
  if (strList != null && substr != null) {
    for (String str : strList) {
      if (str != null) { // 跟下面的if语句可以合并在一起
        if (str.contains(substr)) {
          matchedStrings.add(str);
        }
      }
    }
  }
  return matchedStrings;
}
```

* 使用编程语言提供的 continue、break、return 关键字，提前退出嵌套。代码示例如下所示：

```java
// 重构前的代码
public List<String> matchStrings(List<String> strList,String substr) {
  List<String> matchedStrings = new ArrayList<>();
  if (strList != null && substr != null){ 
    for (String str : strList) {
      if (str != null && str.contains(substr)) {
        matchedStrings.add(str);
        // 此处还有10行代码...
      }
    }
  }
  return matchedStrings;
}

// 重构后的代码：使用continue提前退出
public List<String> matchStrings(List<String> strList,String substr) {
  List<String> matchedStrings = new ArrayList<>();
  if (strList != null && substr != null){ 
    for (String str : strList) {
      if (str == null || !str.contains(substr)) {
        continue; 
      }
      matchedStrings.add(str);
      // 此处还有10行代码...
    }
  }
  return matchedStrings;
}
```

* 调整执行顺序来减少嵌套。具体的代码示例如下所示：

```java
// 重构前的代码
public List<String> matchStrings(List<String> strList,String substr) {
  List<String> matchedStrings = new ArrayList<>();
  if (strList != null && substr != null) {
    for (String str : strList) {
      if (str != null) {
        if (str.contains(substr)) {
          matchedStrings.add(str);
        }
      }
    }
  }
  return matchedStrings;
}

// 重构后的代码：先执行判空逻辑，再执行正常逻辑
public List<String> matchStrings(List<String> strList,String substr) {
  if (strList == null || substr == null) { //先判空
    return Collections.emptyList();
  }

  List<String> matchedStrings = new ArrayList<>();
  for (String str : strList) {
    if (str != null) {
      if (str.contains(substr)) {
        matchedStrings.add(str);
      }
    }
  }
  return matchedStrings;
}
```

* 将部分嵌套逻辑封装成函数调用，以此来减少嵌套。具体的代码示例如下所示：

```java
// 重构前的代码
public List<String> appendSalts(List<String> passwords) {
  if (passwords == null || passwords.isEmpty()) {
    return Collections.emptyList();
  }
  
  List<String> passwordsWithSalt = new ArrayList<>();
  for (String password : passwords) {
    if (password == null) {
      continue;
    }
    if (password.length() < 8) {
      // ...
    } else {
      // ...
    }
  }
  return passwordsWithSalt;
}

// 重构后的代码：将部分逻辑抽成函数
public List<String> appendSalts(List<String> passwords) {
  if (passwords == null || passwords.isEmpty()) {
    return Collections.emptyList();
  }

  List<String> passwordsWithSalt = new ArrayList<>();
  for (String password : passwords) {
    if (password == null) {
      continue;
    }
    passwordsWithSalt.add(appendSalt(password));
  }
  return passwordsWithSalt;
}

private String appendSalt(String password) {
  String passwordWithSalt = password;
  if (password.length() < 8) {
    // ...
  } else {
    // ...
  }
  return passwordWithSalt;
}
```

除此之外，常用的还有通过使用多态来替代 if-else、switch-case 条件判断的方法。

## 学会使用解释性变量

常用的用解释性变量来提高代码的可读性的情况有下面 2 种。

* 常量取代魔法数字。示例代码如下所示：

```java
public double CalculateCircularArea(double radius) {
  return (3.1415) * radius * radius;
}

// 常量替代魔法数字
public static final Double PI = 3.1415;
public double CalculateCircularArea(double radius) {
  return PI * radius * radius;
}
```

* 使用解释性变量来解释复杂表达式。示例代码如下所示：

```java
if (date.after(SUMMER_START) && date.before(SUMMER_END)) {
  // ...
} else {
  // ...
}

// 引入解释性变量后逻辑更加清晰
boolean isSummer = date.after(SUMMER_START)&&date.before(SUMMER_END);
if (isSummer) {
  // ...
} else {
  // ...
} 
```

