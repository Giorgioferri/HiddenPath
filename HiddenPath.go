package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var site string

func test(site string) {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("error")
		return
	}
	defer resp.Body.Close()
	fmt.Println("status code:", resp.StatusCode)
	if resp.StatusCode == 200 || resp.StatusCode != 404 {
		fmt.Printf("found! [%d] %s\n", resp.StatusCode, site)
	} else {
		fmt.Printf("not found [%d] %s\n", resp.StatusCode, site)
	}
}

func scan(site string, path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("hi")
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		g := scanner.Text()
		url := fmt.Sprintf("%s/%s", site, g)
		test(url)

	}

}

func main() {
	sito := flag.String("site", "url", "use this for set a site")
	wordlist := flag.String("wordlist", "wordlist.txt", "use this for set a site")
	flag.Parse()

	tests := fmt.Sprintf("%s", *sito)
	scan(tests, *wordlist)

}
