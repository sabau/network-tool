package main

import (
    "fmt"
    "time"
    "github.com/anvie/port-scanner"
)

func main(){
     t := time.Duration(1)*time.Second
     startPort := 1720
     endPort := 1721
     ps := portscanner.NewPortScanner("140.242.46.49", t)

     fmt.Printf("Checking your device connectivity %d-%d...\n", startPort, endPort)
     openedPorts := ps.GetOpenedPort(startPort, endPort)
     for i := 0; i < len(openedPorts); i++ {
        port := openedPorts[i]
        fmt.Print(" ", port, " [open]")
        fmt.Println("  -->  ", ps.DescribePort(port))
     }

     ps = portscanner.NewPortScanner("81.174.21.18", t)

     fmt.Printf("Checking destination reachability %d-%d...\n", startPort, endPort)
     openedPorts = ps.GetOpenedPort(startPort, endPort)

     for i := 0; i < len(openedPorts); i++ {
        port := openedPorts[i]
        fmt.Print(" ", port, " [open]")
        fmt.Println("  -->  ", ps.DescribePort(port))
     }

}

