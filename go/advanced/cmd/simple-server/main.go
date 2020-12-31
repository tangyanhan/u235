package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		f, fh, err := r.FormFile("content")
		if err != nil {
			log.Println("failed to get form file content:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("Received file, filename=", fh.Filename)
		filePath := path.Join(os.TempDir(), fh.Filename)
		log.Println("Save file to:", filePath)
		dst, err := os.Create(filePath)
		if err != nil {
			log.Println("failed to create dst file at", filePath, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, f); err != nil {
			log.Println("failed to copy file to file", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("Success"))
	})

	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		log.Fatalln(err)
	}
}
