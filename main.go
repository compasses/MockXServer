package main

import (
	"log"

	"github.com/Compasses/MockXServer/middleware"
)

const banner string = `

			Mock Server

`

func main() {
	log.Println(banner)
	log.Printf("Git commit:%s\n", Version)
	log.Printf("Build time:%s\n", Compile)
	middle := middleware.NewMiddleware()
	middle.Run()
}
