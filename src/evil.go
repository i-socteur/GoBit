package main

import (
//	"bitmessage"
	"bufio"
	"fmt"
	"net"
)

func main() {
	//get external address	
	var ExtAddr string
	conn, err := net.Dial("tcp", "checkip.dyndns.org:80")
	if err != nil {
		fmt.Println("Error getting extAddr:", err)
	}

	fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: checkip.dyndns.org\r\nUser-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)\r\nConnection: close\r\n\r\n")
	reader := bufio.NewReader(conn)
	reader.ReadString('<')
	reader.ReadString(':')
	s, _ := reader.ReadString('<')
	ExtAddr = s[1:len(s)-2]
	fmt.Println(ExtAddr)
	

	//get peer

}
