package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"time"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
)


func buildAnnounceRequest(connectionId uint64,torrent *gotorrentparser.Torrent,p uint16) []byte{
	res:=make([]byte,98)

	// connectionID
	conId:=make([]byte,8)
	binary.BigEndian.PutUint64(conId,connectionId)
	copy(res[0:],conId)

	// action 
	action:=make([]byte,4)
	binary.BigEndian.PutUint32(action,1)
	copy(res[8:],action)

	// transactionID
	transactionID:=make([]byte,4)
	binary.BigEndian.PutUint32(transactionID,rand.Uint32())
	copy(res[12:],transactionID)

	// info Hash
	infoHash,_:=hex.DecodeString(torrent.InfoHash)
	copy(res[16:],infoHash)

	// peer_id

	copy(res[36:],PEER_ID)

	// key
	key := make([]byte, 4)
	rand.Read(key)
	copy(res[88:], key)

	// num_want
	num_want := make([]byte, 4)
	binary.BigEndian.PutUint32(num_want, 4294967295)
	copy(res[92:], num_want)

	// port
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, p)
	copy(res[96:], transactionID)

	return res
}

func buildConn()[]byte{

	// connectionId
	connectionId:=make([]byte,8)
	binary.BigEndian.PutUint64(connectionId,0x41727101980)

	// action 
	action:=make([]byte,4)
	binary.BigEndian.PutUint32(action,0)

	// TransactionId
	tld:=rand.Uint32()
	transactionId:=make([]byte,4)
	binary.BigEndian.PutUint32(transactionId,tld)

	buff:=connectionId
	buff = append(buff, action...)
	buff =append(buff, transactionId...)

	return buff

}

func parseAccounceResponse(resp []byte ,n int) annResp{
	var res annResp

	res.action=binary.BigEndian.Uint32(resp[0:4])
	res.transactionId=binary.BigEndian.Uint32(resp[4:8])
	res.leechers=binary.BigEndian.Uint32(resp[12:16])
	res.seeders=binary.BigEndian.Uint32(resp[8:12])
	
	temp:=resp[20:]

	for i:=0;i<(n-20);i+=6{
		if i+6<len(temp){
			var k Peer
			for j := i; j < i+4; j++ {
				k.ip += strconv.Itoa(int(temp[j]))
				if j < i+3 {
					k.ip += "."
				}
			}
			k.port = binary.BigEndian.Uint16(temp[i+4 : i+6])
			res.peers = append(res.peers, k)
		}
	}

	return res
}

func parseConnResp(resp []byte) connResp{
	var res connResp

	res.action =binary.BigEndian.Uint32(resp[0:4])
	res.transactionId=binary.BigEndian.Uint32(resp[4:8])
	res.connectionId=binary.BigEndian.Uint64(resp[8:16])

	return res
}

func handleConnection(k int, buff []byte, torrent *gotorrentparser.Torrent, peers *[]Peer){

	URL,err := url.Parse(torrent.Announce[k])
	if err!=nil {
		println("URL Parse failed:", err.Error())
		return
	}

	connection,err := net.Dial("udp",URL.Host)

	if err!=nil {
		println("Connection failed:", err.Error())
		return
	}

	defer connection.Close()

	err = connection.SetDeadline(time.Now().Add(5 * time.Second))
	if err!=nil {
		println("Connection SetDeadLine Error=", err.Error())
		return
	}

	connection.Write(buff)

	// buffer data we recieve

	recieved := make([]byte,16)

	_,err = connection.Read(recieved)
	
	if err!=nil {
		println("Unable to read Connect Data ", err.Error())
		return
	}

	resp:=parseConnResp(recieved)

	println("Connect Response Recieved")

	// connnect 
	if resp.action == 0 {
		req := buildAnnounceRequest(resp.connectionId, torrent, 6881)
		connection.Write(req)
		received := make([]byte, 1048576)
		n, err := connection.Read(received)
		println("Announce response size = ", n)
		if err != nil {
			println("Announce Read data failed:", err.Error())
		} else {
			resp := parseAccounceResponse(received, n)
			*peers = append(*peers, resp.peers...)
			if len(*peers) < len(resp.peers) {
				*peers = resp.peers
			}
		}
	}

}


func getPeers(torrent *gotorrentparser.Torrent)[]Peer{
	buff :=buildConn()
	urls :=torrent.Announce

	var peers []Peer

	for i:= range urls{
		if(urls[i][0:3]=="udp"){
			handleConnection(i,buff,torrent,&peers)
		}
	}

	newpeers := make([]Peer,0)

	for _,i:=range peers{
		peer:=i.ip+fmt.Sprintf("%v",i.port)
		if !listOfPeers[peer]{
			listOfPeers[peer]=true
			newpeers = append(newpeers, i)
		}
	}

	return newpeers
}