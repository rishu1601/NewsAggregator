package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
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

//SiteMapIndex contains a slice of all the locations,i.e. sitemap tags
type SiteMapIndex struct {
	Locations []Location `xml:"sitemap"`
}

//Location contains data at each loc tag
type Location struct {
	Loc string `xml:"loc"`
}

//This function returns the loc as a string type
func (l Location) String() string {
	return fmt.Sprintf(l.Loc)
}

func main() {
	url := "https://www.washingtonpost.com/news-sitemaps/index.xml"
	resp, _ := getResp(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	// stringBody := convertByteToString(resp)
	resp.Body.Close()
	var s SiteMapIndex
	xml.Unmarshal(bytes, &s)
	//fmt.Println(s.Locations)
	for _, sites := range s.Locations { //Range goes through index and values
		fmt.Println(sites)
	}
}
