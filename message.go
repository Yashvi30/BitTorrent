package main

import (
	"encoding/binary"
	"encoding/hex"
)

func buildHandshake(infoHash string,peerId []byte)[]byte{
	res:=make([]byte,68)

	//pstrlen
	copy(res[0:],[]byte{19})

	// pstr
	copy(res[1:20],"BitTorrent protocol")

	// reserved bytes
	copy(res[20:],[]byte{0,0,0,0,0,0,0,0})

	// infoHash
	info_hash,_:=hex.DecodeString(infoHash)
	copy(res[28:],info_hash)

	//peerId

	copy(res[48:],peerId)

	return res
}

func keepAlive(peerConn *PeerConnection)error{
	res:=make([]byte,4)

	copy(res[0:],[]byte{0,0,0,0})

	_,err:=peerConn.conn.Write(res)

	return err
}

func sendChoke(peerConn *PeerConnection)error{
	res:=make([]byte,5)

	copy(res[0:4],[]byte{0,0,0,1})

	copy(res[4:],[]byte{0})

	_,err:=peerConn.conn.Write(res)

	return err
}

func sendUnchoke(peerConn *PeerConnection)error{
	res:=make([]byte,5)

	copy(res[0:4],[]byte{0,0,0,1})

	copy(res[4:],[]byte{1})

	_,err:=peerConn.conn.Write(res)

	return err
}

func sendInterested(peerConn *PeerConnection)error{
	res:=make([]byte,5)

	copy(res[0:],[]byte{0,0,0,1})

	copy(res[4:],[]byte{2})

	_,err:=peerConn.conn.Write(res)

	return err
}

func notInterested(peerConn *PeerConnection)error{
	res:=make([]byte,5)

	copy(res[0:],[]byte{0,0,0,1})

	copy(res[4:],[]byte{3})

	_,err:=peerConn.conn.Write(res)

	return err
}

func sendHave(peerConn *PeerConnection,index uint32)error{
	res:=make([]byte,9)

	copy(res[0:4],[]byte{0,0,0,5})

	copy(res[4:],[]byte{4})

	indexBit:=make([]byte,4)
	binary.BigEndian.PutUint32(indexBit,index)
	copy(res[5:],indexBit)

	_,err:=peerConn.conn.Write(res)

	return err
}

func sendBitfield(peerConn *PeerConnection,bitfield []byte)error{
	res:=make([]byte,5+len(bitfield))

	length:=make([]byte,4)
	binary.BigEndian.PutUint32(length,uint32(len(bitfield)+1))
	copy(res[0:],length)

	copy(res[4:],[]byte{5})

	copy(res[5:],bitfield);

	_,err:=peerConn.conn.Write(res)

	return err
}

func sendRequest(peerConn *PeerConnection,index uint32,offset uint32,length uint32) error{
	res:=make([]byte,17)

	copy(res[0:],[]byte{0,0,0,13})

	copy(res[4:],[]byte{6})

	indexBit:=make([]byte,4)
	binary.BigEndian.PutUint32(indexBit,index)
	copy(res[5:],indexBit)

	off_set:=make([]byte,4)
	binary.BigEndian.PutUint32(off_set,offset)
	copy(res[9:],off_set)

	leng:=make([]byte,4)
	binary.BigEndian.PutUint32(leng,length)
	copy(res[13:],leng)

	_, err := peerConn.conn.Write(res)
	return err
}

func sendPiece(peerConn *PeerConnection,index uint32,offset uint32,block []byte)error{
	res:=make([]byte,13+len(block))

	copy(res[0:],[]byte{0,0,0,9})

	copy(res[4:],[]byte{7})

	indexBit:=make([]byte,4)
	binary.BigEndian.PutUint32(indexBit,index)
	copy(res[5:],indexBit)

	off_set:=make([]byte,4)
	binary.BigEndian.PutUint32(off_set,offset)
	copy(res[9:],off_set)

	copy(res[13:],block)

	_, err := peerConn.conn.Write(res)
	return err
}

func sendCancel(peerConn *PeerConnection,index uint32,offset uint32,length uint32) error{
	res:=make([]byte,17)

	copy(res[0:],[]byte{0,0,0,13})

	copy(res[4:],[]byte{8})

	indexBit:=make([]byte,4)
	binary.BigEndian.PutUint32(indexBit,index)
	copy(res[5:],indexBit)

	off_set:=make([]byte,4)
	binary.BigEndian.PutUint32(off_set,offset)
	copy(res[9:],off_set)

	leng:=make([]byte,4)
	binary.BigEndian.PutUint32(leng,length)
	copy(res[13:],leng)

	_, err := peerConn.conn.Write(res)
	return err
}

func sendPort(peerConn *PeerConnection,port uint16)error{

	res:=make([]byte,7)

	copy(res[0:],[]byte{0,0,0,3})

	copy(res[4:],[]byte{9})

	Port:=make([]byte,2)
	binary.BigEndian.PutUint16(Port,port)
	copy(res[5:],Port)

	_, err := peerConn.conn.Write(res)
	return err
}