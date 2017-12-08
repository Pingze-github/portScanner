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
	"sort"
	"sync"
)

func scanRange(ip string, ports []int, gomax int, timeout int) ([]int, []int) {

	var mutex sync.Mutex

	finish := make(chan int)
	channel := make(chan int, gomax)
	var openPorts []int
	var timeoutPorts []int

	scan := func (ip string, port int) {
		// time.Sleep(time.Duration(timeout)) // 模拟超时
		address := ip + ":" + strconv.Itoa(port)
		_, err := net.DialTimeout("tcp", address, time.Duration(timeout))
		if err != nil {
			//fmt.Println(err)
			if strings.Contains(err.Error(), "timeout") {
				//fmt.Println("TIMEOUT")
				timeoutPorts = append(timeoutPorts, port)
			}
		} else {
			mutex.Lock()
			openPorts = append(openPorts, port)
			defer mutex.Unlock()
		}
		i := <- channel
		if i == 1 {
			finish <- 0
		}
	}
	num := len(ports);
	for i := 0; i < num; i++ {
		if i == num - 1 {
			channel <- 1
		} else {
			channel <- 0
		}
		go scan(ip, ports[i])
	}
	<- finish
/*	fmt.Println(len(refusedPorts))
	fmt.Println("拒绝", refusedPorts)
	fmt.Println(len(timeoutPorts))
	fmt.Println("超时", timeoutPorts)
	fmt.Println(len(openPorts))*/
	return openPorts, timeoutPorts
}

func join(intArr []int) string {
	sort.Sort(sort.IntSlice(intArr))
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
	timeout int
	ports []int
)

func init() {
	flag.StringVar(&ip, "ip", "", "IP to detect")
	flag.StringVar(&port, "port", "", "Ports to detect")
	flag.IntVar(&max, "max", 5000, "Max concurrent number")
	flag.IntVar(&timeout, "timeout", 3000, "connect timeout (ms)")
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
		portStrings := strings.Split(port, ",")
		for i := 0; i < len(portStrings); i++ {
			portInt, _ := strconv.Atoi(portStrings[i])
			ports = append(ports, portInt)
		}
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
	if len(ports) == 0 {
		for port := portRange[0]; port <= portRange[1]; port++ {
			ports = append(ports, port)
		}
	}
	openPorts, timeoutPorts := scanRange(ip, ports, max, timeout * 1e6)
	if len(timeoutPorts) > (portRange[1] - portRange[0]) * 9 / 10 {
		fmt.Printf("Error: timeout")
	} else {
		fmt.Printf(join(openPorts))
	}
	//fmt.Println((time.Now().UnixNano() - start) / 1e6, "ms")
}


