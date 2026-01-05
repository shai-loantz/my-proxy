package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HandleHealthCheck(writer http.ResponseWriter, _ *http.Request) {
	log.Println("health check was called")
	fmt.Fprint(writer, "healthy :)")
}
