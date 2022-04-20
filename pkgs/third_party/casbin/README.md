# Casbin

在 Casbin 中, 访问控制模型被抽象为基于 PERM (Policy, Effect, Request, Matcher) 的一个文件。可以通过组合可用的模型来定制访问控制模型。
例如，可以在一个 model 中结合 RBAC 角色和 ABAC 属性，并共享一组 policy 规则。

## PERM 模式

### Request
定义请求参数。基本请求是一个元组对象，至少需要主题(访问实体)、对象(访问资源) 和动作(访问方式)

例如，一个请求可能长这样： `r={sub,obj,act}`

它实际上定义了我们应该提供访问控制匹配功能的参数名称和顺序。

### Policy
定义访问策略模式。事实上，它是在政策规则文件中定义字段的名称和顺序。

例如： `p={sub, obj, act}` 或 `p={sub, obj, act, eft}`

> 注：如果未定义 eft (policy result)，则策略文件中的结果字段将不会被读取， 和匹配的策略结果将默认被允许。


### Matcher

匹配请求和政策的规则。

例如： `m = r.sub == p.sub && r.act == p.act && r.obj == p.obj` 这个简单和常见的匹配规则意味着如果请求的参数(sub,obj,act)
可以在策略中匹配到资源和方法，那么策略结果（`p.eft`）便会返回。 策略的结果将保存在 `p.eft` 中。


### Effect

它可以被理解为一种模型，在这种模型中，对匹配结果再次作出逻辑组合判断。

例如： `e = some (where (p.eft == allow))`

这句话意味着，如果匹配的策略结果有一些是允许的，那么最终结果为真。

另一个示例： `e = some (where (p.eft == allow)) && !some(where (p.eft == deny)` 此示例组合的逻辑含义是：
如果有符合允许结果的策略且没有符合拒绝结果的策略， 结果是为真。 换言之，当匹配策略均为允许（没有任何否认）是为真（更简单的是，既允许又同时否认，拒绝就具有优先地位)。

## Model 语法

```bash
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
p2 = sub, act
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
```

`[request_definition]` 部分用于 request 的定义，它明确了 `e.Enforce(...)` 函数中参数的含义。

你可以自定义请求表单, 如果不需要指定特定资源，则可以这样定义 `sub`、`act` ，或者如果有两个访问实体, 则为 `sub`、`sub2`、`obj`、`act`。

`[policy_definition]` 部分是对 policy 的定义，policy 规则的具体描述：

```bash
p, alice, data1, read
p2, bob, write-all-objects
```

policy 部分的每一行称之为一个策略规则， 每条策略规则通常以形如 p, p2 的 policy type 开头。 如果存在多个 policy 定义，那么我们会根据前文提
到的 policy type 与具体的某条定义匹配。上面的 policy 的绑定关系将会在 matcher 中使用，如下：

```bash
(alice, data1, read) -> (p.sub, p.obj, p.act)
(bob, write-all-objects) -> (p2.sub, p2.act)
```

`[policy_effect]` 是策略效果的定义

`e = some(where (p.eft == allow))`
上面的策略效果表示如果有任何匹配的策略规则 `allow`, 最终效果是 `allow` (aaka allow-override). `p.eft` 是政策的效果，它可以 `allow` 或 `deny`。 是可选的，
默认值是 `allow`。 因为我们没有在上面指定它，所以它使用默认值。

`e = !some(where (p.eft == deny))`
这意味着如果没有匹配的政策规则为 `deny`, 最终效果是 `allow` (别名为拒绝). `some` 表示：如果存在一个匹配的策略规则。 `any` 意味着：
所有匹配的政策规则(这里不使用)。 策略效果甚至可以与逻辑表达式相关联：

`some(where (p.eft == allow)) && !some(where (p.eft == deny))`

这意味着至少有一个匹配的策略规则 `allow`，并且没有匹配的 `deny` 的策略规则。 因此，允许和拒绝授权都得到支持，拒绝则被推翻。

`[matchers]` 是策略匹配程序的定义

## RBAC

`[role_definition]` 是 RBAC 角色继承关系的定义

