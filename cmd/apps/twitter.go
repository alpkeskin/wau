package apps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/fatih/color"
)

type twitterResponse struct {
	Valid bool   `json:"valid"`
	Msg   string `json:"msg"`
	Taken bool   `json:"taken"`
}

func Twitter(wg *sync.WaitGroup, email string, showFalse bool) {
	defer wg.Done()
	var endpoint string = "https://api.twitter.com/i/users/email_available.json"

	data := url.Values{}
	data.Set("email", email)

	r, err := http.Get(endpoint + "?" + data.Encode())
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode == 200 {
		body, _ := ioutil.ReadAll(r.Body)
		var response twitterResponse
		json.Unmarshal(body, &response)
		if response.Taken {
			fmt.Println("Twitter \U0001f440")
		} else {
			if showFalse {
				fmt.Println("Twitter", color.RedString(" [Not here!]"))
			}
		}
	} else {
		color.Red("Couldn't check Twitter!")
	}
}
