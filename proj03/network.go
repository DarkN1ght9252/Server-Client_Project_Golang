package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func send(pkt *Packet, conn *net.UDPConn, addr net.Addr) (int, error) {
	buffer := new(bytes.Buffer)
	buffer.Grow(int(pkt.hdr.len) + 6)

	// pack the data
	var lenflag uint16
	lenflag = uint16(pkt.hdr.len)
	lenflag = lenflag << 8
	lenflag += uint16(pkt.hdr.flag)
	binary.Write(buffer, binary.BigEndian, lenflag)
	binary.Write(buffer, binary.BigEndian, pkt.hdr.seqno)
	binary.Write(buffer, binary.BigEndian, pkt.hdr.ackno)
	binary.Write(buffer, binary.BigEndian, pkt.dat)

	var n int
	var err error
	if addr == nil {
		// sender, addr not needed
		n, err = conn.Write(buffer.Bytes())
	} else {
		// receiver, need to know who we are replying to
		n, err = conn.WriteTo(buffer.Bytes(), addr)
	}
	return n, err
}

func recv(conn *net.UDPConn, timeout int) (*Packet, net.Addr, bool) {
	var ack bool = false
	var success bool = true
	var addr net.Addr
	if timeout > 0 {
		ack = true
	}
	var buffer []byte = make([]byte, MAX_PACKET_SIZE)
	var pkt *Packet = &Packet{}

	if ack {
		// read with timer (sender)
		deadline := time.Now().Add(time.Second)
		conn.SetReadDeadline(deadline)
		_, _, err := conn.ReadFrom(buffer)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				// we timed out...
				fmt.Println("Warning, timeout on ack")
				return nil, nil, !success
			}
			fmt.Println("Error reading")
			return nil, nil, !success
		}
	} else {
		// read without timer (receiver)
		var err error
		_, addr, err = conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error reading")
			return nil, nil, !success
		}
	}

	// unpack the bytes
	reader := bytes.NewBuffer(buffer)
	pkt.hdr.len, _ = reader.ReadByte()
	pkt.hdr.flag, _ = reader.ReadByte()
	seqnohigh, _ := reader.ReadByte()
	seqnolow, _ := reader.ReadByte()
	pkt.hdr.seqno = (uint16(seqnohigh) << 8) + uint16(seqnolow)
	acknohigh, _ := reader.ReadByte()
	acknolow, _ := reader.ReadByte()
	pkt.hdr.ackno = (uint16(acknohigh) << 8) + uint16(acknolow)
	pkt.dat = buffer[6:]

	return pkt, addr, success
}
