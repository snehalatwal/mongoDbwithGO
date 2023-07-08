package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/snehaMongoDb/router"
)

func main() {
	fmt.Println("MongoDB API")
	r := router.Router()
	log.Fatal(http.ListenAndServe(":4000", r))
}
