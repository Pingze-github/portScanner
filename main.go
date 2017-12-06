package main

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"time"
)

func scanRange(ip string, portRange [2]int, gomax int, timeout int) []int {

	finish := make(chan int)
	channel := make(chan int, gomax)
	var openPorts []int

	scan := func (ip string, port int) {
		// time.Sleep(time.Duration(timeout)) // 模拟超时
		address := ip + ":" + strconv.Itoa(port)
		_, err := net.DialTimeout("tcp", address, time.Duration(timeout))
		if err != nil {
			fmt.Println(err)
		} else {
			openPorts = append(openPorts, port)
		}
		i := <- channel
		if i == 1 {
			finish <- 0
		}
	}

	for port := portRange[0]; port <= portRange[1]; port++ {
		if port == portRange[1] {
			channel <- 1
		} else {
			channel <- 0
		}
		go scan(ip, port)
	}
	<- finish
	return openPorts
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ip := "127.0.0.1"
	var portRange = [2]int{1, 65535}
	gomax := 10000
	start := time.Now().UnixNano()
	openPorts := scanRange(ip, portRange, gomax, 3e9)
	fmt.Println(openPorts)
	fmt.Println((time.Now().UnixNano() - start) / 1e6)
}

