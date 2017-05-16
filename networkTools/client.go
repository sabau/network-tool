package networkTools

import (
	"fmt"
	"time"
	"github.com/sabau/port-scanner"
	"bufio"
	"net"
	"math"
	"sync"
	"log"
	"strconv"
	"strings"
)
var (
	timeoutDuration = 5 * time.Second
)

func MachineCheck(ip string, tcpPorts []int, udpPorts []int, logger *log.Logger){
	ps := portscanner.NewPortScanner(ip, timeoutDuration)
	fmt.Println("TCP Check, please wait.")
	logger.Print("TCP Check, please wait.\r\n")
	for i := 0; i < len(tcpPorts); i++ {
		status := ""
		if ps.IsOpen(tcpPorts[i]) {
			status = "[Open]"
		}else{
			status = "[Closed]"
		}
		fmt.Println(" ", tcpPorts[i], " ", status, "  -->  ", ps.DescribePort(tcpPorts[i]))
		logger.Print(" ", tcpPorts[i], " ", status, "  -->  ", ps.DescribePort(tcpPorts[i]), "\r\n")
	}

	if len(udpPorts) > 0 {
		updClosed := make(chan int, udpPorts[1] - udpPorts[0])
		var wg sync.WaitGroup
		for i := udpPorts[0]; i < udpPorts[1]; i++{
			wg.Add(1)
			go func(i int, ps *portscanner.PortScanner, c chan int) {
				defer wg.Done()
				if ! (ps.IsOpenUDP(i)) {
					c <- i
				}
			} (i, ps, updClosed)
		}
		wg.Wait()
		close(updClosed)
		for i := range updClosed {
			fmt.Printf("ICMP CLOSED PORT: %d\n",i)
			logger.Printf("ICMP CLOSED PORT: %d\r\n",i)
		}
	}
}

func IperfCheck(iperfIp string, udpPorts []int, logger *log.Logger){
	var wg sync.WaitGroup
	//try with our server
	updClosed := make(chan int, udpPorts[1] - udpPorts[0])
	errors := make(chan string, udpPorts[1] - udpPorts[0])
	exp_increment := 1000
	step := 250
	fmt.Printf("%d->", udpPorts[0])
	for i := udpPorts[0]; i < udpPorts[1]; i+=step {
		init_length := len(updClosed)
		fmt.Printf("%d", i + int(math.Min(float64(step), float64(udpPorts[1]-i))))
		for j := i; (j < i+step) && j < udpPorts[1]; j++ {
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
			} (j, iperfIp, updClosed, errors)
		}
		wg.Wait()
		fmt.Print("->")
		time.Sleep(500 * time.Millisecond)
		if len(updClosed) >= init_length + step && i + step + exp_increment < udpPorts[1] {
			i += exp_increment
			exp_increment += exp_increment
		}
	}
	fmt.Println(" UDP Connectivity check DONE")
	close(updClosed)
	close(errors)
	fmt.Print("CLOSED PORTS:")
	logger.Print("CLOSED PORT: \r\n")
	s := []string{}
	for i := range updClosed {
		fmt.Printf(" %d",i)
		s = append(s, strconv.Itoa(i))
	}
	logger.Print(strings.Join(s, " "), "\r\n")
	fmt.Println("")

	//for e := range errors {
	//	fmt.Println("Error description: " + e)
	//}
}

func clientUDP(ip string, port int) (bool, string) {
	p :=  make([]byte, 2048)

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d",ip,port))
	if err  != nil {
		return false, fmt.Sprintf("Dial Error:  %v", err)
	}
	fmt.Fprintf(conn, "QuiVIDEO %s:%d", ip, port)
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	reader := bufio.NewReader(conn)
	_, err = reader.Read(p)

	if err == nil {
		//if (! bytes.Contains(p, []byte("QuiVIDEO"))) {
		//	fmt.Printf("%s\n", p)
		//}
	} else {

		return false, fmt.Sprintf("Receiving error %v", err)
	}
	conn.Close()
	return true, ""
}