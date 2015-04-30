

package main


import (
	"fmt"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
//	"database/sql"
	_ "github.com/lib/pq"
)



var db gorm.DB



func pg_connect() {
	var url = "postgres://" + config.Db.Username + ":" + config.Db.Password + "@" + config.Db.Hostname + ":" + strconv.Itoa(config.Db.Port) + "/" + config.Db.Database
	tmp_db, err := gorm.Open("postgres", url)
	tmp_db.SingularTable(true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Connected to database")
		db = tmp_db
	}

	// ping
	ping := db.DB().Ping()
	if ping != nil {
		fmt.Println(ping)
		os.Exit(1)
	}
}




// create all tables
func pg_create_tables() {
	db.CreateTable(&GUser{})
	fmt.Println("Created GUser table")
	db.CreateTable(&GOrga{})
	fmt.Println("Created GOrga table")
	db.CreateTable(&GRepo{})
	fmt.Println("Created GRepo table")
}

// drop all tables
// use it only in dev mode
func pg_drop_tables() {
	db.DropTable(&GUser{})
	fmt.Println("Droped GUser table")
	db.DropTable(&GOrga{})
	fmt.Println("Droped GOrga table")
	db.DropTable(&GRepo{})
	fmt.Println("Droped GRepo table")
}
// update mapping and model
func pg_auto_migrate() {
	db.AutoMigrate(&GUser{}, &GOrga{}, &GRepo{})
	fmt.Println("Updated GRepo schema")
}




func pg_create(entry interface{}) {
	db.Create(entry)
}

func pg_findOne(entry interface{}, whereCondition string) {
	if whereCondition == "" {
		db.First(entry)
	} else {
		db.Where(whereCondition).First(entry)
	}
}

func pg_findAll(entries interface{}, whereCondition string) {
	if whereCondition == "" {
		db.Find(entries)
	} else {
		db.Where(whereCondition).Find(entries)
	}
}

