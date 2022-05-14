---
title: DomainEvents
layout: post
category: architect
author: 夏泽民
---
Spring Data之@DomainEvents注解
背景
在对一个Entity进行save操作时，往往需要触发后续的业务流程，通常采用如下做法

public void saveUser(){
	User user = ...
	user = repository.save(user);

	doSomething(user);
}

public void action(){
	User user = ...
	saveUser(user);
	doSomething(user);
}

其中有一些注意事项，例如

doSomething与saveUser在同一个事务中，需要考虑doSomething中的异常对repository.save(user)的影响
doSomething与saveUser不在同一个事务中，那么在doSomething中查询user时将查询不到，因为saveUser的事务还未提交。
这种情况则需要将doSomething上移到调用saveUser同级的地方调用这种情况则需要将doSomething上移到调用saveUser同级的地方调用
DomainEvents
近日在Spring Data的官方手册中看到@DomainEvents的介绍。官方解释是由Repositoty管理的Entity是源于聚合根（ aggregate roots）的，在领域驱动设计系统中，可以通过聚合根发出领域事件。在Spring Data中可以通过@DomainEvents注解在聚合根的方法上，从而可以简单快捷的发出事件。下面就来看一下，DomainEvents的具体使用效果。
首先定义一个普通的Entity

@Data
@Entity
@Table(name = "t_user")
@AllArgsConstructor
@NoArgsConstructor
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String firstName;

    private String lastName;

    private Integer age;

	//该方法会在userRepository.save()调用时被触发调用
    @DomainEvents
    Collection<UserSaveEvent> domainEvents() {
        return Arrays.asList(new UserSaveEvent(this.id));
    }

}

其中UserSaveEvent的定义如下

@Data
@AllArgsConstructor
public class UserSaveEvent {

    private Long id;

}
1
2
3
4
5
6
7
再定义一个UserService消费发出的事件

@Service
public class UserService {

    @Autowired
    private UserRepository userRepository;

	//接受User发出的类型为UserSaveEvent的DomainEvents事件
    @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
    public void event(UserSaveEvent event){
        System.out.println(userRepository.getOne(event.getId()));
    }

}

其中@TransactionalEventListener注解的phase有多个选项

BEFORE_COMMIT
AFTER_COMMIT
AFTER_ROLLBACK
AFTER_COMPLETION
看名字就知道它们的作用和区别了，因为事件是repository.save发出的，这里就涉及到了事务。通过phase的不同选项，就能选择是在事务提交前获取事件，还是提交后，或者混滚的时候。

运行一下单元测试

@Before
public void before(){

    userRepository.saveAll(Arrays.asList(
            new User(null,"刘","一", 20),
            new User(null,"陈","二", 20),
            new User(null,"张","三", 20),
            new User(null,"李","四", 20),
            new User(null,"王","五", 20),
            new User(null,"赵","六", 20),
            new User(null,"孙","七", 20),
            new User(null,"周","八", 20)
    ));
}

控制台输出

Hibernate: insert into t_user (id, age, first_name, last_name) values (null, ?, ?, ?)
Hibernate: insert into t_user (id, age, first_name, last_name) values (null, ?, ?, ?)
Hibernate: insert into t_user (id, age, first_name, last_name) values (null, ?, ?, ?)
Hibernate: insert into t_user (id, age, first_name, last_name) values (null, ?, ?, ?)
User(id=1, firstName=刘, lastName=一, age=20)
User(id=2, firstName=陈, lastName=二, age=20)
User(id=3, firstName=张, lastName=三, age=20)
User(id=4, firstName=李, lastName=四, age=20)

上面是使用的phase = TransactionPhase.AFTER_COMMIT，即事务提交后响应事件，所以userRepository.getOne(event.getId())能查询到user对象。如果改成TransactionPhase.BEFORE_COMMIT呢

@TransactionalEventListener(phase = TransactionPhase.BEFORE_COMMIT)
public void event(UserSaveEvent event){
    System.out.println(userRepository.getOne(event.getId()));
}

其实效果是一样的也能查询到user，难道BEFORE_COMMIT没起作用？没提交事务前按理是查询不到的才对。
其实是因为session的缓存，因为event方法并没有添加@Async注解异步，也没有@Transactional(value = Transactional.TxType.REQUIRES_NEW)开启新事务，所以这时与发送事件的repository.save还在一个事务内。

如果给event方法开启新事务

@TransactionalEventListener(phase = TransactionPhase.BEFORE_COMMIT)
@Transactional(value = Transactional.TxType.REQUIRES_NEW)
public void event(UserSaveEvent event){
    System.out.println(userRepository.getOne(event.getId()));
}

这样查询就会报错，因为查不到了

org.springframework.orm.jpa.JpaObjectRetrievalFailureException: Unable to find com.learn.data.entity.User with id 1; nested exception is javax.persistence.EntityNotFoundException: Unable to find com.learn.data.entity.User with id 1
1
再将phase改成TransactionPhase.AFTER_COMMIT试试

@TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
@Transactional(value = Transactional.TxType.REQUIRES_NEW)
public void event(UserSaveEvent event){
    System.out.println(userRepository.getOne(event.getId()));
}

控制输出

Hibernate: insert into t_user (id, age, first_name, last_name) values (null, ?, ?, ?)
Hibernate: insert into t_user (id, age, first_name, last_name) values (null, ?, ?, ?)
Hibernate: insert into t_user (id, age, first_name, last_name) values (null, ?, ?, ?)
Hibernate: insert into t_user (id, age, first_name, last_name) values (null, ?, ?, ?)
Hibernate: select user0_.id as id1_0_0_, user0_.age as age2_0_0_, user0_.first_name as first_na3_0_0_, user0_.last_name as last_nam4_0_0_ from t_user user0_ where user0_.id=?
User(id=1, firstName=刘, lastName=一, age=20)
Hibernate: select user0_.id as id1_0_0_, user0_.age as age2_0_0_, user0_.first_name as first_na3_0_0_, user0_.last_name as last_nam4_0_0_ from t_user user0_ where user0_.id=?
User(id=2, firstName=陈, lastName=二, age=20)
Hibernate: select user0_.id as id1_0_0_, user0_.age as age2_0_0_, user0_.first_name as first_na3_0_0_, user0_.last_name as last_nam4_0_0_ from t_user user0_ where user0_.id=?
User(id=3, firstName=张, lastName=三, age=20)
Hibernate: select user0_.id as id1_0_0_, user0_.age as age2_0_0_, user0_.first_name as first_na3_0_0_, user0_.last_name as last_nam4_0_0_ from t_user user0_ where user0_.id=?
User(id=4, firstName=李, lastName=四, age=20)

