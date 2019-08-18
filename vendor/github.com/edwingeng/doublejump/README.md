# Overview
This is a revamped [Google's jump](https://arxiv.org/pdf/1406.2294.pdf) consistent hash. It overcomes the shortcoming of the original implementation - not being able to remove nodes. [Here is the main idea behind doublejump](https://docs.google.com/presentation/d/e/2PACX-1vTHyFGUJ5CBYxZTzToc_VKxP_Za85AeZqQMNGLXFLP1tX0f9IF_z3ys9-pyKf-Jj3iWpm7dUDDaoFyb/pub?start=false&loop=false&delayms=3000).

# Benchmark
```
BenchmarkDoubleJumpWithoutLock/10-nodes             50000000            27.6 ns/op
BenchmarkDoubleJumpWithoutLock/100-nodes            30000000            42.7 ns/op
BenchmarkDoubleJumpWithoutLock/1000-nodes           30000000            54.1 ns/op

BenchmarkDoubleJump/10-nodes                        20000000            72.9 ns/op
BenchmarkDoubleJump/100-nodes                       20000000            86.1 ns/op
BenchmarkDoubleJump/1000-nodes                      20000000            97.9 ns/op

BenchmarkStathatConsistent/10-nodes                  5000000           301 ns/op
BenchmarkStathatConsistent/100-nodes                 5000000           334 ns/op
BenchmarkStathatConsistent/1000-nodes                3000000           444 ns/op

BenchmarkSerialxHashring/10-nodes                    5000000           280 ns/op
BenchmarkSerialxHashring/100-nodes                   5000000           340 ns/op
BenchmarkSerialxHashring/1000-nodes                  3000000           427 ns/op
```

# Example
```
h := NewHash()
for i := 0; i < 10; i++ {
    h.Add(fmt.Sprintf("node%d", i))
}

fmt.Println(h.Len())
fmt.Println(h.LooseLen())

fmt.Println(h.Get(1000))
fmt.Println(h.Get(2000))
fmt.Println(h.Get(3000))

h.Remove("node3")
fmt.Println(h.Len())
fmt.Println(h.LooseLen())

fmt.Println(h.Get(1000))
fmt.Println(h.Get(2000))
fmt.Println(h.Get(3000))

// Output:
// 10
// 10
// node9
// node2
// node3
// 9
// 10
// node9
// node2
// node0
```

# Acknowledgements
The implementation of the original algorithm is credited to [dgryski](https://github.com/dgryski/go-jump).
