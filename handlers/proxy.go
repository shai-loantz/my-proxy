package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HandleProxy(writer http.ResponseWriter, request *http.Request) {
	log.Println("proxy handler was called")
	fmt.Fprint(writer, "Implement me")
}
