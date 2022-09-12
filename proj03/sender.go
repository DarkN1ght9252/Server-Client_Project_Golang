package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func sender(filename *string, conn *net.UDPConn) int {
	var seqno uint16 = 0
	data, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Println("Error, cannot read file.")
		return 2
	}

	var pkt *Packet

	for start := 0; start < len(data); start += 255 {
		end := start + 255
		if end > len(data) {
			end = len(data)
		}
		fmt.Println("Update seqno")
		seqno++
		fmt.Println("Create Send Info")
		pkt = make_data_pkt(data[start:end], seqno)

		// TODO: send DATA and get ACK
		for {
			fmt.Println("Send Info to Rec")
			send(pkt, conn, nil)
			fmt.Println("Rcv response")
			retpkt, _, ok := recv(conn, 1)
			if !ok {
				return 3
			}
			if isACK(retpkt, seqno) {
				break
			}
		}
	}

	// TODO: send FIN and get FINACK
	pkt = make_fin_pkt(seqno)

	for {
		send(pkt, conn, nil)
		retpkt, _, ok := recv(conn, 1)
		if !ok {
			return 3
		}
		if isACK(retpkt, seqno+1) {
			pkt = make_ack_pkt(seqno + 1)
			send(pkt, conn, nil)
			break
		}
	}

	// TODO: return 0 for success, 3 for failure
	return 0
}

func make_data_pkt(data []byte, seqno uint16) *Packet {
	var pkt *Packet = &Packet{}

	pkt.hdr.flag = DATA
	pkt.hdr.seqno = seqno
	pkt.hdr.len = uint8(len(data))
	pkt.dat = data

	return pkt
}

func make_fin_pkt(seqno uint16) *Packet {
	var pkt *Packet = &Packet{}

	pkt.hdr.seqno = seqno
	pkt.hdr.flag = FIN

	return pkt
}

func isACK(pkt *Packet, expected uint16) bool {
	// TODO: return true if ACK (including FINACK) and ackno is what is expected
	fmt.Printf("%d : %d\n", pkt.hdr.ackno, expected)
	if (pkt.hdr.flag == FINACK || pkt.hdr.flag == ACK) && pkt.hdr.ackno == expected {
		return true
	}
	return false
}
