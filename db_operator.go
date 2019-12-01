package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

const (
	kCreateAttrInfo = `CREATE TABLE IF NOT EXISTS "attr_info" (
                            "attr" INTEGER NOT NULL PRIMARY KEY,
                            "name" TEXT)`

	kCreateViewInfo = `CREATE TABLE IF NOT EXISTS "view_info" (
                            "view" TEXT, 
                            "instance" TEXT,
                            PRIMARY KEY(view, instance))`

	kCreateInstanceTable = `CREATE TABLE IF NOT EXISTS "%s" (
								"instance" TEXT,
								"time" INTEGER,
								"attr" INTEGER,
                          		"counter" INTEGER,
								PRIMARY KEY(instance, time, attr))`

	kUpdateInstanceTable = `INSERT INTO %s(instance, time, attr, counter) values("%s", %d, %d, %d)
                            ON CONFLICT(instance, time, attr) DO UPDATE SET counter = counter + %d`
)

type db struct {
	handle *sql.DB
}

func (d *db) open(name string) error {
	var err error
	d.handle, err = sql.Open("sqlite3", name)
	if err != nil {
		log.Printf("%q: %s\n", err, name)
		return err
	}
	_, err = d.handle.Exec(kCreateAttrInfo)
	if err != nil {
		log.Printf("%q: %s\n", err, kCreateAttrInfo)
		return err
	}

	_, err = d.handle.Exec(kCreateViewInfo)
	if err != nil {
		log.Printf("%q: %s\n", err, kCreateViewInfo)
		return err
	}
	return nil
}

func (d *db) updateInstanceAttr(instance string, t int64, attr int, counter int) (err error) {
	var tableName string
	if tableName, err = d.createInstanceTable(); err != nil {
		return err
	}

	sqlStmt := fmt.Sprintf(kUpdateInstanceTable, tableName, instance, t, attr, counter, counter)

	_, err = d.handle.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	return nil
}

func (d *db) createInstanceTable() (name string, err error) {
	name = "attr_report_" + time.Now().Format("2006_01_02")
	intanceTableCreate := fmt.Sprintf(kCreateInstanceTable, name)

	_, err = d.handle.Exec(intanceTableCreate)
	if err != nil {
		log.Printf("%q: %s\n", err, intanceTableCreate)
		return name, err
	}

	return name, nil
}

type attrCounter struct {
	T       int64 `json:"t"`
	Counter int   `json:"counter"`
}

type attrCounterArr struct {
	Counters []attrCounter `json:"counters"`
}

func (d *db) queryAttrSequence(sqlStmt string) (results attrCounterArr) {
	log.Println(sqlStmt)

	rows, err := d.handle.Query(sqlStmt)
	if err != nil {
		return results
	}
	defer rows.Close()

	for rows.Next() {
		r := attrCounter{}
		err = rows.Scan(&r.T, &r.Counter)
		if err != nil {
			continue
		}
		results.Counters = append(results.Counters, r)
	}
	return results
}

func (d *db) queryByInstaceAndAttr(table string, instance string, attr int) (results attrCounterArr) {
	sqlStmt := fmt.Sprintf("SELECT time,counter FROM attr_report_%s where instance=\"%s\" and attr=%d;",
		table, instance, attr)

	return d.queryAttrSequence(sqlStmt)
}

func (d *db) queryByViewAndAttr(table string, view string, attr int)(results attrCounterArr) {
	sqlStmt := fmt.Sprintf(
		`SELECT time,sum(counter) FROM
		(SELECT instance AS view_instance FROM view_info WHERE view="%s") 
		INNER JOIN attr_report_%s ON view_instance = attr_report_%s.instance
		where attr=%d
		GROUP BY time`, view, table, table, attr)

	return d.queryAttrSequence(sqlStmt)
}

type attrName struct {
	Attr int    `json:"attr"`
	Name string `json:"name"`
}

type attrArr struct {
	Attr []attrName `json:"attr_name"`
}

func (d *db) queryAttr(sqlStmt string) (results attrArr) {
	rows, err := d.handle.Query(sqlStmt)
	if err != nil {
		return results
	}
	defer rows.Close()

	for rows.Next() {
		attr := 0
		name := ""
		err = rows.Scan(&attr, &name)
		if err != nil {
			continue
		}
		results.Attr = append(results.Attr, attrName{attr, name})
	}
	return results
}

func (d *db) queryAttrByInstance(table string, instance string) (results attrArr) {
	sqlStmt := fmt.Sprintf(
		`SELECT attr, name FROM 
		(SELECT DISTINCT attr AS attr_report_attr FROM attr_report_%s WHERE instance="%s")
		INNER JOIN attr_info ON attr_report_attr = attr_info.attr`, table, instance)

	return d.queryAttr(sqlStmt)
}

func (d *db) queryAttrByView(table string, view string) (results attrArr) {
	sqlStmt := fmt.Sprintf(
		`SELECT attr, name FROM
		(SELECT DISTINCT attr AS attr_report_attr FROM 
  			(SELECT instance AS view_instance FROM view_info WHERE view="%s")
  			INNER JOIN attr_report_%s ON view_instance = attr_report_%s.instance)
		INNER JOIN attr_info ON attr_report_attr = attr_info.attr`, view, table, table)

	return d.queryAttr(sqlStmt)
}
