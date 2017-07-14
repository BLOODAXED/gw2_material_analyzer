package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yasvisu/gw2api"
	"math"
	"os"
	"strconv"
	"time"
)

func FetchMaterialPrices(db *sql.DB) {

	api := gw2api.NewGW2Api()

	idList := []int{}

	materialIds, err := db.Query("SELECT `id` FROM `materials`")
	if err != nil {
		fmt.Println("Material Id query Error:", err)
	} else {
		for materialIds.Next() {
			var (
				id int
			)
			err := materialIds.Scan(&id)
			if err != nil {
				fmt.Println("Recipe Id scan Error:", err)
			} else {
				idList = append(idList, id)
			}
		}
		materialIds.Close()
	}
	allMaterialPrices := []gw2api.ArticlePrice{}
	queue := make(chan []gw2api.ArticlePrice)
	remaining := math.Ceil(float64(len(idList)) / 200)
	//concurrently get api data
	for i := 0; i <= len(idList); i += 200 {
		go func(i int) {
			if i+200 > len(idList) {
				prices, err := api.CommercePriceIds(idList[i:len(idList)]...)
				if err != nil {
					fmt.Println("Failed to get recipes "+strconv.Itoa(idList[i])+" to "+strconv.Itoa(len(idList)), err)
					os.Exit(9)
				} else {
					queue <- []gw2api.ArticlePrice(prices[:])
				}
			} else {
				prices, err := api.CommercePriceIds(idList[i : i+200]...)
				if err != nil {
					fmt.Println("Failed to get recipes "+strconv.Itoa(idList[i])+" to "+strconv.Itoa(idList[i+200]), err)
					os.Exit(9)
				} else {
					queue <- []gw2api.ArticlePrice(prices[:])
				}
			}
		}(i)
	}

	for t := range queue {
		allMaterialPrices = append(allMaterialPrices, t...)
		if remaining--; remaining == 0 {
			close(queue)
		}
	}
	for _, x := range allMaterialPrices {
		go func(x gw2api.ArticlePrice) {
			updatePrice, err := db.Query("UPDATE `materials` SET `materials`.`sell_price`=?, `materials`.`buy_price`=? WHERE `materials`.`id`=?", x.Sells.UnitPrice, x.Buys.UnitPrice, x.ID)
			if err != nil {
				fmt.Println("Failed to update material ID:"+string(x.ID)+"because", err)
			} else {
				updatePrice.Close()
			}
		}(x)
	}
	fmt.Println("material end")
}

func FetchRecipePrices(db *sql.DB) {

	api := gw2api.NewGW2Api()

	idList := []int{}

	recipeIds, err := db.Query("SELECT `recipe` FROM `recipes`")
	if err != nil {
		fmt.Println("Recipe Id query Error:", err)
	} else {
		for recipeIds.Next() {
			var (
				id int
			)
			err := recipeIds.Scan(&id)
			if err != nil {
				fmt.Println("Recipe Id scan Error:", err)
			} else {
				idList = append(idList, id)
			}
		}
		recipeIds.Close()
	}
	allRecipePrices := []gw2api.ArticlePrice{}
	queue := make(chan []gw2api.ArticlePrice)
	remaining := math.Ceil(float64(len(idList)) / 200)
	//concurrently get api data
	for i := 0; i <= len(idList); i += 200 {
		go func(i int) {
			if i+200 > len(idList) {
				prices, err := api.CommercePriceIds(idList[i:len(idList)]...)
				if err != nil {
					fmt.Println("Failed to get recipes "+strconv.Itoa(idList[i])+" to "+strconv.Itoa(len(idList)), err)
					os.Exit(9)
				} else {
					queue <- []gw2api.ArticlePrice(prices[:])
				}
			} else {
				prices, err := api.CommercePriceIds(idList[i : i+200]...)
				if err != nil {
					fmt.Println("Failed to get recipes "+strconv.Itoa(idList[i])+" to "+strconv.Itoa(idList[i+200]), err)
					os.Exit(9)
				} else {
					queue <- []gw2api.ArticlePrice(prices[:])
				}
			}
		}(i)
	}
	db.SetMaxOpenConns(100)
	for t := range queue {
		allRecipePrices = append(allRecipePrices, t...)
		if remaining--; remaining == 0 {
			close(queue)
		}
	}
	for _, x := range allRecipePrices {
		go func(x gw2api.ArticlePrice) {
			updatePrice, err := db.Query("UPDATE `recipes` SET `recipes`.`individual_sell_price`=?, `recipes`.`individual_buy_price`=? WHERE `recipes`.`recipe`=?", x.Sells.UnitPrice, x.Buys.UnitPrice, x.ID)
			if err != nil {
				fmt.Println("Failed to update single price recipe output ID:"+strconv.Itoa(x.ID)+"because", err)
			} else {
				updatePrice.Close()
				//fmt.Println("close update price")
				updateFullPrice, err := db.Query("UPDATE `recipes` SET `recipes`.`recipe_sell_price`=`recipes`.`individual_sell_price`*`recipes`.`amount_created`, `recipes`.`recipe_buy_price`=`recipes`.`individual_buy_price`*`recipes`.`amount_created` WHERE `recipes`.`recipe`=?", x.ID)
				if err != nil {
					fmt.Println("Failed to update all price recipe output ID: "+strconv.Itoa(x.ID)+" because", err)
				} else {
					updateFullPrice.Close()
				}
			}
		}(x)
	}
	fmt.Println("recipe end")
}

func LastUpdate(db *sql.DB) bool {

	lessThanTenMinutes := false
	lastTime, err := db.Query("SELECT `last_update` FROM `tp_update` WHERE `id`=0")
	if err != nil {
		fmt.Println("Error", err)
	} else {
		for lastTime.Next() {
			var (
				date string
			)
			err := lastTime.Scan(&date)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				lastUpdate, _ := time.Parse("2006-01-02 15:04:05", date)
				duration := time.Since(lastUpdate)
				if duration.Minutes() >= 10 {
					lessThanTenMinutes = true
				}
			}
		}
		defer lastTime.Close()
	}
	return lessThanTenMinutes
}
