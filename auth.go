package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/automattic/go/jaguar"
	"github.com/valyala/fastjson"
)

// runSetup prompts the user for the necessary info to
// configure and run. It can be triggered directly using
// --init or will get triggered if testSetup fails
func runSetup() {
	var user, pass string

	// prompt user for site
	conf.SiteURL = promptForURL("Enter URL for site: ")

	// prompt for username
	fmt.Print("Enter username: ")
	_, err := fmt.Scanf("%s", &user)
	if err != nil {
		log.Fatal("What happened?", err)
	}

	// prompt for password
	fmt.Print("Enter password: ")
	_, err = fmt.Scanf("%s", &pass)
	if err != nil {
		log.Fatal("What happened?", err)
	}

	// make JWT call to fetch token
	url := strings.Join([]string{conf.SiteURL, "wp-json", "jwt-auth/v1/token"}, "/")
	j := jaguar.New()
	j.Url(url)
	j.Params.Add("username", user)
	j.Params.Add("password", pass)
	resp, err := j.Method("POST").Send()
	if err != nil {
		log.Fatal("API error authentication", err)
	}

	if resp.StatusCode == 403 {
		log.Fatal("Error authenticating, try again.")
	}

	if resp.StatusCode == 404 {
		log.Fatal("Auth API not found. JWT Auth plugin installed and activated?")
	}

	conf.Token = fastjson.GetString(resp.Bytes, "data", "token")
	if conf.Token == "" {
		log.Fatal("Unable to find token in json response")
	}

	if conf.Token == "" {
		log.Fatal("No authentication token.", resp.StatusCode, string(resp.Bytes))
	}

	// write out config
	jsonConf, err := json.Marshal(conf)
	if err != nil {
		log.Warn("JSON Encoding Error", err)
	} else {
		err = ioutil.WriteFile(configfile, jsonConf, 0644)
		if err != nil {
			log.Warn("Error writing config file "+configfile, err)
		} else {
			log.Debug("Wrote config file " + configfile)
		}
	}
}

// testSetup confirms everything is configured and working
// includes the local directories, blog config, and auth
func testSetup() bool {
	if conf.SiteURL == "" {
		log.Warn("Site URL not set")
		return false
	}

	if conf.Token == "" {
		log.Warn("Authentication token not set")
		return false
	}

	j := getApiFetcher("jwt-auth/v1/token/validate")
	resp, err := j.Method("POST").Send()
	if err != nil {
		log.Warn("Error in Auth validation API", err)
	}

	if resp.StatusCode == 403 {
		log.Warn("Authentication error. Try running --init ", string(resp.Bytes))
		return false
	}

	return true
}

func promptForURL(prompt string) string {
	var input string

	fmt.Print(prompt)
	_, err := fmt.Scanf("%s", &input)
	if err != nil {
		log.Fatal("What happened?", err)
	}

	input = strings.TrimSuffix(input, "/")

	_, err = url.ParseRequestURI(input)
	if err != nil {
		log.Warn("Error with URL. Be sure to include http:// prefix")
		return promptForURL(prompt)
	}

	return input
}
