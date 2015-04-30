

package main


import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)



type Postgresql struct {
	Hostname	string
	Port		int
	Username	string
	Password	string
	Database	string
}

type Github_scrapping struct {
	Repo_min_forks	int
	Repo_min_stars	int
	Min_time_between_requests int
	Max_time_between_requests int
}

type Configuration struct {
	Db		Postgresql
	Github_scrapping	Github_scrapping
}

var config = Configuration{}


func init() {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Println("error")
	}

	e = json.Unmarshal(file, &config)
	if e != nil {
		fmt.Println("error")
	}
}
