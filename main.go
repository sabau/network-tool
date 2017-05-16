package main

import (
	"fmt"
	"flag"
	"github.com/sabau/network-tool/networkTools"
	"log"
	"os/user"
	"os"
)

var (
	serverMode = flag.Bool("s", false, "run in server mode")
	udpPorts = []int{50000, 65535}
	portalPorts = []int{80, 443, 17992}
	routerPorts = []int{443, 17990}
	gatewayTcpPorts = []int{5060, 1720}
	allTcpPorts = []int{80, 443, 1720, 5060, 17990, 17992}
	iperfIp = "192.121.180.221"
	errorlog *os.File
	logger *log.Logger
)

func main(){
	fmt.Println("QuiVIDEO Network-Health tool")

	myself, err := user.Current()
	if err != nil {
		panic(err)
	}
	homedir := myself.HomeDir
	log_file := homedir+"/Desktop/network-tool.log"

	errorlog, err := os.OpenFile(log_file,  os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer errorlog.Close()

	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	logger = log.New(errorlog, "applog: ", log.Lshortfile|log.LstdFlags)

	logger.Print("QuiVIDEO Network-Health tool INIT\r\n")


	flag.Parse()
	if *serverMode == true {
		fmt.Println("QuiVIDEO Server Mode")
		for i := udpPorts[0]; i < udpPorts[1]; i++ {
			go func(i int) {
				networkTools.Server(i)
			} (i)
		}
	}else {

		var ip string
		fmt.Println("Vidyo UDP media range analysis")
		logger.Print("Vidyo UDP media range analysis:\r\n")
		networkTools.IperfCheck(iperfIp, udpPorts, logger)

		ip = "192.121.180.132"
		fmt.Printf("Portal: %s\n", ip)
		logger.Printf("Portal: %s\r\n", ip)
		networkTools.MachineCheck(ip, portalPorts, []int{}, logger)

		ip = "192.121.180.133"
		fmt.Printf("Router: %s\n", ip)
		logger.Printf("Router: %s\r\n", ip)
		networkTools.MachineCheck(ip, routerPorts, udpPorts, logger)

		ip = "192.121.180.142"
		fmt.Printf("Router: %s\n", ip)
		logger.Printf("Router: %s\r\n", ip)
		networkTools.MachineCheck(ip, routerPorts, udpPorts, logger)


		ip = "192.121.180.134"
		fmt.Printf("Gateway: %s\n", ip)
		logger.Printf("Gateway: %s\r\n", ip)
		networkTools.MachineCheck(ip, gatewayTcpPorts, udpPorts, logger)
	}
	var input string
	fmt.Println("\n\nPress enter to exit")
	fmt.Scanln(&input)
}

