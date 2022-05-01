package apps

import (
	"encoding/json"
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

type vscoResponse struct {
	Status int `json:"status"`
}

func getVSID() string {
	var url string = "https://vsco.co"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	res, _ := client.Do(req)
	if res.StatusCode == 200 {
		r, _ := regexp.Compile("vs_anonymous_id=.+;")
		var pattern string
		for _, v := range res.Header {
			for _, v2 := range v {
				if r.FindString(v2) != "" {
					pattern = r.FindString(v2)
					break
				}
			}
		}
		var i int = strings.Index(pattern, ";")
		return pattern[16:i]
	}
	return ""
}

func Vsco(wg *sync.WaitGroup, email string, showFalse bool) {
	defer wg.Done()
	var token string = getVSID()
	if token == "" {
		color.Red("Couldn't get VS ID from VSCO!")
	} else {
		var endpoint string = "https://vsco.co/ajx/user/doForgotPassword"

		data := url.Values{}
		data.Set("email", email)

		client := &http.Client{}
		r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
		if err != nil {
			log.Fatal(err)
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
		r.Header.Add("Cookie", "vs_anonymous_id="+token+";")

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
			var response vscoResponse
			json.Unmarshal(body, &response)
			if response.Status == 1 {
				fmt.Println("Vsco \U0001f440")
			} else {
				if showFalse {
					fmt.Println("Vsco", color.RedString(" [Not here!]"))
				}
			}
		} else {
			color.Red("Couldn't check VSCO!")
		}
	}
}
