package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"os"
)

func makeTables(db *sql.DB, database string) {

	/*db, err := sql.Open("mysql", user+":"+password+"@"+address+"/"+database)
	if err != nil {
		fmt.Println("error", err)
	}
	defer db.Close()*/

	fmt.Println(db.Ping())

	recipeChk, err := db.Query("Select COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'gw2' AND TABLE_NAME = 'recipes'")
	if err != nil {
		fmt.Println("Error", err)
	} else {
		for recipeChk.Next() {
			var (
				count int
			)
			err = recipeChk.Scan(&count)
			if count == 0 {
				fmt.Println("make recipes table")
				makeRecipeTbl, err := db.Query("CREATE TABLE `" + database + "`.`recipes` ( `id` INT NOT NULL , `recipe` TEXT NOT NULL ,  `mat_1_id` INT, `count_1` INT, `mat_2_id` INT, `count_2` INT, `mat_3_id` INT, `count_3` INT, `mat_4_id` INT, `count_4` INT, `amount_created` INT NOT NULL , `individual_sell_price` INT,  `individual_buy_price` INT, `recipe_sell_price` INT, `recipe_buy_price` INT, PRIMARY KEY (`id`))")
				if err != nil {
					fmt.Println("error ", err)
				} else {
					fmt.Println(makeRecipeTbl)
					makeRecipeTbl.Close()
				}
			} else {
				fmt.Println("recipes table already exists")
			}

		}
	}
	recipeChk.Close()
	materialsChk, err := db.Query("Select COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'gw2' AND TABLE_NAME = 'materials'")
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("gotanswer ")
		for materialsChk.Next() {
			var (
				count int64
			)
			err = materialsChk.Scan(&count)
			if count == 0 {
				fmt.Println("make materials table")
				makeMaterialsTbl, err := db.Query("CREATE TABLE `" + database + "`.`materials` ( `id` INT NOT NULL , `sell_price` INT, `buy_price` INT, PRIMARY KEY (`id`))")
				if err != nil {
					fmt.Println("error ", err)
				} else {
					fmt.Println(makeMaterialsTbl)
					makeMaterialsTbl.Close()
				}
			} else {
				fmt.Println("materials table exists")
			}
		}
	}
	materialsChk.Close()
	itemsChk, err := db.Query("Select COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'gw2' AND TABLE_NAME = 'items'")
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("gotanswer ")
		for itemsChk.Next() {
			var (
				count int64
			)
			err = itemsChk.Scan(&count)
			if count == 0 {
				fmt.Println("make materials table")
				makeItemsTbl, err := db.Query("CREATE TABLE `" + database + "`.`items` ( `id` INT NOT NULL , `name` TEXT, PRIMARY KEY (`id`))")
				if err != nil {
					fmt.Println("error ", err)
				} else {
					fmt.Println(makeItemsTbl)
					makeItemsTbl.Close()
				}
			} else {
				fmt.Println("items table exists")
			}
		}
	}
	itemsChk.Close()
	tpUpdateChk, err := db.Query("Select COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'gw2' AND TABLE_NAME = 'tp_update'")
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("gotanswer ")
		for tpUpdateChk.Next() {
			var (
				count int64
			)
			err = tpUpdateChk.Scan(&count)
			if count == 0 {
				fmt.Println("make materials table")
				makeTPUpdateTbl, err := db.Query("CREATE TABLE `" + database + "`.`tp_update` ( `id` INT NOT NULL, `last_update` TIMESTAMP NOT NULL )")
				if err != nil {
					fmt.Println("error ", err)
				} else {
					fmt.Println(makeTPUpdateTbl)
					makeTPUpdateTbl.Close()
					firstDate, err := db.Query("INSERT INTO `tp_update` (`last_update`) VALUES ('0000-00-00 00:00:00')")
					if err != nil {
						fmt.Println("error", err)
					} else {
						firstDate.Close()
					}

				}
			} else {
				fmt.Println("tp_update table exists")
			}
		}
	}
	tpUpdateChk.Close()

}

func updateTimestamp(db *sql.DB) {

	//db, err := sql.Open("mysql", user+":"+password+"@"+address+"/"+database)
	//if err != nil {
	//	fmt.Println("error", err)
	//}
	//defer db.Close()

	updateTStmp, err := db.Query("UPDATE `tp_update` SET `last_update`=CURRENT_TIMESTAMP WHERE `id` = '0'")
	if err != nil {
		fmt.Println("error", err)
	} else {
		updateTStmp.Close()
	}

}
