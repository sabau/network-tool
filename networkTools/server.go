package networkTools


import (
	"fmt"
	"net"
	"os"
	"bytes"
)


func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_,err := conn.WriteToUDP([]byte(fmt.Sprintf("QuiVIDEO: %d\n", addr.Port)), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

/* A Simple function to verify error */
func checkError(err error) {
	if err  != nil {
		fmt.Printf("Some error  %v", err)
		os.Exit(0)
	}
}


func Server(port int) {
	p := make([]byte, 2048)

	/* Lets prepare a address at any address at port 10001*/
	ServerAddr,err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%d",port))
	//checkError(err)
	if err  != nil {
		fmt.Printf("Resolve Error  %v\n\n", err)
		return
	}

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	//checkError(err)
	if err  != nil {
		fmt.Printf("Listen Error  %v\n\n", err)
		return
	}
	defer ServerConn.Close()


	for {
		_,remoteaddr,err := ServerConn.ReadFromUDP(p)
		if err !=  nil {
			fmt.Printf("Read message error  %v\n\n", err)
			continue
		}
		if (bytes.Contains(p, []byte("QuiVIDEO"))) {
			fmt.Printf("%s \n", p)
			go sendResponse(ServerConn, remoteaddr)
		}
	}
}

