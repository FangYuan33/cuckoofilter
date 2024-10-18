## Cuckoo Filter

[![GoDoc](https://godoc.org/github.com/seiflotfy/cuckoofilter?status.svg)](https://godoc.org/github.com/seiflotfy/cuckoofilter) [![CodeHunt.io](https://img.shields.io/badge/vote-codehunt.io-02AFD1.svg)](http://codehunt.io/sub/cuckoo-filter/?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

布谷鸟过滤器是用于近似集合成员查询的布隆过滤器的替代品。虽然布隆过滤器是众所周知的空间高效的数据结构，用于处理诸如“某元素 x 是否在集合中？”之类的查询，但它们不支持删除操作。为了支持删除操作的变种（如计数布隆过滤器）通常需要更多的空间。

布谷鸟过滤器提供了动态添加和删除项目的灵活性。布谷鸟过滤器基于布谷鸟哈希（因此得名布谷鸟过滤器）。它本质上是一个存储每个键指纹的布谷鸟哈希表。布谷鸟哈希表可以非常紧凑，因此对于需要低误报率（< 3%）的应用，布谷鸟过滤器可能比传统的布隆过滤器占用更少的空间。

有关算法和引用的详细信息，请暂时使用这篇文章。

["Cuckoo Filter: Better Than Bloom" by Bin Fan, Dave Andersen and Michael Kaminsky](https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf)

### 实现细节

上面引用的论文中留给了若干参数供选择。在这个实现中：

1. 每个元素有 2 个可能的桶索引
2. 每个桶的静态大小为 4 个指纹
3. 指纹的静态大小为 8 位

1 和 2 被作者建议为最佳选择。3 的选择取决于所需的假阳率。给定目标误报率 `r` 和桶大小 `b`，他们建议选择指纹大小 `f` 使用以下公式：

    f >= log2(2b/r) 位

在这个代码库中使用 8 位指纹大小时，你可以预期 `r ~= 0.03`。[其他实现](https://github.com/panmari/cuckoofilter) 使用 16 位指纹大小，对应的误报率为 `r ~= 0.0001`。

### Example

```go
package main

import "fmt"
import cuckoo "github.com/seiflotfy/cuckoofilter"

func main() {
  cf := cuckoo.NewFilter(1000)
  cf.InsertUnique([]byte("geeky ogre"))

  // Lookup a string (and it a miss) if it exists in the cuckoofilter
  cf.Lookup([]byte("hello"))

  count := cf.Count()
  fmt.Println(count) // count == 1

  // Delete a string (and it a miss)
  cf.Delete([]byte("hello"))

  count = cf.Count()
  fmt.Println(count) // count == 1

  // Delete a string (a hit)
  cf.Delete([]byte("geeky ogre"))

  count = cf.Count()
  fmt.Println(count) // count == 0
  
  cf.Reset()    // reset
}
```
