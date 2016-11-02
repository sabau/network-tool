package networkTools

import (
	"fmt"
	"time"
	"github.com/sabau/port-scanner"
	"bufio"
	"net"
	//"bytes"
	"math"
	"sync"
)
var (
	timeoutDuration = 1 * time.Second
)

func MachineCheck(ip string, tcpPorts []int, udpPorts []int){
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
		for i := udpPorts[0]; i < udpPorts[1]; i++{
			go func(i int, ps *portscanner.PortScanner) {
				if !(ps.IsOpenUDP(i)){
					updClosed <- " * CLOSED [ICMP] " + string(i)
				} else {
					updClosed <- ""
				}
			}(i, ps)
		}
		msg := <- updClosed
		fmt.Println(msg)
	}
}

func IperfCheck(iperfIp string, udpPorts []int){
	var wg sync.WaitGroup
	//try with our server
	updClosed := make(chan int, udpPorts[1] - udpPorts[0])
	errors := make(chan string, udpPorts[1] - udpPorts[0])
	for i := udpPorts[0]; i < udpPorts[1]; i+=1000 {
		fmt.Printf("%d->%d", i, i + int(math.Min(float64(1000), float64(udpPorts[1]-i))))
		for j := i; (j < i+1000) && j < udpPorts[1]; j++ {
			wg.Add(1)
			go func(j int, iperfIp string, c chan int, e chan string) {
				defer wg.Done()
				ok, err := clientUDP(iperfIp, j)
				if len(err) > 0 {
					e <- err
				}
				if ! (ok) {
					c <- j
				}
			}(j, iperfIp, updClosed, errors)
		}
		wg.Wait()
		fmt.Println(" DONE")
	}
	close(updClosed)
	close(errors)
	for i := range updClosed {
		fmt.Printf("CLOSED PORT: %d\n",i)
	}

	for e := range errors {
		fmt.Println("Error description: " + e)
	}
}

func clientUDP(ip string, port int) (bool, string) {
	p :=  make([]byte, 2048)

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d",ip,port))
	if err  != nil {
		return false, fmt.Sprintf("Dial Error:  %v", err)
	}
	fmt.Fprintf(conn, "QuiVIDEO %s:%d", ip, port)
	reader := bufio.NewReader(conn)
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	_, err = reader.Read(p)
	if err == nil {
		//if (! bytes.Contains(p, []byte("QuiVIDEO"))) {
		//	fmt.Printf("%s\n", p)
		//}
	} else {
		return false, fmt.Sprintf("Receiving error %v\n", err)
	}
	conn.Close()
	return true, ""
}