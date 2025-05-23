package main

import "go_logistics/router"

func main() {
	server := router.Router()
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
