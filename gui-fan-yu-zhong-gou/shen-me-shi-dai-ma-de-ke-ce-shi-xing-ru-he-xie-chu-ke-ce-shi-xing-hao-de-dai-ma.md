# 什么是代码的可测试性？如何写出可测试性好的代码？

## 什么是代码的可测试性？

粗略地讲，所谓代码的可测试性，就是针对代码编写单元测试的难易程度。对于一段代码，如果很难为其编写单元测试，或者单元测试写起来很费劲，需要依靠单元测试框架中很高级的特性，那往往就意味着代码设计得不够合理，代码的可测试性不好。

## 编写可测试性代码的最有效手段

依赖注入是编写可测试性代码的最有效手段。通过依赖注入，我们在编写单元测试的时候，可以通过 mock 的方法解依赖外部服务，这也是我们在编写单元测试的过程中最有技术挑战的地方。

## 常见的 Anti-Patterns

常见的测试不友好的代码有下面这 5 种：

* 代码中包含未决行为逻辑
* 滥用可变全局变量
* 滥用静态方法
* 使用复杂的继承关系
* 高度耦合的代码

## 实战

Transaction 是经过抽象简化之后的一个电商系统的交易类，用来记录每笔订单交易的情况。Transaction 类中的 execute\(\) 函数负责执行转账操作，将钱从买家的钱包转到卖家的钱包中。真正的转账操作是通过调用 WalletRpcService RPC 服务来完成的。除此之外，代码中还涉及一个分布式锁 DistributedLock 单例类，用来避免 Transaction 并发执行，导致用户的钱被重复转出。

```java
public class Transaction {
  private String id;
  private Long buyerId;
  private Long sellerId;
  private Long productId;
  private String orderId;
  private Long createTimestamp;
  private Double amount;
  private STATUS status;
  private String walletTransactionId;
  
  // ...get() methods...
  
  public Transaction(String preAssignedId, Long buyerId, Long sellerId, Long productId, String orderId) {
    if (preAssignedId != null && !preAssignedId.isEmpty()) {
      this.id = preAssignedId;
    } else {
      this.id = IdGenerator.generateTransactionId();
    }
    if (!this.id.startWith("t_")) {
      this.id = "t_" + preAssignedId;
    }
    this.buyerId = buyerId;
    this.sellerId = sellerId;
    this.productId = productId;
    this.orderId = orderId;
    this.status = STATUS.TO_BE_EXECUTD;
    this.createTimestamp = System.currentTimestamp();
  }
  
  public boolean execute() throws InvalidTransactionException {
    if ((buyerId == null || (sellerId == null || amount < 0.0) {
      throw new InvalidTransactionException(...);
    }
    if (status == STATUS.EXECUTED) return true;
    boolean isLocked = false;
    try {
      isLocked = RedisDistributedLock.getSingletonIntance().lockTransction(id);
      if (!isLocked) {
        return false; // 锁定未成功，返回false，job兜底执行
      }
      if (status == STATUS.EXECUTED) return true; // double check
      long executionInvokedTimestamp = System.currentTimestamp();
      if (executionInvokedTimestamp - createdTimestap > 14days) {
        this.status = STATUS.EXPIRED;
        return false;
      }
      WalletRpcService walletRpcService = new WalletRpcService();
      String walletTransactionId = walletRpcService.moveMoney(id, buyerId, sellerId, amount);
      if (walletTransactionId != null) {
        this.walletTransactionId = walletTransactionId;
        this.status = STATUS.EXECUTED;
        return true;
      } else {
        this.status = STATUS.FAILED;
        return false;
      }
    } finally {
      if (isLocked) {
       RedisDistributedLock.getSingletonIntance().unlockTransction(id);
      }
    }
  }
}
```

在 Transaction 类中，主要逻辑集中在 execute\(\) 函数中，所以它是我们测试的重点对象。为了尽可能全面覆盖各种正常和异常情况，针对这个函数，我设计了下面 6 个测试用例。

* 正常情况下，交易执行成功，回填用于对账（交易与钱包的交易流水）用的 walletTransactionId，交易状态设置为 EXECUTED，函数返回 true。
* buyerId、sellerId 为 null、amount 小于 0，返回 InvalidTransactionException。
* 交易已过期（createTimestamp 超过 14 天），交易状态设置为 EXPIRED，返回 false。
* 交易已经执行了（status==EXECUTED），不再重复执行转钱逻辑，返回 true。钱包（WalletRpcService）转钱失败，交易状态设置为 FAILED，函数返回 false。
* 交易正在执行着，不会被重复执行，函数直接返回 false。

测试用例设计完了。现在看起来似乎一切进展顺利。但是，事实是，当我们将测试用例落实到具体的代码实现时，你就会发现有很多行不通的地方。对于上面的测试用例，重点来看其中的 1 和 3。

现在，我们就来看测试用例 1 的代码实现。具体如下所示：

```java
public void testExecute() {
  Long buyerId = 123L;
  Long sellerId = 234L;
  Long productId = 345L;
  Long orderId = 456L;
  Transction transaction = new Transaction(null, buyerId, sellerId, productId, orderId);
  boolean executedResult = transaction.execute();
  assertTrue(executedResult);
}
```

execute\(\) 函数的执行依赖两个外部的服务，一个是 RedisDistributedLock，一个 WalletRpcService。这就导致上面的单元测试代码存在下面几个问题。

* 如果要让这个单元测试能够运行，我们需要搭建 Redis 服务和 Wallet RPC 服务。搭建和维护的成本比较高。
* 我们还需要保证将伪造的 transaction 数据发送给 Wallet RPC 服务之后，能够正确返回我们期望的结果，然而 Wallet RPC 服务有可能是第三方（另一个团队开发维护的）的服务，并不是我们可控的。换句话说，并不是我们想让它返回什么数据就返回什么。
* Transaction 的执行跟 Redis、RPC 服务通信，需要走网络，耗时可能会比较长，对单元测试本身的执行性能也会有影响。
* 网络的中断、超时、Redis、RPC 服务的不可用，都会影响单元测试的执行。

单元测试主要是测试程序员自己编写的代码逻辑的正确性，并非是端到端的集成测试，它不需要测试所依赖的外部系统（分布式锁、Wallet RPC 服务）的逻辑正确性。所以，如果代码中依赖了外部系统或者不可控组件，比如，需要依赖数据库、网络通信、文件系统等，那我们就需要将被测代码与外部系统解依赖，而这种解依赖的方法就叫作“mock”。所谓的 mock 就是用一个“假”的服务替换真正的服务。mock 的服务完全在我们的控制之下，模拟输出我们想要的数据。

那如何来 mock 服务呢？mock 的方式主要有两种，手动 mock 和利用框架 mock。利用框架 mock 仅仅是为了简化代码编写，每个框架的 mock 方式都不大一样。我们这里只展示手动 mock。



