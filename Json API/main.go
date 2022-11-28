package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/user", makeHttpHandler(handleGetUserByID))
	fmt.Println("Server started at 2000")
	http.ListenAndServe(":2000", nil)
}
