# portScanner
高效率TCP端口扫描器。Go语言编写。

### 编译
```go build main.go```
#### 交叉编译
Linux x64:
```
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go
```
Windows x32:
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

### 和(Nmap)[https://nmap.org/man/zh/]对比
1. 并发上，portScanner可以达到极高并发；Nmap不并发。portScanner占用的端口资源远多于Nmap（等于并发数），
但检测速度也远高于Nmap。
2. 检测方式上，Nmap可以使用TCP-SYN检测（快）、TCP-connect、UDP、TCP+version检测，
功能更全。但除第一种外，都很慢。
<br>portScanner正在准备加入SYN检测的功能，届时端口占用更少，检测更快。
<br>portScanner尚不支持UDP检测，后续可能会加入。
<br>portScanner不会加入对端口号和服务类型的映射，也不会加入对具体服务版本的检测。
这超出了portScanner的功能范围。
3. Nmap是一个功能全面的，注重隐蔽性的目标探测工具。
<br>portScanner的定位是一个微型、高并发、功能直接的端口开闭扫描工具。


