package main

import (
	"bitmessage"
	"bufio"
	"fmt"
	"net"
)

//https://en.bitcoin.it/wiki/Satoshi_Client_Node_Discovery

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
	ExtAddr = s[1 : len(s)-2]
	fmt.Println(ExtAddr)

	//get peer
	//	addrs, err := net.LookupIP("dnsseed.bitcoin.dashjr.org")
	//	if err != nil {
	//		fmt.Println("Error looking up bitcoin.sipa.be:", err)
	//	}
	//	fmt.Println(addrs)
	//dns is dead, use preseeded address (really evil)
	list := bitmessage.AddressList{0, make([]*bitmessage.NetworkAddress, 0, 30)}
	peer := &bitmessage.NetworkAddress{}
	peer.SetTimestampNow(0)
	peer.SetIP(net.IPv4(0xc0, 0xa8, 0x01, 0x2b)) //local daemon	
	peer.SetPort(8333)
	list.AddAddress(peer)
	fmt.Println(peer.IP)

	//connect to peer
	bm := new(bitmessage.BitMessage) //todo: make more idiomatic (golang)
	bm.SetMagic(bitmessage.MainNetMagic)
	bm.SetCommand("version")

	vm := new(bitmessage.VersionMessage)
	vm.SetVersion(40000)
	vm.SetServices(bitmessage.NODE_NETWORK)
	vm.SetTimestampNow()
	vm.SetAddrYou(net.IPv4(0x00, 0x00, 0x00, 0x00), 1, 0)
	vm.SetAddrMe(net.IPv4(0x00, 0x00, 0x00, 0x00), 1, 0)
	vm.SetRandomNonce()
	vm.SetSubVersionNull()
	vm.SetStartHeight(uint32(1))

	bm.SetPayloadVersion(vm)

	compiled := bm.Compile()

	//send version and connect to peer
	conn, err = net.Dial("tcp", "192.168.1.43:8333")
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
	}
	fmt.Println(compiled)
	conn.Write(compiled)
	fmt.Println("Reading")
	msgChan := make(chan *bitmessage.BitMessage)
	go bitmessage.DecodeMessages(conn, msgChan)
	for bm := range msgChan {
		fmt.Println(bm.GiveMessageType())
	}
	reader = bufio.NewReader(conn)
	s, err = reader.ReadString('k')
	if err != nil {
		fmt.Println("Error reading response:", err)
	}
	fmt.Println(s)
}
