package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yasvisu/gw2api"
	"math"
	"os"
	//"time"
)

type RecipeDBEntry struct {
	ID, ItemID, FirstMatID, FirstMatCount, SecondMatID, SecondMatCount, ThirdMatID, ThirdMatCount, FourthMatID, FourthMatCount, AmountMade int
}

func generateRecipes(user string, password string, address string, database string) {

	api := gw2api.NewGW2Api()

	db, err := sql.Open("mysql", user+":"+password+"@"+address+"/"+database)
	if err != nil {
		fmt.Println("error", err)
	}
	defer db.Close()

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(150)
	db.SetConnMaxLifetime(0)
	recipeIds, err := api.Recipes()
	if err != nil {
		fmt.Println("Failed to get recipe ids", err)
	}

	fmt.Println(len(recipeIds))

	//var recipes []gw2api.Recipe
	allRecipeInfo := []gw2api.Recipe{}
	var dbEntries []RecipeDBEntry

	queue := make(chan []gw2api.Recipe)
	remaining := math.Ceil(float64(len(recipeIds)) / 200)
	//concurrently get api data
	for i := 0; i <= len(recipeIds); i += 200 {
		go func(i int) {
			if i+200 > len(recipeIds) {
				recipes, err := api.RecipeIds(recipeIds[i:len(recipeIds)]...)
				if err != nil {
					fmt.Println("Failed to get recipes "+string(i)+" to "+string(len(recipeIds)), err)
					os.Exit(9)
				} else {
					queue <- []gw2api.Recipe(recipes[:])
				}
			} else {
				recipes, err := api.RecipeIds(recipeIds[i : i+200]...)
				if err != nil {
					fmt.Println("Failed to get recipes "+string(i)+" to "+string(i+200), err)
					os.Exit(9)
				} else {
					queue <- []gw2api.Recipe(recipes[:])
				}
			}
		}(i)
	}

	for t := range queue {
		allRecipeInfo = append(allRecipeInfo, t...)
		if remaining--; remaining == 0 {
			close(queue)
		}
	}
	fmt.Println(len(allRecipeInfo))
	for _, x := range allRecipeInfo {
		temp := RecipeDBEntry{}
		temp.ID = x.ID
		temp.ItemID = x.OutputItemID
		switch len(x.Ingredients) {
		case 2:
			temp.FirstMatID = x.Ingredients[0].ItemID
			temp.FirstMatCount = x.Ingredients[0].Count
			temp.SecondMatID = x.Ingredients[1].ItemID
			temp.SecondMatCount = x.Ingredients[1].Count
			/*temp.ThirdMatID = nil
			temp.ThirdMatCount = nil
			temp.FourthMatID = nil
			temp.FourthMatCount = nil*/

		case 3:
			temp.FirstMatID = x.Ingredients[0].ItemID
			temp.FirstMatCount = x.Ingredients[0].Count
			temp.SecondMatID = x.Ingredients[1].ItemID
			temp.SecondMatCount = x.Ingredients[1].Count
			temp.ThirdMatID = x.Ingredients[2].ItemID
			temp.ThirdMatCount = x.Ingredients[2].Count
			/*temp.FourthMatID = nil
			temp.FourthMatCount = nil*/

		case 4:
			temp.FirstMatID = x.Ingredients[0].ItemID
			temp.FirstMatCount = x.Ingredients[0].Count
			temp.SecondMatID = x.Ingredients[1].ItemID
			temp.SecondMatCount = x.Ingredients[1].Count
			temp.ThirdMatID = x.Ingredients[2].ItemID
			temp.ThirdMatCount = x.Ingredients[2].Count
			temp.FourthMatID = x.Ingredients[3].ItemID
			temp.FourthMatCount = x.Ingredients[3].Count

		default:
			temp.FirstMatID = x.Ingredients[0].ItemID
			temp.FirstMatCount = x.Ingredients[0].Count
			/*temp.SecondMatID = nil
			temp.SecondMatCount = nil
			temp.ThirdMatID = nil
			temp.ThirdMatCount = nil
			temp.FourthMatID = nil
			temp.FourthMatCount = nil*/
		}
		temp.AmountMade = x.OutputItemCount
		dbEntries = append(dbEntries, temp)
	}
	fmt.Println(len(dbEntries))

	for _, x := range dbEntries {
		rows, err := db.Query("INSERT INTO `gw2`.`recipes` (`id`, `recipe`, `mat_1_id`, `count_1`, `mat_2_id`, `count_2`, `mat_3_id`, `count_3`, `mat_4_id`, `count_4`, `amount_created`) VALUES (?,?,?,?,?,?,?,?,?,?,?)", x.ID, x.ItemID, x.FirstMatID, x.FirstMatCount, x.SecondMatID, x.SecondMatCount, x.ThirdMatID, x.ThirdMatCount, x.FourthMatID, x.FourthMatCount, x.AmountMade)
		if err != nil {
			fmt.Println("error", err)
		}
		rows.Close()
		//time.Sleep(100 * time.Millisecond)
	}

}
