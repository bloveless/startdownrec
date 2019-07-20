package startdownrec

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Run is used by google cloud to kick off the function
func Run(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Hostname string `json:"hostname"`
		Status   string `json:"status"`
	}

	decErr := json.NewDecoder(r.Body).Decode(&d)
	if decErr != nil || d.Status == "" || d.Hostname == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, wErr := w.Write([]byte("Invalid JSON body"))
		if wErr != nil {
			log.Fatal("Error responding with error")
		}
		return
	}

	fmt.Println("Hostname: " + d.Hostname)
	fmt.Println("Status: " + d.Status)

	w.WriteHeader(200)
	_, wErr := w.Write([]byte("Status recorded"))
	if wErr != nil {
		log.Fatal("Error responding with success")
	}
}
