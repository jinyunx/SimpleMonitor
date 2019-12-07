package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	go logRotate()
}

func getLogName() string {
	now := time.Now()
	return "log/monitor_" + now.Format("2006-01-02") + ".log"
}

func logRotate() {
	if err := os.MkdirAll("log", os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	lastLogName := getLogName()
	logFile, err := os.OpenFile(lastLogName,
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Log file open error : %v", err)
	}

	log.SetOutput(logFile)

	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		if lastLogName != getLogName() {
			lastLogName = getLogName()
			logFile.Close()
			logFile, err = os.OpenFile(lastLogName,
				os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

			if err != nil {
				log.Fatalf("Log file open error : %v", err)
			}

			log.SetOutput(logFile)
		}
	}
}

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
		fmt.Println("send data")

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
