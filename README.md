# go-gc
Illustration of Go Garbage collector work

# Running
```go
$ go run main.go
```
Press Ctrl+C when you feel like stopping; the memory will be cleaned and you'll see the resulting output with freed allocated memory:
```
Alloc = 0 MiB   TotalAlloc = 0 MiB      Sys = 68 MiB    NumGC = 0
Alloc = 7 MiB   TotalAlloc = 7 MiB      Sys = 68 MiB    NumGC = 1
Alloc = 15 MiB  TotalAlloc = 15 MiB     Sys = 68 MiB    NumGC = 2
Alloc = 15 MiB  TotalAlloc = 15 MiB     Sys = 68 MiB    NumGC = 2
Alloc = 0 MiB   TotalAlloc = 15 MiB     Sys = 68 MiB    NumGC = 3
```
