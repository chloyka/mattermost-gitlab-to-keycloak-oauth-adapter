package main

import "net/http"

type responseWriterWithStatusCode struct {
	http.ResponseWriter
	statusCode int
	debugMode  bool
	body       []byte
}

func (w *responseWriterWithStatusCode) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWithStatusCode) Write(b []byte) (int, error) {
	if w.debugMode {
		w.body = b
	}

	return w.ResponseWriter.Write(b)
}

func newResponseWriterWithStatusCode(w http.ResponseWriter, debugMode bool) http.ResponseWriter {
	return &responseWriterWithStatusCode{ResponseWriter: w, debugMode: debugMode}
}
