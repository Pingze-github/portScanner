package main

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"time"
	"flag"
	"regexp"
	"strings"
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
			//fmt.Println(err)
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

func join(intArr []int) string {
	str := ""
	for i := 0; i < len(intArr); i++ {
		str += strconv.Itoa(intArr[i])
		if i != len(intArr) - 1 {
			str += ","
		}
	}
	return str
}

var (
	ip string
	port string
	max int
	portRange [2]int
)

func init() {
	flag.StringVar(&ip, "ip", "", "IP to detect")
	flag.StringVar(&port, "port", "", "Ports to detect")
	flag.IntVar(&max, "max", 5000, "Max concurrent number")
	flag.Parse()
	match1, _ := regexp.MatchString("^\\d+-\\d+$", port)
	match2 := strings.Contains(port, ",")
	if match1 {
		index := strings.Index(port, "-")
		start, err := strconv.Atoi(string([]rune(port)[:index]))
		if err != nil {
			fmt.Println("Error: param port")
			return
		}
		end, err := strconv.Atoi(string([]rune(port)[index+1:]))
		if err != nil {
			fmt.Println("Error: param port")
			return
		}
		portRange = [2]int{start, end}
	} else if match2 {
		splits := strings.Split(port, ",")
		if len(splits) < 2 {
			fmt.Println("Error: param port")
			return
		}
		start, _ := strconv.Atoi(splits[0])
		end, _ := strconv.Atoi(splits[1])
		// TODO 这种输入格式时，检测输入的几个端口，而不是范围
		portRange = [2]int{start, end}
	} else {
		start, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("Error: param port")
			return
		}
		end, _ := strconv.Atoi(port)
		portRange = [2]int{start, end}
	}
}

func main() {
	//fmt.Println(ip, portRange, max)
	runtime.GOMAXPROCS(runtime.NumCPU())
	//start := time.Now().UnixNano()
	openPorts := scanRange(ip, portRange, max, 3e9)
	fmt.Println(join(openPorts))
	//fmt.Println((time.Now().UnixNano() - start) / 1e6, "ms")
}


