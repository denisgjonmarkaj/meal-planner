package main

import (
    "log"
    "math"
    "math/rand"
    "time"
)

// Strutture dei dati ampliate
type FoodRules struct {
    Name            string   `json:"name"`
    StandardPortion float64  `json:"standardPortion"`
    Unit            string   `json:"unit"`
    CaloriesPer100g float64  `json:"caloriesPer100g"`
    Category        string   `json:"category"` // protein, carb, fat, vegetable, fruit, beverage
    MealTypes       []string `json:"mealTypes"`
    MinPortion      float64  `json:"minPortion"`
    MaxPortion      float64  `json:"maxPortion"`
    Required        bool     `json:"required"` // se è un componente obbligatorio del pasto
    Frequency       int      `json:"frequency"` // quante volte può apparire in un giorno
}

type MealRules struct {
    RequiredCategories []string          // categorie che devono essere presenti nel pasto
    CategoryLimits    map[string]float64 // limiti calorici per categoria
    MinProtein        float64            // grammi minimi di proteine
    MinCarbs          float64            // grammi minimi di carboidrati
    MaxFat            float64            // grammi massimi di grassi
}

var mealRules = map[string]MealRules{
    "colazione": {
        RequiredCategories: []string{"beverage", "carb", "protein"},
        CategoryLimits: map[string]float64{
            "carb":     200,
            "protein":  150,
            "fat":      100,
            "beverage": 50,
        },
        MinProtein: 15,
        MinCarbs:   30,
        MaxFat:     15,
    },
    "pranzo": {
        RequiredCategories: []string{"carb", "protein", "vegetable"},
        CategoryLimits: map[string]float64{
            "carb":      300,
            "protein":   250,
            "vegetable": 100,
            "fat":       150,
        },
        MinProtein: 30,
        MinCarbs:   60,
        MaxFat:     25,
    },
    "cena": {
        RequiredCategories: []string{"protein", "vegetable", "carb"},
        CategoryLimits: map[string]float64{
            "protein":   300,
            "vegetable": 100,
            "carb":      200,
            "fat":       100,
        },
        MinProtein: 35,
        MinCarbs:   45,
        MaxFat:     20,
    },
}

func generateMealWithUserIngredients(mealType string, userIngredients []string, targetCalories float64) Meal {
    rand.Seed(time.Now().UnixNano())
    
    var items []Food
    var totalCalories float64 = 0
    rules := mealRules[mealType]
    
    // 1. Prima aggiungi gli elementi obbligatori per il tipo di pasto
    for _, category := range rules.RequiredCategories {
        if !containsCategory(items, category) {
            // Cerca tra gli ingredienti dell'utente prima
            added := false
            for _, ing := range userIngredients {
                if rule, exists := foodRules[ing]; exists {
                    if rule.Category == category && isAppropriateForMeal(rule, mealType) {
                        calories := calculateCalories(rule.StandardPortion, rule.CaloriesPer100g)
                        if totalCalories + calories <= targetCalories {
                            items = append(items, Food{
                                Name:     rule.Name,
                                Quantity: rule.StandardPortion,
                                Unit:     rule.Unit,
                                Calories: math.Round(calories),
                            })
                            totalCalories += calories
                            added = true
                            break
                        }
                    }
                }
            }
            
            // Se non trovato tra gli ingredienti dell'utente, aggiungi uno standard
            if !added {
                addStandardCategoryItem(&items, &totalCalories, category, mealType, targetCalories)
            }
        }
    }
    
    // 2. Aggiungi altri ingredienti dell'utente se appropriati e nei limiti calorici
    remainingCalories := targetCalories - totalCalories
    if remainingCalories > 0 {
        for _, ing := range userIngredients {
            if rule, exists := foodRules[ing]; exists {
                if isAppropriateForMeal(rule, mealType) {
                    calories := calculateCalories(rule.StandardPortion, rule.CaloriesPer100g)
                    if totalCalories + calories <= targetCalories {
                        // Verifica i limiti per categoria
                        if categoryLimit, ok := rules.CategoryLimits[rule.Category]; ok {
                            currentCatCals := getCategoryCalories(items, rule.Category)
                            if currentCatCals + calories <= categoryLimit {
                                items = append(items, Food{
                                    Name:     rule.Name,
                                    Quantity: rule.StandardPortion,
                                    Unit:     rule.Unit,
                                    Calories: math.Round(calories),
                                })
                                totalCalories += calories
                            }
                        }
                    }
                }
            }
        }
    }

    log.Printf("Generated meal %s with %.2f calories (target: %.2f)", 
        mealType, totalCalories, targetCalories)
    
    return Meal{
        Items:    items,
        Calories: math.Round(totalCalories),
    }
}

func addStandardCategoryItem(items *[]Food, totalCalories *float64, category string, mealType string, targetCalories float64) {
    // Mappa degli alimenti standard per categoria e tipo di pasto
    standardItems := map[string]map[string]string{
        "colazione": {
            "beverage": "caffe",
            "carb":    "panbauletto",
            "protein": "prosciutto_cotto",
        },
        "pranzo": {
            "carb":      "riso_venere",
            "protein":   "tonno_naturale",
            "vegetable": "zucchine",
        },
        "cena": {
            "protein":   "tacchino",
            "vegetable": "broccoli",
            "carb":      "pane_integrale",
        },
    }

    if mealStandards, ok := standardItems[mealType]; ok {
        if standardItem, ok := mealStandards[category]; ok {
            if rule, exists := foodRules[standardItem]; exists {
                calories := calculateCalories(rule.StandardPortion, rule.CaloriesPer100g)
                if *totalCalories + calories <= targetCalories {
                    *items = append(*items, Food{
                        Name:     rule.Name,
                        Quantity: rule.StandardPortion,
                        Unit:     rule.Unit,
                        Calories: math.Round(calories),
                    })
                    *totalCalories += calories
                }
            }
        }
    }
}

func isAppropriateForMeal(rule FoodRules, mealType string) bool {
    for _, allowedMeal := range rule.MealTypes {
        if allowedMeal == mealType {
            return true
        }
    }
    return false
}