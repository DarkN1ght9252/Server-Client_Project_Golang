package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
)

func receiver(filename *string, conn *net.UDPConn) int {
	var expected uint16 = 0
	var pkt *Packet
	var buffer []byte
	reader := bytes.NewBuffer(buffer)
	// recieve
	for {
		// TODO: receive DATA and send ACK if exepcted seqno arrives
		// NOTE: Don't forget to write the data
		// NOTE: You'll need the addr returned from recv in order to
		// send back to the sender.
		fmt.Println("Get Send Info")
		rcv, addr, ok := recv(conn, 0)

		if !ok {
			return 3
		}

		//Writes Data to File
		if rcv.hdr.seqno > expected {
			fmt.Println("Update Expected")
			expected++
			reader.Write(rcv.dat)
			pkt = make_ack_pkt(expected)
		}
		fmt.Println("Send ACK")
		// if rand.Intn(10) < 2 {
		send(pkt, conn, addr)
		// }

		// TODO: break out of infinte loop after FINACK
		if rcv.hdr.flag == FIN {
			for {
				pkt = make_finack_pkt(expected + 1)

				rcv, addr, ok := recv(conn, 0)

				if !ok {
					return 3
				}

				if rcv.hdr.flag == ACK {

					break
				} else {
					send(pkt, conn, addr)
				}
			}
			ioutil.WriteFile(*filename, reader.Bytes(), 0666)
			break
		}
	}

	return 0
}

func make_ack_pkt(ackno uint16) *Packet {
	var pkt *Packet = &Packet{}

	pkt.hdr.ackno = ackno
	pkt.hdr.flag = ACK

	return pkt
}

func make_finack_pkt(ackno uint16) *Packet {
	var pkt *Packet = &Packet{}

	pkt.hdr.flag = FINACK
	pkt.hdr.ackno = ackno

	return pkt
}
