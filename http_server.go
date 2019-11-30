package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func httpServer(d *db, port string) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))
	mux.HandleFunc("/r", func(writer http.ResponseWriter, request *http.Request) {
		handleDataRequest(d, writer, request)
	})
	mux.HandleFunc("/attr", func(writer http.ResponseWriter, request *http.Request) {
		handleAttrRequest(d, writer, request)
	})
	err := http.ListenAndServe(":" + port, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func getQuery(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func getQueryInt(r *http.Request, key string) int {
	i, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		return 0
	}
	return i
}

func handleDataRequest(d *db, w http.ResponseWriter, r *http.Request) {
	instance := getQuery(r, "instance")
	if len(instance) == 0 {
		fmt.Fprint(w, "Instance cannot be emty")
	}

	attr := getQueryInt(r, "attr")
	if attr == 0 {
		fmt.Fprint(w, "Attr cannot be 0")
	}

	now := time.Now().Unix()
	date := now - now % (3600 * 24)
	tableName := time.Unix(date, 0).Format("2006_01_02")
	result := d.queryByInstaceAndAttr(tableName, instance, attr)
	resultJson, _ := json.Marshal(&result)
	fmt.Fprint(w, string(resultJson))
}

func handleAttrRequest(d *db, w http.ResponseWriter, r *http.Request) {
	instance := getQuery(r, "instance")
	if len(instance) == 0 {
		fmt.Fprint(w, "Instance cannot be emty")
	}

	now := time.Now().Unix()
	date := now - now % (3600 * 24)
	tableName := time.Unix(date, 0).Format("2006_01_02")
	result := d.queryByInstace(tableName, instance)
	resultJson, _ := json.Marshal(&result)
	fmt.Fprint(w, string(resultJson))
}