现在能查询到了，但控制台里面的查询打出了select语句，与没有添加@Transactional时是不一样了，没有@Transactional注解时是没有select语句的，说明JPA查询的是seesion缓存并没有真正执行查询。

结束
@DomainEvents和@TransactionalEventListener的组合使用，给我们处理实体保存后触发事件。特别是异步事件（给event方法加上@Async，同时开启@EnableAsync）是非常简便的，它是一种领域驱动的思想，让代码显得更加的内聚。
<!-- more -->
https://www.cnblogs.com/daxnet/archive/2012/12/27/2836372.html

在最近的一次代码签入中，Byteart Retail已经可以支持领域事件（Domain Events）的定义和处理了。在这篇文章中，我将详细介绍领域事件机制在Byteart Retail案例中的具体实现。

在进行领域建模的时候，我们就已经知道保证领域模型纯净度的必要性。简而言之，领域模型中的各个对象都应该是POCO（POJO）对象，而不应向其添加任何与技术架构相关的内容。Udi Dahan曾经说过：“The main assertion being that you do *not* need to inject anything into your domain entities. Not services. Not repositories. Nothing.”。因此，在之前有朋友提出过，是否可以在Domain Model中访问仓储？现在看来，答案是否定的。那么Domain Service呢？当然也不行。顺便提一下，在当前版本的Byteart Retail中的Domain Service访问了仓储，这是一个不太合理的做法，在下个版本中我将进行改进。那么，如果在某些业务需求下，需要访问这些技术层面的东西，又该怎么办呢？比如当系统管理员完成销售订单的发货操作时，希望向客户发送一份电子邮件。此时就要用到领域事件。

领域事件是应用系统中众多事件的一种分类。企业级应用程序事件大致可以分为三类：系统事件、应用事件和领域事件。领域事件的触发点在领域模型（Domain Model）中，故以此得名。通过使用领域事件，我们可以实现领域模型对象状态的异步更新、外部系统接口的委托调用，以及通过事件派发机制实现系统集成。在进行实际业务分析的过程中，如果在通用语言中存在“当a发生时，我们就需要做到b。”这样的描述，则表明a可以定义成一个领域事件。领域事件的命名一般也就是“产生事件的对象名称+完成的动作的过去式”的形式，比如：订单已经发货的事件（OrderDispatchedEvent）、订单已被收货和确认的事件（OrderConfirmedEvent）等。在当前的Byteart Retail案例的源代码中，就引入了这两种领域事件。事实上针对该案例而言，还有很多地方可以使用领域事件，比如当客户地址变更时，可以通过事件处理器来更新所有该事件发生前所有未发货订单的客户收货地址等。当然，为了简单起见，案例仅演示了上述两种事件。

另外，领域事件本身具有自描述性。它不仅能够表述系统发生了什么事情，而且还能够描述发生事件的动机。例如AddressChangedEvent可以衍生出两个派生类：ContactMovedEvent和AddressCorrectedEvent，虽然这两种事件都会导致地址信息的变更，但它们所表述的动机是不同的：前者体现了地址变更是因为联系人的地址发生了改变，而后者则体现了地址变更是因为地址信息原本是错的，现在被更正过来了。

现在，我们开始逐步讨论领域事件在Byteart Retail案例中的实现方式。

定义一个领域事件
通常，我们会为领域事件定义一个接口（IDomainEvent接口），所有实现了该接口的类型都被认为是一个领域事件的类型。为了能够向事件处理器等事件管理机构提供完善的信息，我们可以在这个接口中设置一些属性，比如事件发生的时间戳、事件来源以及事件的ID值等等，当然这些内容都是根据具体的项目需求而定的。在Byteart Retail案例中，又定义了一个抽象类（DomainEvent类），该类实现了IDomainEvent接口，同时在这个类中提供了一个带参构造函数，它接受一个代表事件来源（Event Source）的领域实体作为参数，因此，在整个Byteart Retail中约定，所有领域事件类型都继承于DomainEvent类型，以便强制每个类型都需要提供一个相同参数类型的带参构造函数。这样做的好处是，每当开发人员初始化一个领域事件，都必须设置其产生的事件来源，在开发上达成了一种契约，有效地降低了错误的产生。

比如，上文所提到的OrderDispatchedEvent定义如下：

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
/// <summary>
/// 表示当针对某销售订单进行发货时所产生的领域事件。
/// </summary>
public class OrderDispatchedEvent : DomainEvent
{
    #region Ctor
    /// <summary>
    /// 初始化一个新的<c>OrderDispatchedEvent</c>类型的实例。
    /// </summary>
    /// <param name="source">产生领域事件的事件源对象。</param>
    public OrderDispatchedEvent(IEntity source) : base(source) { }
    #endregion
 
    #region Public Properties
    /// <summary>
    /// 获取或设置订单发货的日期。
    /// </summary>
    public DateTime DispatchedDate { get; set; }
    #endregion
}
在这个事件定义中，构造函数接受一个IEntity类型的参数，以表示产生当前事件的实体对象，此外，它还包含了订单发货的日期信息。

领域事件的派发和处理
处理领域事件的机制称为“事件处理器（Event Handler）”，而领域事件的派发，我们则是通过“事件聚合器（Event Aggregator）”实现的。接下来，我们讨论这两个部分的具体实现过程。

