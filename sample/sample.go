package main

import (
	"github.com/antik10ud/mymux"
	"net/http"
	"fmt"
)

var echoRoute mymux.Route

func main() {
	handler := mymux.NewRouterTemplateHandler()
	handler.ErrorHandler(myErrorHandler)
	handler.RegisterType("caps", "[A-Z]{1,64}")
	echoRoute = handler.AppendRoute("GET", "/echo/Text:caps", mymux.Adapt(caps))
	server := &http.Server{
		Addr:    "0.0.0.0:23127",
		Handler: handler,
	}
	server.ListenAndServe()

}

func caps(w http.ResponseWriter, r *http.Request) {
	vars := mymux.GetVars(r)
	text := vars["Text"]
	url := echoRoute.URL(mymux.URLVars{"Text": "OTHER"})

	w.Write([]byte(fmt.Sprintf("You say %s, you can try %s ", text, url)))
}

func myErrorHandler(w http.ResponseWriter, status int, detail string) {
	w.WriteHeader(status)
	w.Write([]byte(detail))
}
