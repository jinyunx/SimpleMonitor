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
	err := http.ListenAndServe(":"+port, mux)
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
	view := getQuery(r, "view")
	if len(instance) == 0 && len(view) == 0 {
		fmt.Fprint(w, "Instance and view cannot both be emty")
		return
	}

	attr := getQueryInt(r, "attr")
	if attr == 0 {
		fmt.Fprint(w, "Attr cannot be 0")
		return
	}

	key := view
	f := d.queryByViewAndAttr
	if len(view) == 0 {
		key = instance
		f = d.queryByInstaceAndAttr
	}

	now := time.Now()
	t := now.Unix()
	today := t - t%(3600*24)
	t = now.Add(-time.Hour * 24).Unix()
	yestoday := t - t%(3600*24)
	t = now.Add(-time.Hour * 24 * 7).Unix()
	lastWeek := t - t%(3600*24)

	type threeDayAttr struct {
		Today    attrCounterArr `json:"today"`
		Yestoday attrCounterArr `json:"yestoday"`
		LastWeek attrCounterArr `json:"last_week"`
	}

	threeDay := threeDayAttr{}

	tableName := time.Unix(today, 0).Format("2006_01_02")
	threeDay.Today = f(tableName, key, attr)

	tableName = time.Unix(yestoday, 0).Format("2006_01_02")
	threeDay.Yestoday = f(tableName, key, attr)

	tableName = time.Unix(lastWeek, 0).Format("2006_01_02")
	threeDay.LastWeek = f(tableName, key, attr)

	resultJson, _ := json.Marshal(&threeDay)
	fmt.Fprint(w, string(resultJson))
}

func handleAttrRequest(d *db, w http.ResponseWriter, r *http.Request) {
	instance := getQuery(r, "instance")
	view := getQuery(r, "view")
	if len(instance) == 0 && len(view) == 0 {
		fmt.Fprint(w, "Instance and view cannot both be emty")
		return
	}

	now := time.Now().Unix()
	date := now - now%(3600*24)
	tableName := time.Unix(date, 0).Format("2006_01_02")

	result := attrArr{}
	if len(view) != 0 {
		result = d.queryAttrByView(tableName, view)
	} else {
		result = d.queryAttrByInstance(tableName, instance)
	}
	resultJson, _ := json.Marshal(&result)
	fmt.Fprint(w, string(resultJson))
}
