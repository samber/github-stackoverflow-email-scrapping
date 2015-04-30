
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"math/rand"
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
)


func init() {
	// we use a random value to calculate the time between two request on github
	rand.Seed(time.Now().UTC().UnixNano())
}




/****************************
* gets all top github repos *
****************************/
func scrape_repos() {
	// let's scrape most starred repos
	fmt.Println("Search repos by stars")
	for i := 0; true; i++ {
		min_stars, _ := get_searchPage_repo("stars", i)
		// bug -> retry
		if min_stars == -1 {
			i--
			continue;
		}
		if min_stars < config.Github_scrapping.Repo_min_stars {
			fmt.Println("Search repos by stars - Page: " + strconv.Itoa(i) + " - Min star: " + strconv.Itoa(min_stars))
			break
		}
	}

	// let's scrape most starred repos
	fmt.Println("Search repos by forks")
	for i := 0; true; i++ {
		_, min_forks := get_searchPage_repo("forks", i)
		// bug -> retry
		if min_forks == -1 {
			i--
			continue;
		}
		if min_forks < config.Github_scrapping.Repo_min_forks {
			break
		}
	}
}


/**
 * gets all repos in a search by "forks" or "stars"
*/
func get_searchPage_repo(orderBy string, page int) (int, int) {
	// check param errors
	if orderBy != "stars" && orderBy != "forks" {
		fmt.Println("get_searchPage_repo must have 'stars' or 'forks' as first parameter")
		os.Exit(1)
	}
	if page <= 0 {
		page = 1
	}

	// get the page
	url := "https://github.com/search?o=desc&q=stars%3A%3E1&s=" + orderBy + "&type=Repositories&p=" + strconv.Itoa(page)
	doc := get_html_page(url)
	if doc == nil {
		fmt.Println("Failed to parse page " + strconv.Itoa(page))
		return -1, -1
	}

	// need to return min value to stop scrapping
	var min_stars = 0
	var min_forks = 0

	// html parsing
	repo_list := doc.Find(".repo-list-item")
	if repo_list == nil {
		fmt.Println("Failed to parse page " + strconv.Itoa(page))
		return -1, -1
	}

	repo_list.Each(func (i int, repo *goquery.Selection) {
		var nbr_stars = 0
		var nbr_forks = 0
		var repo_path = ""
		var repo_name = ""
		var repo_owner = ""

		// find forks and stars numbers in the DOM
		stat_items := repo.Find(".repo-list-stat-item")
		stat_items.Each(func (j int, stat_item *goquery.Selection) {
			statItem_value, exists := stat_item.Attr("aria-label")
			if exists == true {
				stat_item := strings.TrimSpace(stat_item.Text())
				stat_item = strings.Replace(stat_item, ",", "", 3)
				nbr, err := strconv.Atoi(stat_item)
				if err == nil && statItem_value == "Stargazers" {
					nbr_stars = nbr
				}
				if err == nil && statItem_value == "Forks" {
					nbr_forks = nbr
				}
			}
		})

		// find repo name and owner
		repo_path = repo.Find(".repo-list-name a").Text()
		repo_owner = strings.Split(repo_path, "/")[0]
		repo_name = strings.Split(repo_path, "/")[1]

		// update min values to stop the scrapping
		if nbr_stars > 0 && (min_stars == 0 || min_stars > nbr_stars) {
			min_stars = nbr_stars
		}
		if nbr_forks > 0 && (min_forks == 0 || min_forks > nbr_forks) {
			min_forks = nbr_forks
		}

		persist_scrapped_repo(repo_path, repo_owner, repo_name, nbr_stars, nbr_forks)
	})

	return min_stars, min_forks
}








/***************************
* gets users by repo owner *
***************************/
func scrape_repos_owner() {
	var repos []GRepo
	pg_findAll(&repos, "")

	fmt.Println("Search repos owners (user/org)")

	// for each repo in the db, we get the owner
	for i := 0; i < len(repos); i++ {
		scrape_repo_owner(repos[i].Owner)
	}
}

func scrape_repo_owner(owner string) {
	doc := get_html_page("https://github.com/" + owner)

	// org ? user ?
	value, _ := doc.Find("body").Attr("class")
	if strings.Index(value, "org") != -1 {
		scrape_repo_owner_parse_org(owner, doc)
	} else if strings.Index(value, "page-profile") != -1 {
		scrape_repo_owner_parse_user(owner, doc)
	}
}

