
Github and Stack-Overflow email scrapping
=========================================

* Warning

It's my first Go program, so please be clement ;-)

* HOW-TO

- Open src/app.go
- Choose what to scrape :

```
	scrape_repos()                // use github search pages (ordering by stars and forks number)
        scrape_repos_contributors()   // MUST call scrape_repos() before !
	scrape_repos_owner()          // MUST call scrape_repos() before !
	scrape_orga_members()         // MUST call scrape_repos_owners() before !
```

* @TODO

- Develop the stack-overflow equivalent
- Remove redundant writes in db
- Store in the db the origin of the scrapping (repo owner ? commiter ? organization member ?) and associate this information to the repo/organization ID
