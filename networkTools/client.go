package networkTools

import (
	"fmt"
	"time"
	"github.com/sabau/port-scanner"
	"bufio"
	"net"
	"bytes"
)
var (
	timeoutDuration = 1 * time.Second
)

func MachineCheck(ip string, tcpPorts []int, iperfIp string, udpPorts []int){
	ps := portscanner.NewPortScanner(ip, timeoutDuration)
	fmt.Println("TCP Check, please wait.")
	for i := 0; i < len(tcpPorts); i++ {
		status := ""
		if ps.IsOpen(tcpPorts[i]) {
			status = "[Open]"
		}else{
			status = "[Closed]"
		}
		fmt.Println(" ", tcpPorts[i], " ", status, "  -->  ", ps.DescribePort(tcpPorts[i]))
	}

	if len(udpPorts) > 0 {
		updClosed := make(chan string)
		updOCount := make(chan int)
		updCCount := make(chan int)
		updOCount <- 0
		updCCount <- 0
		for i := udpPorts[0]; i < udpPorts[1]; i++{
			go func(i int, ps *portscanner.PortScanner) {
				if !(ps.IsOpenUDP(i)){
					updClosed <- " " + string(i)
					updCCount++
				} else {
					//try with our server
					if ! (clientUDP(iperfIp, i)){
						updClosed <- " IPERF " + string(i)
					}
				}
			}(i, ps)
		}
		msg := <- updClosed
		fmt.Println(msg)
	}
}

func clientUDP(ip string, port int) bool{
	p :=  make([]byte, 2048)

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d",ip,port))
	if err  != nil {
		fmt.Printf("Dial Error:  %v", err)
		return false
	}
	fmt.Fprintf(conn, "QuiVIDEO Check %s:%d", ip, port)
	reader := bufio.NewReader(conn)
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	_, err = reader.Read(p)
	if err == nil {
		if (! bytes.Contains(p, []byte("QuiVIDEO"))) {
			fmt.Printf("%s\n", p)
		}
	} else {
		fmt.Printf("Receiving error %v\n", err)
		return false
	}
	conn.Close()
	return true
}