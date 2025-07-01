package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

var timeout_time float64

type Resp struct {
	Command []string `json:"command"`
	Timeout float64  `json:"timeout"`
}

func CommandHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case "POST":
		var resp Resp

		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		start := time.Now()
		cmdStruct := exec.Command(resp.Command[0], resp.Command[1:]...)
		out, err := cmdStruct.Output()
		duration := time.Since(start)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(out))

		if duration.Seconds() > timeout_time || duration.Seconds() > resp.Timeout {
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "The command was executed successfully")

	default:
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	var address string

	fmt.Scan(&address, &timeout_time)

	http.HandleFunc("/command", CommandHandler)

	fmt.Printf("Server listening localhost:%s \nMax execution time (sec): %f\n", address, timeout_time)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}

}
