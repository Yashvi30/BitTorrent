package main

import (
	"encoding/binary"
	"io"
	"net"
	"time"
)

func messageType(peerConn *PeerConnection, t int) (int32, int32, error) {
	peerConn.conn.SetDeadline(time.Now().Add(time.Duration(t) * time.Second))
	defer peerConn.conn.SetDeadline(time.Time{})

	buff1 := make([]byte, 4)
	_, err := io.ReadFull(peerConn.conn, buff1)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return 0, -2, nil
		}
		println(err.Error())
		return -1, -1, err
	}

	msglen := int32(binary.BigEndian.Uint32(buff1))
	if msglen == 0 {
		return 0, -1, nil
	}

	buff2 := make([]byte, 1)
	_, err = io.ReadFull(peerConn.conn, buff2)
	if err != nil {
		return -1, -1, err
	}

	msgId := int32(uint32(buff2[0]))

	return msglen,msgId,nil

}



func handleHave(peerConn *PeerConnection,length int32)error{
	peerConn.conn.SetDeadline(time.Now().Add(100 * time.Second))
	defer peerConn.conn.SetDeadline(time.Time{})

	buff:=make([]byte,length)
	_, err := io.ReadFull(peerConn.conn, buff)
	if err != nil {
		return err
	}

	index := int32(binary.BigEndian.Uint32(buff))
	(*peerConn.bitfield)[index]=true

	return nil
}


func handleBitfield(peerConn *PeerConnection,length int32) error{
	peerConn.conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer peerConn.conn.SetDeadline(time.Time{})

	buff:=make([]byte,length)
	_,err:=io.ReadFull(peerConn.conn,buff)
	if err!=nil{
		return err
	}

	for i, j := range buff {
		for bit := 0; bit < 8; bit++ {
			if (j&(1<<bit) != 0) && ((i+1)*8-bit-1 < len(*peerConn.bitfield)) {
				(*peerConn.bitfield)[(i+1)*8-bit-1] = true
			}
		}
	}
	return nil
}

func handleCancel(peerConnection *PeerConnection) {
	// TODO
}
func handlePort(peerConnection *PeerConnection) {
	// TODO
}
func handleRequest(peerConnection *PeerConnection) {
	// TODO
}

func handdlePiece(peerConn *PeerConnection,length int32) error{
	peerConn.conn.SetDeadline(time.Now().Add(50 * time.Second))
	defer peerConn.conn.SetDeadline(time.Time{})

	buff:= make([]byte,length)
	_,err:=io.ReadFull(peerConn.conn,buff)
	if err!=nil{
		return err
	}

	index:=int32(binary.BigEndian.Uint32(buff[0:4]))
	offset:=int32(binary.BigEndian.Uint32(buff[4:8]))

	if pieces[index].data==nil{
		return nil
	}

	mutex.Lock()
	copy((*pieces[index].data)[offset:], buff[8:])
	mutex.Unlock()

	return nil
}

func handleMessage(peerConn *PeerConnection,msgId,msglen int32)error{
	switch msgId {
	case -2:
		// timeout
	case -1:
		// keep alive
		return nil
	case 0:
		// choke
		peerConn.choked = true
	case 1:
		// unchoke
		peerConn.choked = false
	case 2:
		// interested
		peerConn.interested = true
	case 3:
		// not interested
		peerConn.interested = false
	case 4:
		// have
		return handleHave(peerConn, msglen-1)
	case 5:
		// bitfield
		return handleBitfield(peerConn, msglen-1)
	case 6:
		// request
		// TODO
	case 7:
		// piece
		return handdlePiece(peerConn, msglen-1)
	case 8:
		// cancel
		// TODO
	case 9:
		// port
		// TODO
	default:
		println("Unknown message id: ", msgId, " with length: ", msglen, "to the peer: ", peerConn.peer.ip, ":", peerConn.peer.port)
		peerConn.conn.SetDeadline(time.Now().Add(100 * time.Second))
		defer peerConn.conn.SetDeadline(time.Time{})
		buff := make([]byte, msglen-1)
		_, err := io.ReadFull(peerConn.conn, buff)
		return err
	}

	return nil
}
