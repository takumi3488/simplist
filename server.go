package main

import (
	"database/sql"
	"log"
	"os"
	"simplist/handler"
	"strings"

	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"
)

func main() {
	// DBに接続する
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	h := &handler.Handler{DB: db}

	// テーブルを作成する
	createTable := `CREATE TABLE IF NOT EXISTS lists (
		key TEXT PRIMARY KEY, items TEXT[] DEFAULT '{}' NOT NULL
	)`
	if _, err := db.Exec(createTable); err != nil {
		panic(err)
	}

	// データを挿入する
	for _, k := range strings.Split(os.Getenv("KEYS"), ",") {
		if k == "" {
			continue
		}
		log.Printf("Inserting key: %s", k)
		if _, err := db.Exec(`INSERT INTO lists (key) VALUES ($1) ON CONFLICT (key) DO NOTHING`, k); err != nil {
			panic(err)
		}
	}

	// Echoを起動する
	e := echo.New()
	e.GET("/", h.GetLists)
	e.GET("/:key", h.GetList)
	e.PUT("/:key", h.UpdateList)
	e.Logger.Fatal(e.Start(":8080"))
}
