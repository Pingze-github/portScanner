# portScanner
高效率端口扫描器。Go语言编写。

### 编译
```go run main.go```
#### 交叉编译
Linux x64:
```
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go
```
Windows x64:
```
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build main.go
```

### 使用
编译为 main 后，

```./main -ip 127.0.0.1 -port 1-65535 -max 10000```

> port 有三种指定方式：
> 1. `100` 指定单个
> 2. `100,101,102` 指定多个
> 3. `1-100` 指定范围

输出所有开端口（tcp）：

```21,80,111,3166```

### 测试

条件：并发数: 10000; 超时时间: 3s; 扫描范围: 1-65535

+ 拒绝型服务器。耗时 7.6s
+ 超时型服务器。耗时 21.2s

(对全部端口超时的服务器全扫描，理论耗时 Ceil(65536/10000)*3s = 21s。实测结果已经极其逼近理论值)


