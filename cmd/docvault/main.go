package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"engineering-document-vault/internal/service"
)

func main() {
	path := flag.String("db", "docvault.db", "database path")
	flag.Parse()
	app, err := service.OpenApplication(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()
	if _, err := os.Stdout.WriteString(fmt.Sprintf("工程项目文档盘已启动，数据库：%s\n", *path)); err != nil {
		log.Fatal(err)
	}
}
