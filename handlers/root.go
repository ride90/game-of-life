package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func (h handler) Root(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	fmt.Fprint(w, "This is root")
}
