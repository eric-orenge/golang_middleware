package main

import (
	"bytes"
	"log"
	"net/http"
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
			w.Write([]byte("Not content type not equal to text/xml"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
func final(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing final")
	w.Write([]byte("OK"))
}
func main() {
	finalHandler := http.HandlerFunc(final)

	http.Handle("/", enforceXMLHandler(finalHandler))
	http.ListenAndServe(":3000", nil)
}
