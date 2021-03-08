# sort

sort 包中实现了：插入排序、归并排序、堆排序和快速排序。但是是内部方法，所以在对数据集合排序时不必考虑应当选择哪一种排序方法，只要
实现了 `sort.Interface` 定义的三个方法就可以对数据集合进 行排序。`sort` 包会根据实际数据自动选择高效的排序算法。

```go
type Interface interface {
	// Len 为集合内元素的总数
	Len() int
	// 如果 index 为 i 的元素小于 index 为 j 的元素，则返回 true，否则 false
	Less(i, j int) bool
	// Swap 交换索引为 i 和 j 的元素
	Swap(i, j int)
}
```

sort 包原生支持 `[]int`、`[]float64` 和 `[]string` 三种内建数据类型切片的排序操作。不必实现 `sort.Interface` 接口的三个方法。

数据集合实现了 `sort.Interface` 定义的三个方法后，即可调用该包的 `Sort()` 方法进行排序。

`Sort()` 方法定义如下：

```go
func Sort(data Interface)
```

`Sort()` 方法惟一的参数就是待排序的数据集合。

该包还提供了一个方法可以判断数据集合是否已经排好顺序，该方法的内部实现依赖于我们自己实现的 `Len()` 和 `Less()` 方法：

```go
func IsSorted(data Interface) bool {
    n := data.Len()
    for i := n - 1; i > 0; i-- {
        if data.Less(i, i-1) {
            return false
        }
    }
    return true
}
```

## Reverse

v1 代码实现的是升序排序，如果要实现降序排序修改 `Less()` 函数：

```go
// 将小于号修改为大于号
func (s StuScores) Less(i, j int) bool {
	return s[i].score > s[j].score
}
```

此外，sort 包提供了 `Reverse()` 方法，可以允许将数据按 `Less()` 定义的排序方式逆序排序，而不必修改 `Less()` 代码。

```go
func Reverse(data Interface) Interface
```

`Reverse()` 返回的一个 `sort.Interface` 接口类型，整个 `Reverse()` 的内部实现比较有趣：

```go
// 定义了一个 reverse 结构类型，嵌入 Interface 接口
type reverse struct {
    Interface
}

// reverse 结构类型的 Less() 方法拥有嵌入的 Less() 方法相反的行为
func (r reverse) Less(i, j int) bool {
    return r.Interface.Less(j, i)
}

// 返回新的实现 Interface 接口的数据类型
func Reverse(data Interface) Interface {
    return &reverse{data}
}
```

## Search
```go
func Search(n int, f func(int) bool) int
```

该函数使用二分查找（**必须是已经生序排序的集合**）的方法，从 `[0, n)` 中取出一个下标值 `index`，`index` 为 `[0, n)` 中最小的，使函数 `f(index)` 为 `true` 的值，并且
`f(index+1)` 也为 `true`。

如果无法找到该 `index` 值，则该方法返回 `n`。

`Search()` 函数一个常用的使用方式是搜索元素 `x` 是否在已经升序排好的切片 `s` 中。

## 3.1.2 *sort*包已经支持的内部数据类型排序

sort 包原生支持 `[]int`、`[]float64` 和 `[]string` 三种内建数据类型切片。

### IntSlice 类型及 []int 排序

由于 `[]int` 切片排序内部实现及使用方法与 `[]float64` 和 `[]string` 类似。

sort 包定义了一个 `IntSlice` 类型，并且实现了 `sort.Interface` 接口：

```go
type IntSlice []int
func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
//IntSlice 类型定义了 Sort() 方法，包装了 sort.Sort() 函数
func (p IntSlice) Sort() { Sort(p) }
//IntSlice 类型定义了 SearchInts() 方法，包装了 SearchInts() 函数
func (p IntSlice) Search(x int) int { return SearchInts(p, x) }
```

`sort.Ints()` 方法使用了该 `IntSlice` 类型：

```go
func Ints(a []int) { Sort(IntSlice(a)) }
```

所以，对 `[]int` 切片排序更常使用 `sort.Ints()`，而不是直接使用 `IntSlice` 类型：
```go
s := []int{5, 2, 6, 3, 1, 4} // 未排序的切片数据
sort.Ints(s)
fmt.Println(s) // 将会输出[1 2 3 4 5 6]
```

如果要使用降序排序，要用 `Reverse()` 方法：
```go
s := []int{5, 2, 6, 3, 1, 4} // 未排序的切片数据
sort.Sort(sort.Reverse(sort.IntSlice(s)))
fmt.Println(s) // 将会输出[6 5 4 3 2 1]
```

如果要查找整数 `x` 在切片 `a` 中的位置，相对于前面提到的 `Search()` 方法，sort 包提供了 `SearchInts()`:
```go
func SearchInts(a []int, x int) int
```

### Float64Slice 类型及[]float64 排序

实现与 `Ints` 类似：
```go
type Float64Slice []float64

func (p Float64Slice) Len() int           { return len(p) }
func (p Float64Slice) Less(i, j int) bool { return p[i] < p[j] || isNaN(p[i]) && !isNaN(p[j]) }
func (p Float64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Float64Slice) Sort() { Sort(p) }
func (p Float64Slice) Search(x float64) int { return SearchFloat64s(p, x) }
```

