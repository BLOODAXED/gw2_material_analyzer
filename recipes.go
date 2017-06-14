package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yasvisu/gw2api"
)

func generateRecipes(user string, password string, address string, database string) {

	api := gw2api.NewGW2Api()

	db, err := sql.Open("mysql", user+":"+password+"@"+address+"/"+database)
	if err != nil {
		fmt.Println("error", err)
	}
	defer db.Close()

	recipeIds, err := api.Recipes()
	if err != nil {
		fmt.Println("Failed to get recipe ids", err)
	}

	fmt.Println(len(recipeIds))

	//var recipes []gw2api.Recipe
	allRecipeInfo := []gw2api.Recipe{}
	for i := 0; i <= len(recipeIds); i += 200 {
		if i+200 > len(recipeIds) {
			recipes, err := api.RecipeIds(recipeIds[i:len(recipeIds)]...)
			if err != nil {
				fmt.Println("Failed to get recipes "+string(i)+" to "+string(len(recipeIds)), err)
			} else {
				allRecipeInfo = append(allRecipeInfo, []gw2api.Recipe(recipes[:])...)
			}
		} else {
			recipes, err := api.RecipeIds(recipeIds[i : i+200]...)
			if err != nil {
				fmt.Println("Failed to get recipes "+string(i)+" to "+string(i+200), err)
			} else {
				allRecipeInfo = append(allRecipeInfo, []gw2api.Recipe(recipes[:])...)
			}
		}
	}
	for _, x := range allRecipeInfo {
		fmt.Printf("%+v\n", x.OutputItemID)
	}

}
