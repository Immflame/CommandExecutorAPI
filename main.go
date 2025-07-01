package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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

		if timeout_time > resp.Timeout {
			timeout_time = resp.Timeout
		}

		if timeout_time <= 0 {
			timeout_time = 1
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout_time)*time.Second)
		defer cancel()

		cmdStruct := exec.CommandContext(ctx, resp.Command[0], resp.Command[1:]...)
		out, err := cmdStruct.Output()

		if ctx.Err() != nil {
			http.Error(w, "Error:"+string(ctx.Err().Error()), http.StatusRequestTimeout)
			return
		}

		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
			return
		}

		fmt.Println(string(out))

		if err != nil {
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("Command execution error: %s", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "The command was executed successfully: %s", string(out))

	default:
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("2 args")
		return
	}

	address := string(os.Args[1])
	s, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Second arg must be float64")
		return
	}
	timeout_time = float64(s)

	http.HandleFunc("/command", CommandHandler)

	fmt.Printf("Server listening localhost:%s \nMax execution time (sec): %f\n", address, timeout_time)
	if err := http.ListenAndServe(":"+address, nil); err != nil {
		fmt.Println(err)
	}
}