事件处理器（Event Handler）
事件处理器的任务是处理捕获的事件，它的职责是相对单一的：只需要对传入的信息进行处理即可。因此，在实现上我们可以将其定义为一个泛型接口，例如在Byteart Retail中，它被定义为IDomainEventHandler<TDomainEvent>接口，TDomainEvent类型参数指定了事件处理器所能够处理的领域事件的类型。一般情况下，该接口只提供一个Handle方法，该方法接受一个类型为TDomainEvent的对象（即领域事件实例）作为参数。所有实现了该接口的类型都被认为是能够处理特定类型领域事件的事件处理器。与领域事件的设计相同，在Byteart Retail中，还提供了一个名为DomainEventHandler<TDomainEvent>的泛型抽象类，该类直接实现了IDomainEventHandler<TDomainEvent>接口，同时实现了一个异步事件处理的方法：HandleAsync。同理，为了达成开发规范，在Byteart Retail中，所有领域事件处理器都应该继承于DomainEventHandler<TDomainEvent>抽象类，并实现其中的抽象方法：Handle方法。由于模板方法模式的支持，开发人员无需考虑异步事件处理的实现（即HandleAsync方法会创建一个用于异步任务处理的Task对象，来执行Handle方法所定义的操作）。

此外，为了简化编程模型，Byteart Retail还支持基于委托的事件处理器。这个设计其实并不是必须的，但在Byteart Retail中，为了简化事件订阅的操作，还是引入了这样一种基于委托的事件处理器。在某些情况下，事件处理逻辑会比较简单，比如仅仅是在捕获到某个事件时更新领域对象的状态，那么对于这样一些应用场景，开发人员就无需为每一个相对简单的事件处理逻辑定义一个单独的事件处理器类型，而只需要让委托的匿名方法来订阅和处理事件即可，这样做不仅简洁而且便于单体测试。有关事件处理器如何去订阅领域事件，我们将在下一小节“事件聚合器”中讨论。还是先让我们来看看Byteart Retail中是如何实现这种基于委托的事件处理器的。

在Byteart Retail中，有一个特殊的领域事件处理器，它与其它领域事件处理器一样，也继承于DomainEventHandler<TDomainEvent>泛型抽象类，但它的特殊性在于，它会在构造函数中接受一个Action<TDomainEvent>类型的委托作为参数，于是，通过一种类似装饰器模式的方式，将Action<TDomainEvent>委托“装饰”成DomainEventHandler<TDomainEvent>类型的对象：

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
/// <summary>
/// 表示代理给定的领域事件处理委托的领域事件处理器。
/// </summary>
/// <typeparam name="TEvent"></typeparam>
internal sealed class ActionDelegatedDomainEventHandler<TEvent> : DomainEventHandler<TEvent>
    where TEvent : class, IDomainEvent
{
    #region Private Fields
    private readonly Action<TEvent> eventHandlerDelegate;
    #endregion
 
    #region Ctor
    /// <summary>
    /// 初始化一个新的<c>ActionDelegatedDomainEventHandler{TEvent}</c>实例。
    /// </summary>
    /// <param name="eventHandlerDelegate">用于当前领域事件处理器所代理的事件处理委托。</param>
    public ActionDelegatedDomainEventHandler(Action<TEvent> eventHandlerDelegate)
    {
        this.eventHandlerDelegate = eventHandlerDelegate;
    }
    #endregion
     
    // 其它函数和属性暂时忽略
}
在此类中Handle方法的实现就非常简单了：

1
2
3
4
5
6
7
8
/// <summary>
/// 处理给定的事件。
/// </summary>
/// <param name="evnt">需要处理的事件。</param>
public override void Handle(TEvent evnt)
{
    this.eventHandlerDelegate(evnt);
}
这种做法的优点是，可以将基于委托的事件处理器当成是普通的事件处理器类型，从而统一了事件订阅和事件派发的接口定义。

需要注意的是，对于ActionDelegatedDomainEventHandler而言，实例之间的相等性并不是由实例本身决定的，而是由其所代理的委托决定的，这对于事件处理器对事件的订阅，以及事件聚合器对事件的派发，都有着重要的影响。根据这个分析，我们就需要重载Equals方法，使用Delegate.Equals方法来判定两个委托的相等性。在Byteart Retail中，IDomainEventHandler<TDomainEvent>接口还实现了IEquatable接口，因此，只需要重载IEquatable接口中定义的Equals方法即可：

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
/// <summary>
/// 获取一个<see cref="Boolean"/>值，该值表示当前对象是否与给定的类型相同的另一对象相等。
/// </summary>
/// <param name="other">需要比较的与当前对象类型相同的另一对象。</param>
/// <returns>如果两者相等，则返回true，否则返回false。</returns>
public override bool Equals(IDomainEventHandler<TEvent> other)
{
    if (ReferenceEquals(this, other))
        return true;
    if ((object)other == (object)null)
        return false;
    ActionDelegatedDomainEventHandler<TEvent> otherDelegate = 
        other as ActionDelegatedDomainEventHandler<TEvent>;
    if ((object)otherDelegate == (object)null)
        return false;
    // 使用Delegate.Equals方法判定两个委托是否是代理的同一方法。
    return Delegate.Equals(this.eventHandlerDelegate, otherDelegate.eventHandlerDelegate);
}
现在我们已经定义好了事件处理器接口以及相关的类，同时也根据需要实现了几个简单的事件处理器（具体代码请参考Byteart Retail案例中ByteartRetail.Domain.Events.Handlers命名空间下的类）。接下来我们要让领域模型能够在业务需要的地方触发领域事件，并让这些事件处理器能够对获得的事件进行处理。在Byteart Retail案例中，这部分内容是使用“事件聚合器”实现的。

事件聚合器（Event Aggregator）
事件聚合器是一种企业应用架构模式，其作用主要是聚合领域模型中的事件处理器，以便事件在触发的时候，被聚合的事件处理器能够对事件进行处理。在Byteart Retail中，事件聚合器的结构如下：

image

在这个设计中，事件聚合器提供了三种接口：Publish、Subscribe和Unsubscribe。Subscribe接口的主要作用是，向事件聚合器注册指定类型事件的处理器，那么对于事件处理器而言，它就是在侦听（订阅）某个事件的发生；而Unsubscribe的作用则正好相反：它会解除某个事件处理器对指定类型事件的侦听，也就是当事件被触发时，不再侦听该事件的事件处理器将不会执行处理任务；至于Publish接口就非常简单了：领域模型使用Publish接口直接向事件聚合器派发事件，事件聚合器在观察到事件发生时，将处理权转交给侦听了该事件的处理器。事件聚合器的引入，使得事件能够被一次派发，多处处理，为应用程序的领域事件处理架构提供了扩展性的同时，也简化了事件订阅过程。

