package main

import (
	"github.com/antik10ud/mymux"
	"net/http"
)

func main() {
	handler := mymux.NewRouterTemplateHandler(myErrorHandler)
	handler.RegisterType("caps", "[A-Z]{1,64}")
	handler.AppendRoute("GET", "/echo/Text:caps", mymux.Adapt(caps))
	server := &http.Server{
		Addr:    "0.0.0.0:23127",
		Handler: handler,
	}
	server.ListenAndServe()

}

func caps(w http.ResponseWriter, r *http.Request) {
	vars := mymux.GetVars(r)
	text := vars["Text"]
	w.Write([]byte(text))
}

func myErrorHandler(w http.ResponseWriter, status int, detail string) {
	w.WriteHeader(status)
	w.Write([]byte(detail))
}
