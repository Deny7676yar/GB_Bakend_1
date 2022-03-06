package api_test

import (
	"bytes"
	"fmt"
	"homework4/internal/api"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetHandler(t *testing.T) {

	// Создаем запрос с указанием нашего хендлера. Так как мы тестируем GET-эндпоинт
	// то нам не нужно передавать тело, поэтому третьим аргументом передаем nil
	req, err := http.NewRequest("GET", "/?name=John", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
	// и используем его для получения ответа
	rr := httptest.NewRecorder()
	handler := &api.Handler{}

	// Наш хендлер соответствует интерфейсу http.Handler, а значит
	// мы можем использовать ServeHTTP и напрямую указать
	// Request и ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := `Parsed query-param with key "name": John`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestUploadHandler(t *testing.T) {
	// открываем файл, который хотим отправить
	uploadDir := "/home/den/GB_Bakend_1/homework4/internal/api/upload"
	file, _ := os.Open(uploadDir + "/testfile.txt")
	defer file.Close()

	// действия, необходимые для того, чтобы засунуть файл в запрос
	// в качестве мультипарт-формы
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	//опять создаем запрос, теперь уже на /upload эндпоинт
	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// создаем ResponseRecorder
	rr := httptest.NewRecorder()

	// создаем заглушку файлового сервера. Для прохождения тестов
	// нам достаточно чтобы он возвращал 200 статус
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok!")
	}))
	defer ts.Close()

	uploadHandler := &api.UploadHandler{
		// таким образом мы подменим адрес файлового сервера
		// и вместо реального, хэндлер будет стучаться на заглушку
		// которая всегда будет возвращать 200 статус, что нам и нужн
		UploadDir: "upload",
		HostAddr:  ts.URL,
	}

	// опять же, вызываем ServeHTTP у тестируемого обработчика
	uploadHandler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)

	}

	expected := `testfile`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFileListHandler(t *testing.T) {
	uploadDir := "/home/den/GB_Bakend_1/homework4/internal/api/upload"
	req, err := http.NewRequest("GET", "/list&ext=.jpg", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := &api.FileListHandler{
		FileList: uploadDir,
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "Name:jpegName	Ext:.jpg	Size:0"

	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
