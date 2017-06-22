package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yasvisu/gw2api"
	"math"
	"os"
)

type ItemDBEntry struct {
	ID   int
	Name string
}

func generateItemDB(db *sql.DB) {

	api := gw2api.NewGW2Api()

	/*db, err := sql.Open("mysql", user+":"+password+"@"+address+"/"+database)
	if err != nil {
		fmt.Println("error", err)
	}
	defer db.Close()*/

	item_codes := []int{}
	uniqueItemIds, err := db.Query("SELECT DISTINCT `recipe` FROM `recipes` WHERE `recipe` != 0 UNION SELECT DISTINCT `mat_1_id` FROM `recipes` WHERE `mat_1_id` != 0 UNION SELECT DISTINCT `mat_2_id` FROM `recipes` WHERE `mat_2_id` != 0 UNION SELECT DISTINCT `mat_3_id` FROM `recipes` WHERE `mat_3_id` != 0 UNION SELECT DISTINCT `mat_4_id` FROM `recipes` WHERE `mat_4_id` != 0")
	if err != nil {
		fmt.Println("error", err)
	} else {
		defer uniqueItemIds.Close()
	}
	for uniqueItemIds.Next() {
		var recipe int
		err := uniqueItemIds.Scan(&recipe)
		if err != nil {
			fmt.Println(err)
		}
		item_codes = append(item_codes, recipe)
	}
	relevantItems := []gw2api.Item{}
	queue := make(chan []gw2api.Item)
	remaining := math.Ceil(float64(len(item_codes)) / 200)
	//concurrently get api data
	for i := 0; i <= len(item_codes); i += 200 {
		go func(i int) {
			if i+200 > len(item_codes) {
				items, err := api.ItemIds("english", item_codes[i:len(item_codes)]...)
				if err != nil {
					fmt.Println("Failed to get items "+string(i)+" to "+string(len(item_codes)), err)
					os.Exit(9)
				} else {
					queue <- []gw2api.Item(items[:])
				}
			} else {
				items, err := api.ItemIds("english", item_codes[i:i+200]...)
				if err != nil {
					fmt.Println("Failed to get items "+string(i)+" to "+string(i+200), err)
					os.Exit(9)
				} else {
					queue <- []gw2api.Item(items[:])
				}
			}
		}(i)
	}

	for t := range queue {
		relevantItems = append(relevantItems, t...)
		if remaining--; remaining == 0 {
			close(queue)
		}
	}
	for _, x := range relevantItems {
		insert, err := db.Query("INSERT INTO `gw2`.`items` (`id`, `name`) VALUES (?,?)", x.ID, x.Name)
		if err != nil {
			fmt.Println(err)
		} else {
			insert.Close()
		}
	}

}

func generateMaterialIdDb(db *sql.DB) {

	//api := gw2api.NewGW2Api()

	/*db, err := sql.Open("mysql", user+":"+password+"@"+address+"/"+database)
	if err != nil {
		fmt.Println("error", err)
	}
	defer db.Close()*/
	//db.SetMaxIdleConns(0)
	//db.SetMaxOpenConns(150)
	//db.SetConnMaxLifetime(0)

	uniqueMaterialIds, err := db.Query("SELECT DISTINCT `mat_1_id` FROM `recipes` WHERE `mat_1_id` != 0 UNION SELECT DISTINCT `mat_2_id` FROM `recipes` WHERE `mat_2_id` != 0 UNION SELECT DISTINCT `mat_3_id` FROM `recipes` WHERE `mat_3_id` != 0 UNION SELECT DISTINCT `mat_4_id` FROM `recipes` WHERE `mat_4_id` != 0")
	if err != nil {
		fmt.Println("error", err)
	}
	defer uniqueMaterialIds.Close()
	for uniqueMaterialIds.Next() {
		var mat_1_id string
		err := uniqueMaterialIds.Scan(&mat_1_id)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(mat_1_id))
		insert, err := db.Query("INSERT INTO `gw2`.`materials` (`id`, `sell price`, `buy_price`) VALUES ('" + string(mat_1_id) + "', 'NULL', 'NULL')")
		if err != nil {
			fmt.Println(err)
		} else {
			insert.Close()
		}
	}

}
