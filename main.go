package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func main() {
	d := db{}
	if err := d.open("monitor.db"); err != nil {
		log.Fatal(err)
	}

	go httpServer(&d, "8080")
	go udpServer(&d, "34218")

	testClient()
}

func testClient() {
	//connect server
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 34218,
	})

	if err != nil {
		fmt.Printf("connect failed, err: %v\n", err)
		return
	}

	for {
		log.Println("send data")

		a := attrData{
			Instance: "instance"+strconv.Itoa(rand.Int()%2+1),
			Time:     time.Now().Unix(),
			Attr:     rand.Int() % 2 + 1,
			Counter:  rand.Int() % 1000,
		}

		d, _ := json.Marshal(&a)
		_, err = conn.Write(d)
		if err != nil {
			fmt.Printf("send data failed, err : %v\n", err)
		}

		time.Sleep(time.Second)
	}
}