在Byteart Retail中，事件聚合器是一个静态类，之所以不设计成实例类，是因为我们无法将其以任何形式注射到领域模型中，更不可能让领域对象提供一个参数为EventAggregator类型的构造函数。这一点与保持领域模型的纯净度有关。Event Aggregator的具体实现代码，请参考ByteartRetail.Domain.Events命名空间下的DomainEventAggregator类。接下来，我们将领域事件的产生、订阅、派发和处理的过程总结一下。

领域事件的订阅、派发和处理
首先，在领域模型参与业务逻辑之前，应用程序架构需要对所需处理的领域事件进行订阅。回顾一下，面向DDD的经典分层架构中，应用层的职责是协调各组件（比如事务、仓储、领域模型等）的任务执行，因此领域事件的订阅也应该在应用层服务被初始化的时候进行。具体到Byteart Retail案例中，就是在应用服务（Application Service）的构造函数中进行。

以OrderServiceImpl类型（该类型位于ByteartRetail.Application.Implementation命名空间下）为例，在构造函数中我们扩展了一个参数：一个IDomainEventHandler<OrderDispatchedEvent>类型的数组，进而在构造函数中，通过使用DomainEventAggregator类，对传入的事件处理器进行订阅操作：

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
public OrderServiceImpl(IRepositoryContext context,
    IShoppingCartRepository shoppingCartRepository,
    IShoppingCartItemRepository shoppingCartItemRepository,
    IProductRepository productRepository,
    IUserRepository customerRepository,
    ISalesOrderRepository salesOrderRepository,
    IDomainService domainService,
    IDomainEventHandler<OrderDispatchedEvent>[] orderDispatchedDomainEventHandlers)
    :base(context)
{
    this.shoppingCartRepository = shoppingCartRepository;
    this.shoppingCartItemRepository = shoppingCartItemRepository;
    this.productRepository = productRepository;
    this.userRepository = customerRepository;
    this.salesOrderRepository = salesOrderRepository;
    this.domainService = domainService;
    this.orderDispatchedDomainEventHandlers.AddRange(orderDispatchedDomainEventHandlers);
 
    foreach (var handler in this.orderDispatchedDomainEventHandlers)
        DomainEventAggregator.Subscribe<OrderDispatchedEvent>(handler);
    DomainEventAggregator.Subscribe<OrderConfirmedEvent>(orderConfirmedEventHandlerAction);
    DomainEventAggregator.Subscribe<OrderConfirmedEvent>(orderConfirmedEventHandlerAction2);
}
构造函数中最后两行是对与OrderConfirmedEvent相关的事件处理委托进行订阅，以演示基于委托的事件处理器的实现方式。这两个委托在OrderServiceImpl类型中，以只读字段（readonly field）的形式进行定义：

1
2
3
4
5
6
7
8
9
10
11
private readonly Action<OrderConfirmedEvent> orderConfirmedEventHandlerAction = e =>
    {
        SalesOrder salesOrder = e.Source as SalesOrder;
        salesOrder.DateDelivered = e.ConfirmedDate;
        salesOrder.Status = SalesOrderStatus.Delivered;
    };
 
private readonly Action<OrderConfirmedEvent> orderConfirmedEventHandlerAction2 = _ =>
    {
         
    };
orderConfirmedEventHandlerAction2的定义无非也就是一个演示而已（演示接下来要讨论的事件处理器退订），因此我也没有在这个匿名方法里填写任何处理逻辑。至于构造函数的IDomainEventHandler<OrderDispatchedEvent>数组参数，则是通过Unity注入的，修改一下服务端的web.config文件即可：

SNAGHTMLbb5949e

接下来，在应用层完成操作后，需要解除事件处理器对事件的订阅（即退订），为了实现这个功能，我修改了IApplicationServiceContract的接口定义，并让ApplicationService类继承于DisposableObject类，之后，在WCF服务上，设置其InstanceContextMode为PerSession，也就是每当WCF客户端建立一次与服务端的连接时，创建一次服务实例，而当客户端关闭并撤销连接时，销毁服务实例。于是，在完成了这些结构调整后，每当一次WCF会话完成后，ApplicationService的Dispose方法就会被调用。那么每个应用层服务的具体实现（OrderServiceImpl、ProductServiceImpl、UserServiceImpl、PostbackServiceImpl）只需根据自己的需要重载Dispose方法，即可在Dispose方法中解除事件处理器对事件的订阅：

1
2
3
4
5
6
7
8
9
10
protected override void Dispose(bool disposing)
{
    if (disposing)
    {
        foreach (var handler in this.orderDispatchedDomainEventHandlers)
            DomainEventAggregator.Unsubscribe<OrderDispatchedEvent>(handler);
        DomainEventAggregator.Unsubscribe<OrderConfirmedEvent>(orderConfirmedEventHandlerAction);
        DomainEventAggregator.Unsubscribe<OrderConfirmedEvent>(orderConfirmedEventHandlerAction2);
    }
}
最后，领域事件的触发就非常简单了：直接调用DomainEventAggregator.Publish即可。整个过程大致可以用下面的序列图描述：

image

至此，我们已经大致了解了Byteart Retail案例中领域事件部分的设计与实现，回顾一下，这些内容包括：领域事件的定义、事件处理器、事件聚合器，以及这些组件之间的相互协作关系。读者朋友如果能够仔细阅读本案例的源代码，相信还能了解到更多深层次的细节问题。然而，事情还没有结束，我们还需要把讨论范围扩大到一个更高的层次：应用事件（Application Event）。虽然它已经超出领域事件的范围，但我还是要在本文中对其进行介绍，因为这个概念很容易造成开发人员对事件类别的混淆。

还有什么问题吗？
在本文最开始的时候提出了一个简单的应用场景：“当系统管理员完成销售订单的发货操作时，希望向客户发送一份电子邮件”，这种需求是最常见不过的了。虽然“完成销售订单的发货”被定义成一个领域事件（事实上它也就是一个领域事件），但处理电子邮件发送的逻辑，却并不是领域事件处理器的任务。通过分析不难得知，领域事件处理器对领域事件的处理，在于整个事务被提交之前。领域事件处理器可以以一种更为复杂的方式来获取或设置领域对象的状态，但对于与事务相关的事件处理过程，领域事件处理器就不是一个很好的选择。试想，如果在领域事件处理器中将电子邮件发送出去了，而接下来的事务提交却失败了，于是就造成了客户所收到的订单状态与实际状态不符的情形。

