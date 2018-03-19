package main

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
)

func enforceXMLHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		if http.DetectContentType(buf.Bytes()) != "text/xml; charset=utf-8" {
			http.Error(w, http.StatusText(415), 415)
			w.Write([]byte("Content type != to text/xml; charset=utf-8\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
func myLoggingHandler(h http.Handler) http.Handler {
	logFile, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return handlers.LoggingHandler(logFile, h)

}
func final(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing final")
	w.Write([]byte("OK"))
}
func main() {
	finalHandler := http.HandlerFunc(final)
	authHandler := httpauth.SimpleBasicAuth("username", "password")

	http.Handle("/", myLoggingHandler(authHandler(enforceXMLHandler(finalHandler))))
	http.ListenAndServe(":3000", nil)
}
