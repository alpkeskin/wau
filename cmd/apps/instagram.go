package apps

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/fatih/color"
)

func getCSRFToken() string {
	var url string = "https://instagram.com"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	res, _ := client.Do(req)
	if res.StatusCode == 200 {
		r, _ := regexp.Compile("csrftoken.+?;")
		var pattern string
		for _, v := range res.Header {
			for _, v2 := range v {
				if r.FindString(v2) != "" {
					pattern = r.FindString(v2)
					break
				}
			}
		}
		var index int = strings.Index(pattern, "=")
		var token string = pattern[index+1 : len(pattern)-1]
		return token
	}
	return ""
}

func Instagram(wg *sync.WaitGroup, email string, showFalse bool) {
	defer wg.Done()
	var token string = getCSRFToken()
	if token == "" {
		color.Red("Couldn't get CSRF token for Instagram!")
	} else {
		var endpoint string = "https://www.instagram.com/accounts/web_create_ajax/attempt/"

		data := url.Values{}
		data.Set("email", email)

		client := &http.Client{}
		r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
		if err != nil {
			log.Fatal(err)
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
		r.Header.Add("Cookie", "csrftoken="+token+";")
		r.Header.Add("X-Csrftoken", token)

		res, err := client.Do(r)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode == 200 {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			match, _ := regexp.MatchString("email_is_taken", string(body))
			if match {
				fmt.Println("Instagram \U0001f440")
			} else {
				if showFalse {
					fmt.Println("Instagram", color.RedString(" [Not here!]"))
				}
			}
		} else {
			color.Red("Couldn't check Instagram!")
		}
	}
}