正确的做法应该是，在领域事件被触发时，将其记录下来，当执行事务提交时，将已记录的领域事件转换成应用事件，并派发到事件总线。这个派发过程可以是同步的，也可以是异步的。接下来的电子邮件发送逻辑就由侦听该事件总线的事件处理器负责执行。这里牵涉到一个分布式事务处理的问题。对于“发送电子邮件”这样的功能，我想，对分布式事务处理的要求应该也没有那么明显：数据库事务提交成功后，直接让基础结构层组件发送电子邮件就可以了，如果发送电子邮件失败，也完全无需回滚数据库事务。大不了客户抱怨说没有收到邮件，系统管理员通过事件日志对发送邮件的功能进行排错即可。但对于某些应用事件，比如客户订房成功后，系统就会将订房成功的事件发送到支付系统，支付系统在多次尝试付款失败后，就需要完成房间退订逻辑，以防止房间被无限制占用，在这些场景下，分布式事务处理就有着一定的必须性（当然你也可以说让支付系统无限制地重试，或者说找Sales Rep进行7x24的跟踪排错来解决事务问题，但我们暂时先不考虑这些解决方案）。

Byteart Retail考虑了这些问题存在的可能性，在事件系统和仓储部分大致进行了以下改动：

引入事件总线系统（IBus接口），应用事件处理器可以侦听该接口来接收需要处理的应用事件；应用层同样可以使用该接口来派发应用事件
实现了一个面向Event Dispatcher的事件总线，通过使用Event Dispatcher，Byteart Retail的事件总线可以支持Sequential、Parallel以及ParallelNoWait三种不同的事件派发方式（详见代码中的注释内容）
更改了AggregateRoot抽象类的实现，引入了存储领域事件的部分
更改了RepositoryContext抽象类的实现，在Commit方法中，不仅执行了仓储本身的提交事务（新的DoCommit方法），而且还会将存储在聚合根中的领域事件派发到事件总线。事件总线定义了其本身是否支持分布式事务处理，RepositoryContext会根据这个设置来决定是否需要启用Distributed Transaction Coordinator（不过貌似Message Queue的解决方案中，也只有MSMQ能够支持MS DTC）
详细的实现部分，我就不在这里一一叙述了，请读者朋友们自己阅读本案例的源代码，尤其是ByteartRetail.Events和ByteartRetail.Events.Handlers命名空间下的类型代码。

执行效果
本文最后，就让我们一起看一下领域事件部分的执行效果。以系统管理员发货为例，按理系统会产生一个OrderDispatchedEvent领域事件，领域模型通过领域事件处理器更新订单的发货日期和状态，与此同时，会将产生的领域事件暂存在聚合根中。当订单更新被提交时，被保存的领域事件将被派发到事件总线，进而邮件发送处理器会捕获到这个事件并发送邮件给客户。

首先，启动Byteart Retail的WCF服务和ASP.NET MVC应用程序，用daxnet/daxnet账户登录，并在账户设置中确保该账户的电子邮件地址设置正确。然后，使用该账户在系统中任意购买一件商品，完成下单后，退出系统，并用admin/admin账户登录，在“管理”->“销售订单管理”页面中，找到刚刚收到的订单，并点击“发货”按钮进行发货处理：

https://www.cnblogs.com/irocker/p/domain-events-pattern-example.html
本文展示的是一个关于网上调查的项目。想象下，当用户完成了一个调查，我们想通知所有人调查已经结束，分配一个人去检查调用问卷。

领域对象
public class Survey
{
    public Guid Id { get; private set; }
    public DateTime EndTime { get; private set; }
    public string QualityChecker { get; set; }

    public Survey()
    {
        this.Id = Guid.NewGuid();
    }

    public void EndSurvey()
    {
        EndTime = DateTime.Now;
        DomainEvent.Raise(new EndOfSurvey() { Survey = this });
    }
}
这个领域对象非常简单，只有一个行为：EndSurvey().

那么这里的DomainEvent是个什么东西呢？它是一个静态类，它发布了一个EndOfSurvey事件。从项目源码中可以看到所有的事件都放在名为Events的文件夹下面。领域对象放在Domain文件夹下面。

EndOfSurvey事件
现在Survey对象希望发布一个EndOfSurvey事件。这个事件的代码如下：

public class EndOfSurvey : IDomainEvent
{
    public Survey Survey { get; set; }
}
EndOfSurvey包含一个Survey实例。它继承自IDomainEvent，这样我们知道他是一个领域事件。本例中所有的事件都要继承自IDomainEvent。 这个接口的定义很简单：

public interface IDomainEvent { }  
DomainEvent类
public static class DomainEvent
{
    public static IEventDispatcher Dispatcher { get; set; }

    public static void Raise<T>(T @event) where T : IDomainEvent
    {
        Dispatcher.Dispatch(@event);
    }

}
源码中的DomainEvent比这个要复杂点，但最重要的便是上面的代码了。

IEventDispatcher是一个ioc容器。它负责找到正确的handler来处理EndOfSurvey事件。

public interface IEventDispatcher
{
    void Dispatch<TEvent>(TEvent eventToDispatch) where TEvent : IDomainEvent;
}
泛型方法Raise<T>能让我们发布无数的事件，Dispatcher自动找出对应的handler。

下面定义一个处理所有事件的handler接口：

public interface IDomainHandler<T> where T : IDomainEvent
{
    void Handle(T @event);
}
我将IEventDispatcher.cs和IDomainHandler.cs都放在一个名为Services的文件夹下面。其他的项目必须提供具体的实现。

domain程序集的代码就是这些了。

定义domain事件handler
我创建了另外一个项目用来写event handler。

EndOfSurveyHandler用来处理EndOfSurvey事件：

