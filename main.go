package main

import (
	"fmt"
	"myserver/api"
	"net/http"
)

func main() {
	fmt.Println("In main")
	srv := api.NewServer()
	http.ListenAndServe(":8080", srv)
}
