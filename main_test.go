package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCommands_SuccessfulExecutionOfTheCommands(t *testing.T) {
	request := bytes.NewBufferString(`{"command": ["cmd.exe","/c","echo","hello"], "timeout": 3}`)

	req := httptest.NewRequest(http.MethodPost, "/command", request)
	rec := httptest.NewRecorder()

	CommandHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Ожидался код %d, но получен %d", http.StatusOK, rec.Code)
	}

	actualBody := strings.TrimSpace(strings.ReplaceAll(rec.Body.String(), "\r\n", "\n"))
	expectedBody := strings.TrimSpace(strings.ReplaceAll("The command was executed successfully: hello", "\r\n", "\n"))
	if !strings.Contains(actualBody, expectedBody) {
		t.Errorf("Expected '%s', but got '%s'\n", expectedBody, actualBody)
	}
}

func TestRequest_InvalidRequestFormatHandling(t *testing.T) {
	request := bytes.NewBufferString(`{"command": ["cmd.exe","/c","echo","hello"], "timeout": "WrongType"}`)

	req := httptest.NewRequest(http.MethodPost, "/command", request)
	rec := httptest.NewRecorder()
	CommandHandler(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код %d, но получен %d", http.StatusBadRequest, rec.Code)
	}
	expectedBody := "Invalid request body\n"
	if rec.Body.String() != expectedBody {

		fmt.Println(rec.Body.String())
		t.Errorf("Ожидалось тело '%s', но получено '%s'\n", expectedBody, rec.Body.String())
	}
}

func TestTimeout_EndOfTime(t *testing.T) {
	request := bytes.NewBufferString(`{"command": ["cmd.exe","/c","echo","hello"], "timeout": 0.0001}`)

	req := httptest.NewRequest(http.MethodPost, "/command", request)
	rec := httptest.NewRecorder()
	CommandHandler(rec, req)
	if rec.Code != http.StatusRequestTimeout {
		t.Errorf("Ожидался код %d, но получен %d", http.StatusRequestTimeout, rec.Code)
	}
	expectedBody := "Error:context deadline exceeded\n"
	if rec.Body.String() != expectedBody {

		fmt.Println(rec.Body.String())
		t.Errorf("Ожидалось тело '%s', но получено '%s'\n", expectedBody, rec.Body.String())
	}
}