public class EndOfSurveyHandler:IDomainHandler<EndOfSurvey>
{
    public void Handle(EndOfSurvey args)
    {
        args.Survey.QualityChecker = "Ivan Amalo";
        // 发送邮件给Ivan，通知他来检查调查问卷
    }
}
如果想使用repository进行一些数据持续化的工作，或者使用一些其他的服务，可以将这些repository和服务通过构造函数注入进来。

EndOfSurveyHandler和EndOfSurvey事件是怎么联系起来的呢？

将所有的代码集成起来
下面要讲Survey.FrontEnd是一个MVC + WebApi应用，这个应用将DomainEvent，Dispatcher，Handler都结合了起来。

这个项目依赖于Ninject.MVC3。

现在我们需要来实现在之前定义的IEventDispatcher。

public class NinjectEventContainer : IEventDispatcher
{
    private readonly IKernel _kernel;

    public NinjectEventContainer(IKernel kernel)
    {
        _kernel = kernel;
    }

    public void Dispatch<TEvent>(TEvent eventToDispatch) where TEvent : IDomainEvent
    {
        foreach (var handler in _kernel.GetAll<IDomainHandler<TEvent>>())
        {
            handler.Handle(eventToDispatch);
        }
    }
}
Dispatch方法使用kernel来查找所有实现了IDomainHandler的handler。在我们的例子中查找的是EndOfSurveyHanlder，然后执行它的Handle()方法。

在NinjectWebCommon.cs中我们定义了handler和event的对应关系。

private static void RegisterServices(IKernel kernel)
{
    DomainEvent.Dispatcher = new NinjectEventContainer(kernel);
    kernel.Bind<IDomainHandler<EndOfSurvey>>().To<EndOfSurveyHandler>();
}   
这就是我们将所有东西集成起来需要做的事情。

测试
我在EndOfSurveyHandler.cs中发布事件的代码那设置了一个断点，来测试事件已经发布，其对应的handler也被执行。

控制器的代码非常简单，如下：

public ActionResult Index()
{
    var survey = new Core.Domain.Survey();
    survey.EndSurvey();

    return View(survey);
}
执行这个action， Ivan Amalo应该被分配成为这个调查问卷的检查者，并且将EndDate设为当前时间

这是使用Axon Framework探索CQRS架构的一系列帖子中的帖子。 建议在继续阅读本系列之前阅读本系列中的前几篇文章，因为这有助于形成围绕所讨论主题的连续性线索。 您可以从使用Axon Framework探索CQRS架构开始：简介

有一个Github项目（exploringCQRSwithAxon），其中包含一个简单的应用程序，伴随本系列中的帖子。 它是一个说明性的应用程序，允许在两个虚构账户之间进行借记，贷记和转账。 很明显这是一个简单的应用程序，它是有目的的。 其目的不是捕获任何复杂的域，而是帮助说明CQRS体系结构的各个组件以及如何使用Axon Framework构建这些不同的组件。

Following with the sample application.
要使项目处于说明此帖中的事物主题的状态：

首先下载项目:

git clone git@github.com:dadepo/exploringCQRSwithAxon.git
1
然后 check out 本次提交的 commit hash:

git checkout 06411af499a8d9dab62e1697820ca1c696f766dd
1
您可以通过执行mvn spring-boot来检查它之后运行应用程序：在根目录中运行，因为应用程序是使用Spring启动构建的。

在上一篇文章中，我们能够利用Axon的构建块来建立一个存储库，从中我们可以检索我们的聚合根。

我们还有适当的命令处理组件，以便我们能够调度最终改变应用程序状态的命令。 但正如该帖子末尾所述，尽管我们能够拥有最终导致状态变化的命令，但我们还没有触及CQRS的核心。

这是因为我们仍然使用相同的模型和基础结构来执行命令处理和查询组件。 正如在使用Axon Framework探索CQRS架构中所介绍的那样：简介CQRS的核心是具有处理写入（命令）和读取（查询）的独特且独立的组件。

那么，我们如何理解CQRS的核心并将命令组件与查询组件分开？ 为了回答这个问题，让我们再看一下我们在介绍性帖子中介绍的CQRS架构图：

cqrs

在图中，我们看到事件总线位于与命令/写入有关的体系结构的一侧和与读取/查询有关的一侧。 因此，为了回答我们的问题，我们利用事件总线来实现我们所寻求的分离。

How does this work? It work thus:
在写入方面，命令会导致domain中的状态更改
models/aggregates
domain中的状态更改会导致domain events捕获已更改的内容
domian 事件发布到事件总线。
在read/query方面，事件处理程序监听这些事件并使用它们传达的信息来维护应用程序状态的反映。 然后将此状态用于应用程序的读取端。
从域模型中的更改发布的事件在技术上称为域事件，Axon带有必要的基础结构，用于在域模型中发生更改时发布这些域事件。

本文的下一部分将介绍对其他示例应用程序所做的代码更改，以便连接并使用必要的Axon组件来实现命令端与查询端的这种分离，基本上是下图中突出显示的组件：

cqrs-event

Overview of code changes
Updating Init Script
由于我们将为查询提供另一个model/storage，因此我们更新了启动脚本以创建一个Account_View表，该表将用于读取应用程序的状态。 我们还插入了两个虚拟账号：表中的acc-one和acc-two。 所以我们的init（）方法现在看起来如此：

@PostConstruct
private void init(){
  // init the tables for commands
  TransactionTemplate transactionTmp = new TransactionTemplate(txManager);
  transactionTmp.execute(new TransactionCallbackWithoutResult() {
      @Override
      protected void doInTransactionWithoutResult(TransactionStatus status) {
          UnitOfWork uow = DefaultUnitOfWork.startAndGet();
          repository.add(new Account("acc-one"));
          repository.add(new Account("acc-two"));
          uow.commit();
      }
  });

  // init the tables for query/view
 JdbcTemplate jdbcTemplate = new JdbcTemplate(dataSource);
 jdbcTemplate
       .execute("create table account_view (account_no VARCHAR , 
                                                      balance FLOAT )");

 jdbcTemplate
       .update("insert into account_view (account_no, balance) values (?, ?)", 
                                              new Object[]{"acc-one", 0.0});

 jdbcTemplate
       .update("insert into account_view (account_no, balance) values (?, ?)", 
                                             new Object[]{"acc-two", 0.0});
}