```bash
[role_definition]
g = _, _
g2 = _, _
```

`g` 是一个 RBAC系统, `g2` 是另一个 RBAC 系统。 `_, _` 表示角色继承关系的前项和后项，即前项继承后项角色的权限。

### Tenant
```bash
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
```

1. Casbin 只存储用户角色的映射关系。
2. Cabin 没有验证用户是否是有效的用户，或者角色是一个有效的角色。 这应该通过认证来解决。
3. **RBAC 系统中的用户名称和角色名称不应相同。因为 Casbin 将用户名和角色识别为字符串， 所以当前语境下 Casbin 无法得出这个字面量到底指代用户 alice 还是角色 alice**。 这时，使用明确的 role_alice ，问题便可迎刃而解。
4. 假设A具有角色 B，B 具有角色 C，并且 A 有角色 C。 这种传递性在当前版本会造成死循环。
5. 定义名字只能是 `r`，`g`，`g2` 之类，否则报错
6. domain `*` is unsupport in `casbin/v2 v2.41.1`

## ABAC

ABAC 是基于属性的访问控制，可以使用主体、客体或动作的属性，而不是字符串本身来控制访问。在ABAC中，可以使用 `struct` (或基于编程语言的类实例) 而不是字符
串来表示模型元素。

属性通常来说分为四类：用户属性（如用户年龄），环境属性（如当前时间），操作属性（如读取）和对象属性（如一篇文章，又称资源属性），所以理论上能够实现非常灵活的权限控制，几乎能满足所有类型的需求。
例如规则：“允许所有班主任在上课时间自由进出校门”这条规则，其中，“班主任”是用户的角色属性，“上课时间”是环境属性，“进出”是操作属性，而“校门”就是对象属性了。

```bash
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == r.obj.Owner
```

在 `matcher` 中使用 `r.obj.Owner` 代替 `r.obj`。 在 `Enforce()` 函数中传递的 `r.obj` 是结构或类实例，而不是字符串。Casbin 使用反射来检索
`obj` 中的成员变量。

`r.obj` struct 示例：
```go
type testResource struct {
    Name  string
    Owner string
}
```

1. 仅有形如 `r.sub`, `r.obj`, `r.act` 等请求元素支持 ABAC。 不能在 policy 元素上使用它，比如 `p.sub`

简单地说，要使用 ABAC，需要做两件事：

1. 在模型匹配器中指定属性。
2. 将元素的结构或类实例作为 Casbin 的 `Enforce()` 的参数传入。

### 适配复杂且大量的 ABAC 规则
上述 ABAC 实例的核心非常简单，但授权系统通常需要非常复杂和大量的 ABAC 规则。 为了满足这一需要，上述实现将在很大程度上使得模型冗长复杂。 因此，
我们可以选择在策略中添加规则代替在模型中添加规则。 这是通过引入一个 `eval()` 功能结构完成的。 下面是此类 ABAC 模型的示例。

```bash
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub_rule, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = eval(p.sub_rule) && r.obj == p.obj && r.act == p.act
```

`p.sub_rule` 是由策略中使用的必要属性组成的结构类型或类类型(用户定义类型)。

这是针对 Enforcement 模型使用的策略 现在就可以使用作为参数传递到 `eval()` 函数的对象实例来定义某些 ABAC 约束条件。

```bash
p, r.sub.Age > 18, /data1, read
p, r.sub.Age < 60, /data2, write
```

## Priority Model
### 隐式优先级加载策略
顺序决定了策略的优先级，策略出现的越早优先级就越高。

```bash
[policy_effect]
e = priority(p.eft) || deny
```

### 显式优先级加载策略
优先级的名称必须是 `priority`，值越小优先级越高。优先级的值如果是非数字字符，将被排在最后，而不是导致报错。

### 基于角色和用户层次结构以优先级加载策略
角色和用户的继承结构只能是多棵树，而不是图。 如果一个用户有多个角色，必须确保用户在不同树上有相同的等级。 如果两种角色具有相同的等级，那么出现早
的策略（相应的角色）优先级更高。

## Root

```bash
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act || r.sub == "root"
```
