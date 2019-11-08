package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type NewsAggPage struct {
	Title string
	News  string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hi,this is go, u are here</h1>")
}
func newsAggregatorHandler(w http.ResponseWriter, r *http.Request) {
	p := NewsAggPage{Title: "Hi, its blown up", News: "Everything is blown up"} //Just a sample
	t, _ := template.ParseFiles("basichtmltemplate.html")                       //Making a template
	t.Execute(w, p)                                                             //p is being passed to the template
}
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg", newsAggregatorHandler)
	http.ListenAndServe(":8000", nil)
}
