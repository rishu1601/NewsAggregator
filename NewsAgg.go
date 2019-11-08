package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func getResp(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	return resp, err
}

//Convert received response to string
func convertByteToString(resp *http.Response) string {
	bytes, _ := ioutil.ReadAll(resp.Body)
	stringBody := string(bytes)
	return stringBody
}

/*
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
<sitemap>
<loc>
https://www.washingtonpost.com/news-sitemaps/politics.xml
</loc>
</sitemap>
<sitemap>
<loc>
https://www.washingtonpost.com/news-sitemaps/opinions.xml
</loc>
</sitemap>
<sitemap>
<loc>
https://www.washingtonpost.com/news-sitemaps/local.xml
</loc>
</sitemap>
</sitemapindex>
*/

//If we directly access loc through sitemap, no need of multiple structs
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

func main() {
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
	//Printing information in the Map
	for idx, data := range newsMap {
		fmt.Println("Title:", idx)
		fmt.Println("\nKeywords:", data.Keywords)
		fmt.Println("\nLocation:", data.Location)
		fmt.Println("\n")
	}
	resp.Body.Close()

}