与 `Sort()`、`IsSorted()`、`Search()` 相对应的三个方法：
```go
func Float64s(a []float64)  
func Float64sAreSorted(a []float64) bool
func SearchFloat64s(a []float64, x float64) int
```

`Float64Slice` 类型定义的 `Less` 方法中，有一个内部函数 `isNaN()`。
`isNaN()` 与 math 包中 `IsNaN()` 实现完全相同，sort 包之所以不使用 `math.IsNaN()`，完全是基于包依赖性的考虑，应当看到，sort 包的实现不依赖与其他任何包。

### StringSlice 类型及[]string 排序

两个 `string` 对象之间的大小比较是基于“字典序”的。

实现与 `Ints` 类似：

```go
type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p StringSlice) Sort() { Sort(p) }
func (p StringSlice) Search(x string) int { return SearchStrings(p, x) }
```

与 `Sort()`、`IsSorted()`、`Search()` 相对应的三个方法：
```go
func Strings(a []string)
func StringsAreSorted(a []string) bool
func SearchStrings(a []string, x string) int
```

## []interface 排序与查找

只要实现了 `sort.Interface` 接口，即可通过 sort 包内的函数完成排序，查找等操作。并且 sort 包已经帮我们把`[]int`,`[]float64`,`[]string` 三种类型都实现
了该接口，我们可以方便的调用。但是这种用法对于其它数据类型的 `slice` 不友好，可能我们需要为大量的 `struct` 定义一个单独的 `[]struct` 类型，再
为其实现 `sort.Interface` 接口。

### 为什么 sort 包可以完成 `[]int` 的排序，而不能完成 `[]struct` 的排序？

因为排序涉及到比较两个变量的值，而 `struct` 可能包含多个属性，程序并不知道你想以哪一个属性或哪几个属性作为衡量大小的标准。如果你能帮助程序完成比较，
并将结果返回， sort 包内的方法就可以完成排序，判断，查找等。

### []interface 的排序

sort 包提供了以下函数：

```go
func Slice(slice interface{}, less func(i, j int) bool) 
func SliceStable(slice interface{}, less func(i, j int) bool) 
func SliceIsSorted(slice interface{}, less func(i, j int) bool) bool 
func Search(n int, f func(int) bool) int
```

通过函数签名可以看到，排序相关的三个函数都接收 `[]interface`，并且需要传入一个比较函数，用于为程序比较两个变量的大小。

#### sort.Slice

该函数完成 `[]interface` 的排序：
```go
people := []struct {
	Name string
	Age  int
}{
	{"Gopher", 7},
	{"Alice", 55},
	{"Vera", 24},
	{"Bob", 75},
}

sort.Slice(people, func(i, j int) bool { return people[i].Age < people[j].Age }) // 按年龄升序排序
fmt.Println("Sort by age:", people)
```

输出结果：
```bash
By age: [{Gopher 7} {Vera 24} {Alice 55} {Bob 75}]
```

#### sort.SliceStable

该函数完成 `[]interface` 的稳定排序：
```go
people := []struct {
	Name string
	Age  int
}{
	{"Gopher", 7},
	{"Alice", 55},
	{"Vera", 24},
	{"Bob", 75},
}

sort.SliceStable(people, func(i, j int) bool { return people[i].Age > people[j].Age }) // 按年龄降序排序
fmt.Println("Sort by age:", people)
```

输出结果：
```bash
By age: [{Bob 75} {Alice 55} {Vera 24} {Gopher 7}]
```

#### sort.SliceIsSorted

该函数判断 `[]interface` 是否为有序：
```go
people := []struct {
    Name string
    Age  int
}{
    {"Gopher", 7},
    {"Alice", 55},
    {"Vera", 24},
    {"Bob", 75},
}

sort.Slice(people, func(i, j int) bool { return people[i].Age > people[j].Age }) // 按年龄降序排序
fmt.Println("Sort by age:", people)
fmt.Println("Sorted:",sort.SliceIsSorted(people,func(i, j int) bool { return people[i].Age < people[j].Age }))
```

输出结果：
```bash
Sort by age: [{Bob 75} {Alice 55} {Vera 24} {Gopher 7}]
Sorted: false
```

sort 包没有为 `[]interface` 提供反序函数，因为我们传入的比较函数已经决定了排序结果是升序还是降序。

#### sort.Search

该函数判断 `[]interface` 是否存在指定元素：

- 升序 slice

sort 包为 `[]int`,`[]float64`,`[]string` 提供的 Search 函数其实也是调用的该函数，因为该函数是使用的二分查找法，所以要求 slice 为升序
排序状态。并且判断条件必须为 `>=`，这也是官方库提供的三个查找相关函数的的写法。

```go
a := []int{2, 3, 4, 200, 100, 21, 234, 56}
x := 21

sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })   // 升序排序
index := sort.Search(len(a), func(i int) bool { return a[i] >= x }) // 查找元素

if index < len(a) && a[index] == x {
	fmt.Printf("found %d at index %d in %v\n", x, index, a)
} else {
	fmt.Printf("%d not found in %v,index:%d\n", x, a, index)
}
```

输出结果：
```bash
found 21 at index 3 in [2 3 4 21 56 100 200 234]
```

- 降序 slice

如果 slice 是降序状态，而我们又不想将其变为升序，只需将判断条件由 `>=` 变更为 `<=` 即可。