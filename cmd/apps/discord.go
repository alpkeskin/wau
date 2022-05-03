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

type discordResponse struct {
	Errors struct {
		Email struct {
			Errors []struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"_errors"`
		} `json:"email"`
	} `json:"errors"`
}

func Discord(wg *sync.WaitGroup, email string, showFalse bool) {
	defer wg.Done()
	var endpoint string = "https://discord.com/api/v9/auth/register"

	var jsonStr = []byte(`{"email":"` + email + `","username":"asdsadsad","password":"q1e31e12r13*","invite":null,"consent":true,"date_of_birth":"1973-05-09","gift_code_sku_id":null,"captcha_key":null,"promotional_email_opt_in":false}`)

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr)) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	r.Header.Add("X-Debug-Options", "bugReporterEnabled")

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode == 400 {
		body, _ := ioutil.ReadAll(res.Body)
		var response discordResponse
		json.Unmarshal(body, &response)
		if len(response.Errors.Email.Errors) > 0 {
			if response.Errors.Email.Errors[0].Code == "EMAIL_ALREADY_REGISTERED" {
				fmt.Println("Discord \U0001f440")
			} else {
				if showFalse {
					fmt.Println("Discord", color.RedString(" [Not here!]"))
				}
			}
		} else {
			if showFalse {
				fmt.Println("Discord", color.RedString(" [Not here!]"))
			}
		}
	} else if res.StatusCode == 429 {
		color.Red("Too many requests to Discord!")
	} else {
		color.Red("Couldn't check Discord!")
	}
}
