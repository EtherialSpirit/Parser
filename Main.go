package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "golang.org/x/net/html"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type siteRoad struct{
	status bool
	url string
}

func main() {

	readyFile := readingFile()
	continueProcessingFile(readyFile)
	fmt.Println("Program end")
}

func readingFile() []string{

	file, err := os.Open("Site.txt")
	_check(err)
	defer file.Close()

	// получить размер файла
	stat, err := file.Stat()
	_check(err)
	// чтение файла
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	_check(err)

	str := string(bs)

	readFile := strings.Split(str, "\r\n")
	return readFile
}

func continueProcessingFile(readyFile []string){

	var clearLink string
	//var wg sync.WaitGroup

	//(a href ?=?")+.{2,}(contact)
	for i :=0; i<len(readyFile);i++ {

		//wg.Add(1)
		//defer wg.Done()

		checkURLFalseLink, err := regexp.MatchString(`wikipedia|google.com|www.paypal.com|vikidia`, strings.ToLower(readyFile[i]))

		if checkURLFalseLink == true {
			clearLink = "This URL doesn`t contain the correct contact."
			fmt.Println(siteRoad{status: false, url: clearLink})
			continue
		}

		checkHTTP, err := regexp.MatchString(`^https?.*`, readyFile[i])
		_check(err)

		if checkHTTP == true {
			receivedURL, err := http.Get(readyFile[i])

			if err != nil {
				fmt.Println(err)
				continue
			}

			if receivedURL.StatusCode != 200 {
				fmt.Println("Get ", readyFile[i], "status code error: ", receivedURL.StatusCode, receivedURL.Status)
				//defer res.Body.Close()
			} else {

				if checkURL(readyFile[i]) == true {
					clearLink = readyFile[i]
					//linkForPrint := siteRoad{status: true, url: clearLink}
					fmt.Println(siteRoad{status: true, url: clearLink})
					//fmt.Println(readyFile[i])
				} else {
					doc, err := goquery.NewDocumentFromReader(receivedURL.Body)
					_check(err)
					contactURL := linkScrape(doc)

					if contactURL !=""{
						clearLink = condition(contactURL, readyFile[i])
						fmt.Println(siteRoad{status: true, url: clearLink})
					}else {
						clearLink = "This URL doesn`t contain the contact."
						fmt.Println(siteRoad{status: false, url: clearLink})
					}
				}
			}
			defer receivedURL.Body.Close()
		}
	}
//wg.Wait()

}

func linkScrape(domHtml *goquery.Document)  string{

	var link string
	domHtml.Find("a[href]").Each(func(i int, s *goquery.Selection) {

		var href, _ = s.Attr("href")

		if checkURL(s.Text()) == true{
			link = href
			return
		} else {

			if checkURL(href) == true {
				link = href
				return
			}
		}

	})

	return string(link)
}

func checkURL(URL string) bool {

	checkURL, err := regexp.MatchString(`[^в]контакты?|\bcontacts?|joindre|\bkontakte?|contactos?|contacta|contacter|kontakty|info\.html`, strings.ToLower(URL))
	_check(err)
	if checkURL == true {
		checkURLFalse, err := regexp.MatchString(`boissiere|контактная|контактный|линз|пар|педал|linsen|len|grill|börse|reiniger|câble|
			rad|thermom|blut|change|pay|board|zahlung|cuota|pago|pai`, strings.ToLower(URL))
		_check(err)
		if checkURLFalse == true {
			return false
		}
	}
	return checkURL
}

func condition(link string, url1 string) string{

	reg, err := regexp.MatchString(`https?|/?www\.`, strings.ToLower(link))
	_check(err)

	if reg ==true{
		return link
	} else {
		re := regexp.MustCompile(".*://|/.*")
		cleanLink := re.ReplaceAllString(url1, "")
		link = "http://"+cleanLink + "/"+link
		return  link
	}

}

func _check(err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
}

