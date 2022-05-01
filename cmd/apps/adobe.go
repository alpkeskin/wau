package apps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/fatih/color"
)

type adobeResponse []struct {
	HasT2ELinked bool `json:"hasT2ELinked"`
}

func Adobe(wg *sync.WaitGroup, email string, showFalse bool) {
	defer wg.Done()
	var endpoint string = "https://auth.services.adobe.com/signin/v2/users/accounts"

	var jsonStr = []byte(`{"username":"` + email + `"}`)

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr)) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	r.Header.Add("X-Ims-Clientid", "adobedotcom2")

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		body, _ := ioutil.ReadAll(res.Body)
		var response adobeResponse
		json.Unmarshal(body, &response)
		if len(response) > 0 {
			fmt.Println("Adobe \U0001f440")
		} else {
			if showFalse {
				fmt.Println("Adobe", color.RedString(" [Not here!]"))
			}
		}
	} else {
		color.Red("Couldn't check Adobe!")
	}
}
