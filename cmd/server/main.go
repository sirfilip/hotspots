package main

import (
	"context"

	"necsam/config"
	"necsam/db"
	nescamhttp "necsam/http"
)

func main() {
	client := db.MongoClient()
	defer client.Disconnect(context.TODO())

	server := nescamhttp.New(client)
	server.Run(config.Get("port"))
}
