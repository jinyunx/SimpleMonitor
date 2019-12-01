package main

import (
	"encoding/json"
	"log"
	"net"
)

type attrData struct {
	Instance string `json:"instance"`
	Time     int64  `json:"time"`
	Attr     int    `json:"attr"`
	Counter  int    `json:"counter"`
}

func udpServer(d *db, port string) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		data := make([]byte, 1500)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println("Read udp data err,", err)
			continue
		}

		go handleUdpData(d, data[:n], remoteAddr)
	}
}

func handleUdpData(d *db, data []byte, remoteAddr *net.UDPAddr) {
	log.Println("get data:", string(data))
	a := attrData{}
	if err := json.Unmarshal(data, &a); err != nil {
		log.Println("Json data err,", err, " ,data:", string(data))
		return
	}

	if a.Instance == "" {
		a.Instance = string(remoteAddr.IP)
	}

	_ = d.updateInstanceAttr(a.Instance, a.Time - (a.Time % 60), a.Attr, a.Counter)
	_ = d.updateInstanceAttr(a.Instance, a.Time - (a.Time % 60) + 60, a.Attr, 0)
}
