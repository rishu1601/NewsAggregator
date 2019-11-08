package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hi,this is go, u are here</h1>")
}
func newsAggregatorHandler(w http.ResponseWriter, r *http.Request) {
	var s SiteMapIndex
	var n News
	newsMap := make(map[string]NewsMap)
	url := "https://www.washingtonpost.com/news-sitemaps/index.xml"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
	for _, sites := range s.Locations { //Range goes through index and values
		sites = strings.TrimSpace(sites)
		resp, err := http.Get(sites)
		fmt.Println(sites)
		if err != nil {
			log.Fatal(err)
		}
		bytes, _ = ioutil.ReadAll(resp.Body)
		xml.Unmarshal(bytes, &n)

		//Go to each site and extract the data with title as key and loc,keywords as value
		for idx := range n.Keywords {
			newsMap[n.Title[idx]] = NewsMap{n.Keywords[idx], n.Loc[idx]}
		}

	}
	p := NewsAggPage{Title: "Hi, its blown up", News: newsMap} //Just a sample
	t, _ := template.ParseFiles("newsAggTemplate.html")        //Making a template
	t.Execute(w, p)                                            //p is being passed to the template
}
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg", newsAggregatorHandler)
	http.ListenAndServe(":8000", nil)
}
