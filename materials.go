package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func main() {
	var user = flag.String("u", os.Getenv("MATDBUSER"), "username for the materials database")
	var password = flag.String("p", os.Getenv("MATDBPASS"), "password for the materials database")
	var database = flag.String("db", os.Getenv("MATDB"), "database name for the materials database")
	var address = flag.String("a", os.Getenv("MATDBADDR"), "address for the materials database")
	var recipes = flag.Bool("setup", false, "truncate and populate the recipes database")
	flag.Parse()
	fmt.Println(*user)
	fmt.Println(*password)
	fmt.Println(*database)
	fmt.Println(*address)

	db, err := sql.Open("mysql", *user+":"+*password+"@"+*address+"/"+*database)
	if err != nil {
		fmt.Println("error", err)
	}
	defer db.Close()

	makeTables(db, *database)
	if *recipes == true {
		generateRecipes(db)
		generateMaterialIdDb(db)
		generateItemDB(db)
	}

}
