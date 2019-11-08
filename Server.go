package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"text/template"
)

var wg sync.WaitGroup

type SiteMapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}
type News struct {
	Title    []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Loc      []string `xml:"url>loc"`
}

type NewsMap struct {
	Keywords string
	Location string
}

type NewsAggPage struct {
	Title string
	News  map[string]NewsMap
}

func newsRoutine(news chan News, sites string) {
	defer wg.Done() //Synchronisation
	var n News
	sites = strings.TrimSpace(sites)
	fmt.Println(sites)
	resp, err := http.Get(sites)
	if err != nil {
		log.Fatal(err)
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &n)
	resp.Body.Close()

	news <- n
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hi,this is go, u are here</h1>")
}

func newsAggregatorHandler(w http.ResponseWriter, r *http.Request) {
	var s SiteMapIndex
	newsMap := make(map[string]NewsMap)
	url := "https://www.washingtonpost.com/news-sitemaps/index.xml"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
	resp.Body.Close()
	queue := make(chan News, 30)

	for _, sites := range s.Locations {
		wg.Add(1) //Adding to waitgroup
		go newsRoutine(queue, sites)
	}

	wg.Wait()
	close(queue)

	for elem := range queue {
		for idx := range elem.Keywords {
			newsMap[elem.Title[idx]] = NewsMap{elem.Keywords[idx], elem.Loc[idx]}
		}
	}
	p := NewsAggPage{Title: "Basic News Aggregator", News: newsMap}
	t, _ := template.ParseFiles("newsAggTemplate.html")
	t.Execute(w, p)
}
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg", newsAggregatorHandler)
	http.ListenAndServe(":8000", nil)
}
