package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var user = flag.String("u", os.Getenv("MATDBUSER"), "username for the materials database")
	var pass = flag.String("p", os.Getenv("MATDBPASS"), "password for the materials database")
	var database = flag.String("db", os.Getenv("MATDB"), "database name for the materials database")
	var address = flag.String("a", os.Getenv("MATDBADDR"), "address for the materials database")
	flag.Parse()
	fmt.Println(*user)
	fmt.Println(*pass)
	fmt.Println(*database)
	fmt.Println(*address)

	makeTables(*user, *pass, *address, *database)

}
