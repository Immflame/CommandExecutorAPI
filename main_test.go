package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// {
//     "command": ["cmd.exe","/c","echo","hello"],
//     "timeout": 1
// } вывод:

// {
//     "command": ["cmd.exe","/c","echo","hello"],
//     "timeout": 0.00001
// } вывод:

type Resp struct {
	Command []string `json:"command"`
	Timeout float64  `json:"timeout"`
}

func TestCommands_SuccessfulExecutionOfTheCommands(t *testing.T) {
	///Успешное выполнение команд

}

func TestRequest_InvalidRequestFormatHandling(t *testing.T) {
	/// Обработка неверно сформированного запроса
	var buf Resp
	req, err := http.NewRequest("Post", "/command", &buf)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CommandHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Invalid request body"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func TestTimeout_EndOfTime(t *testing.T) {
	///Проверяющий запуск команды в ситуации, когда она не заканчивается в отведенное время
}
