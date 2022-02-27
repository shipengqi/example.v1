# Casbin

在 Casbin 中, 访问控制模型被抽象为基于 PERM (Policy, Effect, Request, Matcher) 的一个文件。可以通过组合可用的模型来定制访问控制模型。
例如，可以在一个model中结合RBAC角色和ABAC属性，并共享一组 policy 规则。

## PERM模式

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


## RBAC

`[role_definition]` 是RBAC角色继承关系的定义

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
3. **RBAC 系统中的用户名称和角色名称不应相同。因为Casbin将用户名和角色识别为字符串， 所以当前语境下Casbin无法得出这个字面量到底指代用户 alice 还是角色 alice**。 这时，使用明确的 role_alice ，问题便可迎刃而解。
4. 假设A具有角色 B，B 具有角色 C，并且 A 有角色 C。 这种传递性在当前版本会造成死循环。
5. 定义名字只能是 `r`，`g`，`g2` 之类，否则报错
6. domain `*` is unsupport in `casbin/v2 v2.41.1`

## ABAC

ABAC是基于属性的访问控制，可以使用主体、客体或动作的属性，而不是字符串本身来控制访问。在ABAC中，可以使用 `struct` (或基于编程语言的类实例) 而不是字符
串来表示模型元素。

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
