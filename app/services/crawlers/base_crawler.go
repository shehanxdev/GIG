// https://jdanger.com/build-a-web-crawler-in-go.html
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/collectlinks"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var visited = make(map[string]bool)

func main() {
	flag.Parse()

	args := flag.Args()
	fmt.Println(args)
	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	queue := make(chan string)

	go func() { queue <- args[0] }()

	for uri := range queue {
		f, err := os.OpenFile("tmp/links", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(uri + "\n"); err != nil {
			panic(err)
		}

		enqueue(uri, queue)
	}
}

func enqueue(uri string, queue chan string) {
	fmt.Println("fetching", uri)
	visited[uri] = true
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}
	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err  := ioutil.ReadAll(resp.Body)
	go func() {

		fmt.Println("result: ",string(body))
		fmt.Println("read error is:", err)
	}()

	links := collectlinks.All(resp.Body)
	//f, err := os.OpenFile("tmp/dat1", os.O_APPEND|os.O_WRONLY, 0600)
	//if err != nil {
	//	panic(err)
	//}
	//
	//if _, err = f.WriteString(string(body)+ "\n"); err != nil {
	//	panic(err)
	//}
	//
	//defer f.Close()

	for _, link := range links {
		absolute := fixUrl(link, uri)
		if uri != "" {
			if !visited[absolute] {
				go func() { queue <- absolute }()
			}
		}
	}
}

func fixUrl(href, base string) (string) {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}
