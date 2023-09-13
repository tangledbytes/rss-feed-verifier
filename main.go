package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

// EXCEPTIONS is a list of URLs that are known to be valid RSS feeds but do not
// have the correct Content-Type header.
//
// NOTE: Phil Opperman's blog is a valid RSS feed but YARR cannot parse it as well.
var EXCEPTIONS = []string{
	"https://os.phil-opp.com/atom.xml",
	"https://www.digitalocean.com/blog/rss",
}

// OpmlFormat is the format of an OPML file.
type OpmlFormat struct {
	XMLName xml.Name `xml:"opml"`
	Head    struct {
		Title string `xml:"title"`
	} `xml:"head"`
	Body struct {
		Outlines []OpmlOutline `xml:"outline"`
	} `xml:"body"`
}

// OpmlOutline is the format of an outline (limited to my needs) in an OPML file.
type OpmlOutline struct {
	Outlines []OpmlOutline `xml:"outline"`
	Text     string        `xml:"text,attr"`
	Type     string        `xml:"type,attr"`
	XMLUrl   string        `xml:"xmlUrl,attr"`
	HtmlUrl  string        `xml:"htmlUrl,attr"`
}

// VerifyFeed takes a list of outlines and verifies that the XMLUrl is a valid
// RSS feed.
// VerifyFeed is recursive.
func VerifyFeed(outlines []OpmlOutline) {
	var wg sync.WaitGroup

	for _, outline := range outlines {
		wg.Add(1)
		go func(outline OpmlOutline) {
			if outline.Type == "rss" {
				if !ValidateXMLUrl(outline.XMLUrl) {
					fmt.Println("Invalid:", outline.XMLUrl)
				}
			}
			if outline.Outlines != nil {
				VerifyFeed(outline.Outlines)
			}

			wg.Done()
		}(outline)
	}

	wg.Wait()
}

// ValidateXMLUrl takes a URL and verifies that it is a valid RSS feed.
func ValidateXMLUrl(url string) bool {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	contentType := resp.Header.Get("Content-Type")
	for _, exception := range EXCEPTIONS {
		if strings.Contains(url, exception) {
			fmt.Println("Exception:", url)
			return true
		}
	}

	if !strings.Contains(contentType, "xml") {
		return false
	}

	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: rss-feed-verifier <opml-file>")
		os.Exit(1)
	}

	opmldata, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	data := &OpmlFormat{}
	err = xml.Unmarshal(opmldata, data)
	if nil != err {
		fmt.Println("Error unmarshalling OPML:", err)
		os.Exit(1)
	}

	VerifyFeed(data.Body.Outlines)
}