/**
 * repo owner can be an organisation... *
*/
func scrape_repo_owner_parse_org(owner string, doc *goquery.Document) {
	name := owner
	fullname := doc.Find(".org-header-info .js-username").Text()
	email := doc.Find(".org-header-meta.has-email [itemprop='email']").Text()
	link := doc.Find(".org-header-meta.has-blog [itemprop='url']").Text()
	persist_scrapped_orga(name, fullname, email, link)
}

/**
 * ...or a user
*/
func scrape_repo_owner_parse_user(owner string, doc *goquery.Document) {
	username := owner
	fullname := doc.Find(".vcard-names .vcard-fullname").Text()
	email := doc.Find(".vcard-detail .email").Text()
	link := doc.Find(".vcard-detail .url").Text()
	starred, _ := strconv.Atoi(doc.Find(".vcard-stat-count").Eq(1).Text())
	persist_scrapped_user(username, fullname, email, link, starred)
}






/***********************************
 * gets users by repo contributors *
***********************************/
func scrape_repos_contributors() {
	var repos []GRepo
	pg_findAll(&repos, "")

	fmt.Println("Search repos contributors")

	for i := 0; i < len(repos); i++ {
		// get contributor list
		var dataContributors []DataContributor
		get_json_page("https://github.com/" + repos[i].Path + "/graphs/contributors-data", &dataContributors)

		// get each user page of the contributor
		for j := 0; j < len(dataContributors); j++ {
			scrape_repo_owner(dataContributors[j].Author.Login)
		}
	}
}

// github use an api request to get a json list of contributors
// DataContributor is a partial mapping of contributors
type DataContributor struct {
	Author	struct {
		Login	string
	}
	Total int

}



/*********************
 * gets orga members *
*********************/
func scrape_orga_members() {
	var orgs []GOrga
	pg_findAll(&orgs, "")

	fmt.Println("Search organization members (only orgs that own top repos)")

	// each org in the db
	for i := 0; i < len(orgs); i++ {
		// each page of the member list
		for page := 1; true; page++ {
			doc := get_html_page("https://github.com/orgs/" + orgs[i].Name + "/people?page=" + strconv.Itoa(page))

			// parse team member list
			doc.Find(".member-listing .member-list-item").Each(func (i int, member *goquery.Selection) {
				username := member.Find(".member-username").Text()
				scrape_repo_owner(username)
			})

			// to check if another page exists or not
			nbr_members, _ := strconv.Atoi(doc.Find("nav.pagehead-nav .counter").Text())
			if page * 30 >= nbr_members {
				break;
			}
		}
	}
}





/*************************************************
 * Persists Users, Organizations or Repositories *
**************************************************/

// repos
func persist_scrapped_repo(repo_path string, repo_owner string, repo_name string, starred int, forked int) {
	gRepo := GRepo{
		Path: repo_path,
		Owner: repo_owner,
		Name: repo_name,
		Starred: starred,
		Forked: forked,
	}
	pg_create(&gRepo)
}

// organizations
func persist_scrapped_orga(name string, fullname string, email string, link string) {
	gOrga := GOrga{
		Name: name,
		Fullname: fullname,
		Email: email,
		Link: link,
	}
	pg_create(&gOrga)
}

// users
func persist_scrapped_user(username string, fullname string, email string, link string, starred int) {
	gUser := GUser{
		Username: username,
		Fullname: fullname,
		Email: email,
		Link: link,
		Starred: starred,
	}
	pg_create(&gUser)
}




/***************************************
* Utils fonctions for github scrapping *
***************************************/

// sleep between requests at least 5 seconds - till 20
// small hack to avoid github limit rate
func random_sleep_time() time.Duration {
	min := config.Github_scrapping.Min_time_between_requests
	max := config.Github_scrapping.Max_time_between_requests

	var random = rand.Intn(max - min)
	fmt.Println("Wait " + strconv.Itoa(random + min) + " seconds")
	return time.Second * time.Duration(random + min);
}

// gets html page
func get_html_page(url string) *goquery.Document {
	resp := get_page(url)

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return doc
	}

	return doc
}

// gets json document
func get_json_page(url string, model interface{}) {
	doc := get_page(url)

	// get json format
	body, _ := ioutil.ReadAll(doc.Body)
	json.Unmarshal(body, model)
}

// gets page with fake user-agent (small github hack)
func get_page(url string) *http.Response {
	fmt.Println("Get url: " + url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0")
	resp, _ := client.Do(req)

	// github hack to avoid flood and get a limit rate
	time.Sleep(random_sleep_time())

	return resp
}



