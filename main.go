package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/namtx/glosbe/color"
	"github.com/olekukonko/tablewriter"
)

const USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"

func main() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://glosbe.com/vi/en/"+os.Args[1], nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", USER_AGENT)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Failed to get data from https://glosbe.com, error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"VI", "EN"})
	table.SetRowLine(true)
	table.SetColWidth(55)

	doc.Find("#tmTable > div.tableRow").Each(func(i int, s *goquery.Selection) {
		var source string
		s.Find("div:nth-child(1) > span > span > span").Each(func(i int, ss *goquery.Selection) {
			if ss.HasClass("tm-p-em") {
				source += color.Yellow(ss.Text())
			} else {
				source += color.White(ss.Text())
			}
		})
		var des string
		s.Find("div:nth-child(2) > span > span > span").Each(func(i int, ss *goquery.Selection) {
			if ss.HasClass("tm-p-em") {
				des += color.Red(ss.Text())
			} else {
				des += color.White(ss.Text())
			}
		})

		table.Append([]string{source, des})
	})

	table.Render()
}
