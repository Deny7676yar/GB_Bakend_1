package main

import (
	"homework4/internal/api"
	"log"
	"net/http"
	"time"
)

func main() {
	handler := &api.Handler{}
	http.Handle("/", handler)

	uploadHandler := &api.UploadHandler{
		UploadDir: "upload",
	}
	http.Handle("/upload", uploadHandler)

	fileListHandler := &api.FileListHandler{
		FileList: "upload",
	}
	http.Handle("/upload", fileListHandler)

	dirToSave := http.Dir(uploadHandler.UploadDir)
	fs := &http.Server{
		Addr:         ":8080",
		Handler:      http.FileServer(dirToSave),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(fs.ListenAndServe())

	srv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

	listFile := http.Dir(fileListHandler.FileList)
	fl := &http.Server{
		Addr:         ":8080",
		Handler:      http.FileServer(listFile),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(fl.ListenAndServe())
}
