package main

import (
	"encoding/xml"
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

//If we directly access loc through sitemap, no need of multiple structs
type SiteMapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}
type News struct {
	Title    []string `xml:url>news>title`
	Keywords []string `xml:url>news>keywords`
	Loc      []string `xml:url>loc`
}

func main() {
	var s SiteMapIndex
	var n News
	url := "https://www.washingtonpost.com/news-sitemaps/index.xml"
	resp, _ := getResp(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	// stringBody := convertByteToString(resp)
	resp.Body.Close()
	xml.Unmarshal(bytes, &s)
	//fmt.Println(s.Locations)
	for _, sites := range s.Locations { //Range goes through index and values
		resp, _ := getResp(sites)
		bytes, _ = ioutil.ReadAll(resp.Body)
		xml.Unmarshal(bytes, &n)
	}
}
