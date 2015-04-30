

package main


import (
)



func main() {
	pg_connect()

	//pg_auto_migrate()		// dev mode
	//pg_drop_tables()		// dev mode
	//pg_create_tables()

	scrape_repos()			// use github search pages (ordering by stars and forks number)
	scrape_repos_contributors()	// MUST call scrape_repos() before !
	scrape_repos_owner()		// MUST call scrape_repos() before !
	scrape_orga_members()		// MUST call scrape_repos_owners() before !
}

