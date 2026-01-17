package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	uploadDir = "./public/img"
	maxSize   = 10 * 1024 * 1024 // 10MB
)

type UploadResponse struct {
	URL string `json:"url"`
}

func main() {
	// Создаем директорию для загрузок, если она не существует
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, 0755)
		if err != nil {
			log.Fatal("Не удалось создать директорию для загрузок:", err)
		}
	}

	// Раздача статических файлов из папки public
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Обработчик загрузки изображений
	http.HandleFunc("/upload", uploadHandler)

	port := ":8080"
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Ограничение размера файла
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	if err := r.ParseMultipartForm(maxSize); err != nil {
		http.Error(w, "Файл слишком большой", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Ошибка получения файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Генерация уникального имени файла
	ext := filepath.Ext(handler.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, newFileName)

	// Создание файла на диске
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Копирование содержимого
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Ошибка записи файла", http.StatusInternalServerError)
		return
	}

	// Ответ с URL файла
	url := fmt.Sprintf("/img/%s", newFileName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadResponse{URL: url})
}
