package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommands_SuccessfulExecutionOfTheCommands(t *testing.T) {
	request := bytes.NewBufferString(`{"command": ["cmd.exe","/c","echo","hello"], "timeout": 1}`)

	req := httptest.NewRequest(http.MethodPost, "/command", request)
	rec := httptest.NewRecorder()

	CommandHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Ожидался код %d, но получен %d", http.StatusOK, rec.Code)
	}

	expectedBody := "The command was executed successfully\n"
	if rec.Body.String() != expectedBody {

		fmt.Println(rec.Body.String())
		t.Errorf("Ожидалось тело '%s', но получено '%s'\n", expectedBody, rec.Body.String())
	}
}

func TestRequest_InvalidRequestFormatHandling(t *testing.T) {
	request := bytes.NewBufferString(`{"command": ["cmd.exe","/c","echo","hello"], "timeout": "WrongType"}`)

	req := httptest.NewRequest(http.MethodPost, "/command", request)
	// Создаём тестовый HTTP-ответ
	rec := httptest.NewRecorder()
	// Вызываем хендлер
	CommandHandler(rec, req)
	// Проверяем код ответа
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код %d, но получен %d", http.StatusBadRequest, rec.Code)
	}
	// Проверяем тело ответа
	expectedBody := "Invalid request body\n"
	if rec.Body.String() != expectedBody {

		fmt.Println(rec.Body.String())
		t.Errorf("Ожидалось тело '%s', но получено '%s'\n", expectedBody, rec.Body.String())
	}
}

func TestTimeout_EndOfTime(t *testing.T) {
	///Проверяющий запуск команды в ситуации, когда она не заканчивается в отведенное время
	request := bytes.NewBufferString(`{"command": ["cmd.exe","/c","echo","hello"], "timeout": 0.0001}`)

	req := httptest.NewRequest(http.MethodPost, "/command", request)
	// Создаём тестовый HTTP-ответ
	rec := httptest.NewRecorder()
	// Вызываем хендлер
	CommandHandler(rec, req)
	// Проверяем код ответа
	if rec.Code != http.StatusRequestTimeout {
		t.Errorf("Ожидался код %d, но получен %d", http.StatusRequestTimeout, rec.Code)
	}
	// Проверяем тело ответа
	expectedBody := "Request timeout\n"
	if rec.Body.String() != expectedBody {

		fmt.Println(rec.Body.String())
		t.Errorf("Ожидалось тело '%s', но получено '%s'\n", expectedBody, rec.Body.String())
	}
}
