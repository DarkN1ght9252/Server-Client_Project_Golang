package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// Our cirdt header
/*
| 1 byte | 1 byte |     2 bytes    |     2 bytes    |
| length | flags  |  sequence num. |    ack. num.   |
*/
type Header struct {
	len   uint8
	flag  uint8
	seqno uint16
	ackno uint16
}

// A packet is a header + data (encapsulation)
type Packet struct {
	hdr Header
	dat []byte
}

// Max possible size is 255 + 6, why?
var MAX_PACKET_SIZE int = 255 + 6

// flag values
var DATA uint8 = 0
var ACK uint8 = 1
var FIN uint8 = 2
var FINACK uint8 = FIN | ACK

func main() {
	var sflag = flag.Bool("sender", false, "operate as a sender")
	var rflag = flag.Bool("receiver", false, "operate as a receiver")
	var filename = flag.String("file", "", "filename to transfer")
	flag.Parse()

	if len(*filename) == 0 {
		fmt.Println("Error, must specify a file to read/write to.")
		os.Exit(1)
	}

	if *sflag && *rflag {
		fmt.Println("Error, cannot operate as a sender and receiver at the same time.")
		os.Exit(1)
	} else if *sflag {
		// sender time!
		fmt.Println("send it!")

		raddr, err := net.ResolveUDPAddr("udp", "localhost:9001")
		if err != nil {
			fmt.Println("Error, could not resolve IP address.")
			os.Exit(1)
		}

		conn, err := net.DialUDP("udp", nil, raddr)
		if err != nil {
			fmt.Println("Error, could not Dial UDP.")
			os.Exit(1)
		}
		defer conn.Close()

		ret := sender(filename, conn)
		if ret != 0 {
			os.Exit(ret)
		} else {
			fmt.Println("Sent successfully.")
		}

	} else if *rflag {
		// receiver time!
		fmt.Println("recv it!")

		laddr, err := net.ResolveUDPAddr("udp", "localhost:9001")
		if err != nil {
			fmt.Println("Error, could not resolve IP.")
			os.Exit(1)
		}

		conn, err := net.ListenUDP("udp", laddr)
		if err != nil {
			fmt.Println("Error, could not listen.")
			os.Exit(1)
		}
		defer conn.Close()

		ret := receiver(filename, conn)
		if ret != 0 {
			os.Exit(ret)
		} else {
			fmt.Println("Recv successfully.")
		}

	} else {
		fmt.Println("Error, must specify sender or receiver.")
		os.Exit(1)
	}
}
