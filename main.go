package main

import (
	"fmt"
	"flag"
	"github.com/sabau/network-tool/networkTools"
)

var (
	serverMode = flag.Bool("s", false, "run in server mode")
	udpPorts = []int{50000, 65535}
	portalPorts = []int{80, 443, 17992}
	routerPorts = []int{443, 17990}
	gatewayTcpPorts = []int{5060, 1720}
	allTcpPorts = []int{80, 443, 1720, 5060, 17990, 17992}
	iperfIp = "192.121.180.221"
)

func main(){
	fmt.Println("QuiVIDEO Network-Health tool")
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
		networkTools.IperfCheck(iperfIp, udpPorts)

		ip = "192.121.180.132"
		fmt.Printf("Portal: %s\n", ip)
		networkTools.MachineCheck(ip, portalPorts, []int{})

		ip = "192.121.180.133"
		fmt.Printf("Router: %s\n", ip)
		networkTools.MachineCheck(ip, routerPorts, udpPorts)

		ip = "192.121.180.142"
		fmt.Printf("Router: %s\n", ip)
		networkTools.MachineCheck(ip, routerPorts, udpPorts)


		ip = "192.121.180.134"
		fmt.Printf("Gateway: %s\n", ip)
		networkTools.MachineCheck(ip, gatewayTcpPorts, udpPorts)
	}
	var input string
	fmt.Println("\n\nPress enter to exit")
	fmt.Scanln(&input)
}