Add an Event Bus
我们需要一个 event bus。它的基础架构允许将事件路由到事件处理程序。

您会注意到事件总线可能看起来类似于命令总线，因为它们都是消息调度基础结构。就功能而言，它们在本质上是不同的。

是的，命令总线调度命令，事件总线调度事件，但命令总线使用的命令用于表示在最近的将来需要发生的事情，并且它期望有一个且只有一个命令处理程序将解释并执行命令中捕获的意图。

另一方面，事件总线，路由事件和事件是过去发生的事情的表达，并且事件可能有零个或多个事件处理程序。

EventBus是描述事件调度组件的Axon接口。 Axon带有几个实现，对于我们的示例应用程序，我们将使用SimpleEventBus。因此我们将其连接为Spring bean：

/**
* The simple command bus, an implementation of an EventBus
* mostly appropriate in a single JVM, single threaded use case.
* @return the {@link SimpleEventBus}
*/
@Bean
public SimpleEventBus eventBus() {
  return new SimpleEventBus();
}

我们还需要连接Axon基础设施，以便能够轻松的设置事件处理程序，以响应发布到事件总线的事件。

Axon Framework附带@EventHandler注释，可用于将方法标记为事件处理程序。 方法的第一个参数表示方法应响应的事件类型。

AnnotationEventListenerBeanPostProcessor扫描具有@EventHandler方法的Spring bean，并自动将它们作为事件处理程序注册到事件总线。

AnnotationEventListenerBeanPostProcessor是AnnotationCommandHandlerBeanPostProcessor的事件计数器部分。

我们注册AnnotationEventListenerBeanPostProcessor：

@Bean
AnnotationEventListenerBeanPostProcessor 
                 annotationEventListenerBeanPostProcessor() {
  /**
   * The AnnotationEventListenerBeanPostProcessor 
     finds all beans that has methods annotated with @EventHandler
   * and subscribe them to the eventbus.
   */
  AnnotationEventListenerBeanPostProcessor listener = 
                               new AnnotationEventListenerBeanPostProcessor();
  listener.setEventBus(eventBus());
  return listener;
}
1
2
3
4
5
6
7
8
9
10
11
12
13
[UPDATE]
从版本2.3开始，Axon Framework提供了@AnnotationDriven注释，可以防止必须显式声明AnnotationEventListenerBeanPostProcessor类型的bean。 要使用它，只需使用@AnnotationDriven注释Spring @Configuration类，所有@ CommandHandler和@EventHandler将自动扫描并注册到各自的总线。 随附的示例应用程序已更新（使用d6c9f18750f8f7d4c341c80a07bdf44c5a815783提交）以使用@ AnnotationDriven。

接下来，我们需要对我们的配置进行的最后更新是为GenericJpaReposirtory提供事件总线。

Update the Repository with Event Bus
GenericJpaRepository需要事件总线，因为它将在我们的域对象中发生更改时发布域事件。 我们在下面的配置中提供事件总线：

@Bean
public GenericJpaRepository genericJpaRepository() {
  SimpleEntityManagerProvider entityManagerProvider = 
                      new SimpleEntityManagerProvider(entityManager);

  GenericJpaRepository genericJpaRepository = 
                      new GenericJpaRepository(entityManagerProvider, 
                                                         Account.class);

  /**
   * Configuring the repository with an event bus which allows the repository
   * to be able to publish domain events
   */
  genericJpaRepository.setEventBus(eventBus());
  return genericJpaRepository;
}

domain model的更改
然后我们转到Account.class，这是我们的设置中的Aggregate Root和Aggregate。 我们添加了允许在状态发生变化时发布domain event的代码。

我们的debit方法现在看起来如此：

public void debit(Double debitAmount) {

if (Double.compare(debitAmount, 0.0d) > 0 &&
  this.balance - debitAmount > -1) {
  this.balance -= debitAmount;

 /**
  * A change in state of the Account has occurred which can 
  * be represented by an to an event: i.e.
  * the account has been debited so an AccountDebitedEvent 
  * is created and registered.
  *
  * When the repository stores this change in state, it will 
  * also publish the AccountDebitedEvent
  * to the outside world.
  */
  AccountDebitedEvent accountDebitedEvent = 
              new AccountDebitedEvent(this.id, debitAmount, this.balance);
  registerEvent(accountDebitedEvent);

} else {
  throw new IllegalArgumentException("Cannot debit with the amount");
}

}

25
和 credit 方法:

public void credit(Double creditAmount) {

  if (Double.compare(creditAmount, 0.0d) > 0 &&
       Double.compare(creditAmount, 1000000) < 0) {

       this.balance += creditAmount;

  /**
  * A change in state of the Account has occurred which 
  * can be represented by an to an event: i.e.
  * the account has been credited so an AccountCreditedEvent 
  * is created and registered.
  *
  * When the repository stores this change in state, it will 
  * also publish the AccountCreditedEvent
  * to the outside world.
  */
  AccountCreditedEvent accountCreditedEvent = 
              new AccountCreditedEvent(this.id, creditAmount, this.balance);
  registerEvent(accountCreditedEvent);
 } else {
  throw new IllegalArgumentException("Cannot credit with the amount");
  }
}

我们在其中创建AccountDebitedEvent或AccountCreditedEvent并使用registerEvent()来注册创建的事件。

我们之所以有registerEvent()方法，因为我们的Account类扩展了AbstractAggregateRoot。

registerEvent()公开了Axon机制，用于在域对象保存到存储库时跟踪需要发布到事件总线的域事件。

我们提到AccountDebitEvent和AccountCreditedEvent作为表示帐户debited/credit的事件。 它们是域事件。 因此，代表这些event 的class如下：

public class AccountCreditedEvent {

  private final String accountNo;
  private final Double amountCredited;
  private final Double balance;

  public AccountCreditedEvent(String accountNo, 
                 Double amountCredited, Double balance) {
      this.accountNo = accountNo;
      this.amountCredited = amountCredited;
      this.balance = balance;
  }

  public String getAccountNo() {
      return accountNo;
  }

  public Double getAmountCredited() {
      return amountCredited;
  }

  public Double getBalance() {
      return balance;
  }
}

和

public class AccountDebitedEvent {
  private final String accountNo;
  private final Double amountDebited;
  private final Double balance;

