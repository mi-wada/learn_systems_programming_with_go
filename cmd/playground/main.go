package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json")
		src := map[string]string{
			"hello": "world",
		}
		srcJSON, _ := json.Marshal(src)

		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		multiWriter := io.MultiWriter(log.Writer(), gzipWriter)
		multiWriter.Write(srcJSON)
	})

	http.ListenAndServe(":8080", nil)
}
