

APP=		github_stack_email_scrapping

SRC=		src/app.go				\
		src/config.go				\
		src/pg.go				\
		src/models.go				\
		src/scrapper.go

PKG=		github.com/lib/pq			\
		github.com/PuerkitoBio/goquery		\
		github.com/jinzhu/gorm



GOPATH := ${PWD}/pkg:${GOPATH}
export GOPATH

default:build

build:
	go build -v -o ./bin/$(APP) $(SRC)

fmt:
	go fmt $(SRC)

run:	build
	./bin/$(APP)

vendor_clean:
	rm -dRf ./pkg/*

vendor_get:
	GOPATH=${PWD}/pkg go get -d -u -v $(PKG)

vendor_update: vendor_get
	rm -rf `find ./pkg/src -type d -name .git` \
	    && rm -rf `find ./pkg/src -type d -name .hg` \
	    && rm -rf `find ./pkg/src -type d -name .bzr` \
	    && rm -rf `find ./pkg/src -type d -name .svn`