  public AccountDebitedEvent(String accountNo, 
                 Double amountDebited, Double balance) {
      this.accountNo = accountNo;
      this.amountDebited = amountDebited;
      this.balance = balance;
  }

  public String getAccountNo() {
      return accountNo;
  }

  public Double getAmountDebited() {
      return amountDebited;
  }

  public Double getBalance() {
      return balance;
  }
}

到目前为止我们取得了什么成果？在我们继续之前，让我们先看看它。

我们现在有了域事件（AccountCreditedEvent和AccountDebitedEvent），我们更新了域对象以发布这些域事件，并且我们已经使用必要的基础结构更新了我们的配置，允许发布域事件。 我们需要添加的下一件事是事件处理程序。

Adding Event Handlers
我们添加了两个事件处理程序：AccountDebitedEventHandler

@Component
public class AccountDebitedEventHandler {

@Autowired
DataSource dataSource;

@EventHandler
public void handle AccountDebitedEvent(AccountDebitedEvent event) {

 JdbcTemplate jdbcTemplate = new JdbcTemplate(dataSource);

 // Get the current states as reflected in the event
 String accountNo = event.getAccountNo();
 Double balance = event.getBalance();

 // Update the view
 String updateQuery = "UPDATE account_view SET balance = ? 
                                             WHERE account_no = ?";
 jdbcTemplate.update(updateQuery, new Object[]{balance, accountNo});
 }
}

和 AccountCreditedEventHandler

@Component
public class AccountCreditedEventHandler {

@Autowired
DataSource dataSource;

@EventHandler
public void handleAccountCreditedEvent(AccountCreditedEvent event, 
Message eventMessage, @Timestamp DateTime moment) {

JdbcTemplate jdbcTemplate = new JdbcTemplate(dataSource);

// Get the current states as reflected in the event
String accountNo = event.getAccountNo();
Double balance = event.getBalance();

// Update the view
String updateQuery = "UPDATE account_view SET balance = ? 
                                           WHERE account_no = ?";
jdbcTemplate.update(updateQuery, new Object[]{balance, accountNo});

System.out.println("Events Handled With EventMessage " + 
            eventMessage.toString() + " at " + moment.toString());
}
}

可以看出，事件处理方法使用@EventHandler进行注释。 它们响应的事件类型由带注释的方法的第一个参数指示。

由于我们已在配置中注册了AnnotationEventListenerBeanPostProcessor，因此这些类将被订阅为事件总线的事件处理程序。

那么在这些事件处理方法中会发生什么呢？ 我们从各个事件中提取信息，并使用JDBC更新Account_view表，该表是用于查询/读取操作的表。

在这个结点上需要注意的一件重要事情是，我们将视图数据存储在与用于命令操作的表不同的表中。
我们也没有对视图层使用任何ORM映射（JPA等），我们也没有任何特殊的类在映射其状态时为Account建model。 这种特性是CQRS看待事物的核心。 事实上，使用CQRS，我们的查询层可以使用不同的简化抽象实现，可以轻松地针对查询/读取操作进行优化。

handleAccountCreditedEvent()方法展示了Axon为事件处理提供的一些附加功能。 可以看出，我们有两个额外的参数。 消息eventMessage包含事件消息：metadata，id等和@Timestamp DateTime时刻，它是在事件发布的时刻注入的。

接下来要做的是更新我们的视图（ViewController / Javascript）以使用这个新设置。

更新视图
我们更新ViewController的getAccounts方法以使用普通JDBC来查询帐户的状态：

@Controller
public class ViewController {

 @Autowired
 private DataSource dataSource;

 @RequestMapping(value = "/view", 
                method = RequestMethod.GET, 
                produces = MediaType.APPLICATION_JSON_VALUE)
 @ResponseBody
 public List<Map<String, Double>> getAccounts() {

 JdbcTemplate jdbcTemplate = new JdbcTemplate(dataSource);
   List<Map<String, Double>> queryResult = 
   jdbcTemplate.query("SELECT * from account_view ORDER BY account_no", 
                                                  (rs, rowNum) -> {
   return new HashMap<String, Double>() \{\{
                           put(rs.getString("ACCOUNT_NO"),
                           rs.getDouble("BALANCE"));
   \}\};
});

  return queryResult;
 }
}

轮询 /view endpoint 的JavaScript仍然存在。

通过所有这些更改，当您运行应用程序时，您仍然可以选择信用卡或借记卡，并在余额部分中反映余额。 就在这一次，借记/贷记是通过与用于查看账户当前余额的组件不同的组件完成的。

Overview of the Axon Building Blocks
在这篇文章中，我们讨论了Axon Framework中的一些新构建块。 让我们来迅速回顾一下：

SimpleEventBus
Axon基础设施负责事件路由到事件处理程序。

AnnotationEventListenerBeanPostProcessor
在Spring应用程序中使用手动方式使用Axon。它扫描具有@EventHandler注释的Spring bean，并自动将它们注册为事件处理程序到事件总线。

RegisterEvent method
可以通过在扩展AbstractAggregateRoot的域对象中的方法，以注册用于发布的域事件。

概要
到目前为止，我们已经能够以CQRS的方式，设置应用程序的组件。

您会注意到我们在没有提及或使用事件溯源（Event Sourcing）的情况下完成了这项工作，这表明如果您不需要它时，您可以在不使用事件源的情况下构建CQRS应用程序。

但是如果你想使用Event Sourcing怎么办？ Axon Framework如何提供帮助？ 下一篇文章探索使用Axon Framework的CQRS：应用事件采购答案，询问如何使用Axon Framework在CQRS应用程序中使用事件溯源。

https://blog.csdn.net/quguang65265/article/details/81382319


https://www.cnblogs.com/daxnet/archive/2013/04/30/3052029.html

http://www.360doc.com/content/13/1226/09/10504424_340184984.shtml
https://www.cnblogs.com/uoyo/p/12421553.html


https://www.jdon.com/eda.html
https://www.jdon.com/48068
https://blog.christianposta.com/microservices/why-microservices-should-be-event-driven-autonomy-vs-authority/

https://www.jdon.com/49113

https://www.jdon.com/eda.html
https://www.jdon.com/event.html

https://www.jdon.com/49081
