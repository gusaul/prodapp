package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Data struct {
	ProductID int64  `db:"product_id"`
	Alias     string `db:"product_key"`
}

type DataCount struct {
	Alias  string `db:"product_key"`
	Counts int    `db:"counts"`
}

func main() {
	infoFile, err := os.OpenFile("result.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetFlags(0)
	log.SetOutput(infoFile)

	conn, err := sqlx.Connect("postgres", "host=172.21.252.170 user=ma160401 password=yzKagrkIbj7OEf dbname=tokopedia-product sslmode=disable")
	if err != nil {
		panic(err)
	}

	startProduct := int64(0)
	go func() {
		for {
			var result []Data
			err := conn.Select(&result, "select product_id, product_key from ws_product_alias where product_id > $1 order by product_id limit 100", startProduct)
			if err != nil {
				fmt.Println("ERROR", err)
				break
			}

			aliases := make([]string, len(result))
			for i, d := range result {
				aliases[i] = d.Alias
			}

			var resCount []DataCount
			q, args, err := sqlx.In("select product_key, count(product_id) as counts from ws_product_alias where product_key in (?) group by product_key", aliases)
			if err != nil {
				fmt.Println("ERROR", err)
				break
			}
			err = conn.Select(&resCount, conn.Rebind(q), args...)
			if err != nil {
				fmt.Println("ERROR", err)
				break
			}

			for _, rc := range resCount {
				if rc.Counts > 1 {
					log.Println(rc.Alias, ":", rc.Counts)
				}
			}

			startProduct = result[len(result)-1].ProductID
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("STOPPED AT PRODUCT", startProduct)
}
