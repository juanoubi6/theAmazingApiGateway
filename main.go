package main

import (
	"theAmazingApiGateway/app/common"
	"theAmazingApiGateway/app/router"
)

func main() {
	common.ConnectToDatabase()
	router.CreateRouter()
	router.RunRouter()
}
