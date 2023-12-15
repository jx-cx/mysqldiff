// main.go

package main

import (
	"mysqldiff/routers"
	_ "mysqldiff/routers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := routers.InitRouter()
	router.LoadHTMLGlob("templates/*")
	router.Run(":8080")
}
