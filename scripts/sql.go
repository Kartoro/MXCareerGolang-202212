package main

import (
	"MXCareerGolang-202212/config"
	"MXCareerGolang-202212/internal/util"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"sync"
)

func main() {
	// generate test data
	db, _ := sql.Open("mysql", config.DSN)
	var wg sync.WaitGroup
	ch := make(chan int, 1000)
	for i := 1001; i <= 10000000; i++ {
		ch <- 1
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, e := db.Exec("INSERT INTO user (username, nickname, password, avatar) VALUES (?,?,?,?)",
				i, i, util.SHA256(strconv.Itoa(i)), config.DefaultImagePath)
			if e != nil {
				log.Println(e)
			}
			<-ch
		}(i)
		wg.Wait()
	}
}
