package main

import (
	"fmt"
	"net/http"
	handler "servicediscovery/handlers"
)

func main() {
	handler.RegisterHandlers()
	fmt.Println("Service Started at port 8080")
	http.ListenAndServe(":8080", nil)
}
