package main

import (
    "log"
    "net/http"
    "math"
    "math/rand"
    "time"
    "sort"
    "github.com/gin-gonic/gin"
)

// Strutture di base
type Food struct {
    Name     string  `json:"name"`
    Quantity float64 `json:"quantity"`
    Unit     string  `json:"unit"`
    Calories float64 `json:"calories"`
}

type Meal struct {
    Items    []Food  `json:"items"`
    Calories float64 `json:"calories"`
}

type MealPlan struct {
    Colazione Meal `json:"colazione"`
    Spuntino  Meal `json:"spuntino"`
    Pranzo    Meal `json:"pranzo"`
    Merenda   Meal `json:"merenda"`
    Cena      Meal `json:"cena"`
}

// Strutture per l'organizzazione degli ingredienti
type IngredientCategory struct {
    Name        string   `json:"name"`
    Ingredients []string `json:"ingredients"`
}

type MealIngredients struct {
    MealName    string               `json:"mealName"`
    Categories  []IngredientCategory `json:"categories"`
}

// Strutture delle regole
type FoodRules struct {
    Name            string   `json:"name"`
    StandardPortion float64  `json:"standardPortion"`
    Unit            string   `json:"unit"`
    CaloriesPer100g float64  `json:"caloriesPer100g"`
    Category        string   `json:"category"` // protein, carb, fat, vegetable, fruit, beverage
    Description     string   `json:"description"`
    MealTypes       []string `json:"mealTypes"`
    MinPortion      float64  `json:"minPortion"`
    MaxPortion      float64  `json:"maxPortion"`
    Required        bool     `json:"required"`
    Frequency       int      `json:"frequency"`
}

type MealRules struct {
    RequiredCategories []string
    CategoryLimits    map[string]float64
    MinProtein        float64
    MinCarbs          float64
    MaxFat            float64
}

// Regole dei pasti
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

// Database delle regole alimentari
// Database delle regole alimentari basato sul PDF
var foodRules = map[string]FoodRules{
    // BEVANDE
    "caffe": {
        Name:            "Caffè",
        StandardPortion: 30,
        Unit:            "g",
        CaloriesPer100g: 1,
        Category:        "beverage",
        Description:     "1 Tazzina",
        MealTypes:       []string{"colazione"},
        MinPortion:      30,
        MaxPortion:      30,
        Required:        true,
        Frequency:       1,
    },
    "spremuta_arancia": {
        Name:            "Spremuta di arancia",
        StandardPortion: 200,
        Unit:            "g",
        CaloriesPer100g: 45,
        Category:        "beverage",
        Description:     "1 Bicchiere",
        MealTypes:       []string{"colazione"},
        MinPortion:      200,
        MaxPortion:      200,
        Required:        false,
        Frequency:       1,
    },
    "ace_diet": {
        Name:            "Ace Diet Hero",
        StandardPortion: 250,
        Unit:            "g",
        CaloriesPer100g: 20,
        Category:        "beverage",
        Description:     "Alternativa alla spremuta",
        MealTypes:       []string{"colazione"},
        MinPortion:      250,
        MaxPortion:      250,
        Required:        false,
        Frequency:       1,
    },

    // CARBOIDRATI
    "panbauletto": {
        Name:            "Panbauletto Integrale",
        StandardPortion: 48,
        Unit:            "g",
        CaloriesPer100g: 270,
        Category:        "carb",
        Description:     "Mulino Bianco",
        MealTypes:       []string{"colazione", "merenda"},
        MinPortion:      48,
        MaxPortion:      48,
        Required:        true,
        Frequency:       2,
    },
    "pane_integrale": {
        Name:            "Pane integrale",
        StandardPortion: 120,
        Unit:            "g",
        CaloriesPer100g: 250,
        Category:        "carb",
        Description:     "4 Fette",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      120,
        MaxPortion:      120,
        Required:        true,
        Frequency:       2,
    },
    "crackers_integrali": {
        Name:            "Crackers integrali",
        StandardPortion: 30,
        Unit:            "g",
        CaloriesPer100g: 430,
        Category:        "carb",
        Description:     "1 Pacchetto",
        MealTypes:       []string{"spuntino", "merenda"},
        MinPortion:      30,
        MaxPortion:      30,
        Required:        false,
        Frequency:       2,
    },
    "riso_venere": {
        Name:            "Riso venere",
        StandardPortion: 100,
        Unit:            "g",
        CaloriesPer100g: 340,
        Category:        "carb",
        Description:     "",
        MealTypes:       []string{"pranzo"},
        MinPortion:      100,
        MaxPortion:      100,
        Required:        false,
        Frequency:       1,
    },
    "riso_basmati": {
        Name:            "Riso basmati",
        StandardPortion: 100,
        Unit:            "g",
        CaloriesPer100g: 350,
        Category:        "carb",
        Description:     "Alternativa al riso venere",
        MealTypes:       []string{"pranzo"},
        MinPortion:      100,
        MaxPortion:      100,
        Required:        false,
        Frequency:       1,
    },
    "pasta_integrale": {
        Name:            "Pasta integrale",
        StandardPortion: 120,
        Unit:            "g",
        CaloriesPer100g: 340,
        Category:        "carb",
        Description:     "",
        MealTypes:       []string{"pranzo"},
        MinPortion:      120,
        MaxPortion:      120,
        Required:        false,
        Frequency:       1,
    },

    // PROTEINE
    "prosciutto_cotto": {
        Name:            "Prosciutto cotto",
        StandardPortion: 50,
        Unit:            "g",
        CaloriesPer100g: 145,
        Category:        "protein",
        Description:     "Alta qualità - sgrassato",
        MealTypes:       []string{"colazione"},
        MinPortion:      50,
        MaxPortion:      50,
        Required:        false,
        Frequency:       1,
    },
    "salmone_affumicato": {
        Name:            "Salmone affumicato",
        StandardPortion: 50,
        Unit:            "g",
        CaloriesPer100g: 217,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"colazione", "pranzo"},
        MinPortion:      50,
        MaxPortion:      100,
        Required:        false,
        Frequency:       1,
    },
    "uova_albume": {
        Name:            "Albume d'uovo",
        StandardPortion: 80,
        Unit:            "g",
        CaloriesPer100g: 52,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"colazione"},
        MinPortion:      80,
        MaxPortion:      80,
        Required:        false,
        Frequency:       1,
    },
    "ricotta_light": {
        Name:            "Ricotta Light",
        StandardPortion: 60,
        Unit:            "g",
        CaloriesPer100g: 146,
        Category:        "protein",
        Description:     "Galbani",
        MealTypes:       []string{"colazione"},
        MinPortion:      60,
        MaxPortion:      60,
        Required:        false,
        Frequency:       1,
    },
    "mozzarella_light": {
        Name:            "Mozzarella Light",
        StandardPortion: 125,
        Unit:            "g",
        CaloriesPer100g: 206,
        Category:        "protein",
        Description:     "Santa Lucia",
        MealTypes:       []string{"pranzo"},
        MinPortion:      125,
        MaxPortion:      250,
        Required:        false,
        Frequency:       1,
    },
    "tonno_naturale": {
        Name:            "Tonno al naturale",
        StandardPortion: 160,
        Unit:            "g",
        CaloriesPer100g: 130,
        Category:        "protein",
        Description:     "Mareblu",
        MealTypes:       []string{"pranzo"},
        MinPortion:      160,
        MaxPortion:      160,
        Required:        false,
        Frequency:       1,
    },
    "petto_pollo": {
        Name:            "Petto di pollo",
        StandardPortion: 250,
        Unit:            "g",
        CaloriesPer100g: 165,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      160,
        MaxPortion:      250,
        Required:        false,
        Frequency:       1,
    },
    "tacchino_petto": {
        Name:            "Tacchino petto",
        StandardPortion: 250,
        Unit:            "g",
        CaloriesPer100g: 104,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"cena"},
        MinPortion:      250,
        MaxPortion:      250,
        Required:        false,
        Frequency:       1,
    },
    "pesce_spada": {
        Name:            "Pesce spada",
        StandardPortion: 250,
        Unit:            "g",
        CaloriesPer100g: 144,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"cena"},
        MinPortion:      250,
        MaxPortion:      250,
        Required:        false,
        Frequency:       1,
    },
    "salmone_fresco": {
        Name:            "Salmone fresco",
        StandardPortion: 200,
        Unit:            "g",
        CaloriesPer100g: 208,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"cena"},
        MinPortion:      200,
        MaxPortion:      200,
        Required:        false,
        Frequency:       1,
    },
    "merluzzo": {
        Name:            "Merluzzo o nasello",
        StandardPortion: 250,
        Unit:            "g",
        CaloriesPer100g: 82,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"cena"},
        MinPortion:      250,
        MaxPortion:      250,
        Required:        false,
        Frequency:       1,
    },
    "orata": {
        Name:            "Orata fresca",
        StandardPortion: 300,
        Unit:            "g",
        CaloriesPer100g: 124,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"cena"},
        MinPortion:      300,
        MaxPortion:      300,
        Required:        false,
        Frequency:       1,
    },

    // LATTICINI E FORMAGGI
    "yogurt_greco": {
        Name:            "Yogurt greco magro alla frutta",
        StandardPortion: 150,
        Unit:            "g",
        CaloriesPer100g: 97,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"colazione", "merenda"},
        MinPortion:      150,
        MaxPortion:      170,
        Required:        false,
        Frequency:       2,
    },
    "grana": {
        Name:            "Grana",
        StandardPortion: 30,
        Unit:            "g",
        CaloriesPer100g: 392,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"pranzo", "merenda"},
        MinPortion:      20,
        MaxPortion:      50,
        Required:        false,
        Frequency:       1,
    },

    // LEGUMI
    "lenticchie": {
        Name:            "Lenticchie secche",
        StandardPortion: 70,
        Unit:            "g",
        CaloriesPer100g: 325,
        Category:        "protein",
        Description:     "",
        MealTypes:       []string{"pranzo"},
        MinPortion:      70,
        MaxPortion:      70,
        Required:        false,
        Frequency:       1,
    },

    // VERDURE
    "broccoli": {
        Name:            "Broccolo",
        StandardPortion: 300,
        Unit:            "g",
        CaloriesPer100g: 34,
        Category:        "vegetable",
        Description:     "a testa",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      300,
        MaxPortion:      300,
        Required:        false,
        Frequency:       2,
    },
    "zucchine": {
        Name:            "Zucchine",
        StandardPortion: 300,
        Unit:            "g",
        CaloriesPer100g: 17,
        Category:        "vegetable",
        Description:     "",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      200,
        MaxPortion:      300,
        Required:        false,
        Frequency:       2,
    },
    "carote": {
        Name:            "Carote",
        StandardPortion: 150,
        Unit:            "g",
        CaloriesPer100g: 41,
        Category:        "vegetable",
        Description:     "",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      150,
        MaxPortion:      200,
        Required:        false,
        Frequency:       2,
    },
    "lattuga": {
        Name:            "Lattuga",
        StandardPortion: 80,
        Unit:            "g",
        CaloriesPer100g: 15,
        Category:        "vegetable",
        Description:     "",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      80,
        MaxPortion:      160,
        Required:        false,
        Frequency:       2,
    },
    "pomodori": {
        Name:            "Pomodori da insalata",
        StandardPortion: 200,
        Unit:            "g",
        CaloriesPer100g: 18,
        Category:        "vegetable",
        Description:     "",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      200,
        MaxPortion:      200,
        Required:        false,
        Frequency:       2,
    },
    "melanzane": {
        Name:            "Melanzane",
        StandardPortion: 200,
        Unit:            "g",
        CaloriesPer100g: 25,
        Category:        "vegetable",
        Description:     "",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      200,
        MaxPortion:      300,
        Required:        false,
        Frequency:       2,
    },
    "funghi": {
        Name:            "Funghi coltivati prataioli",
        StandardPortion: 300,
        Unit:            "g",
        CaloriesPer100g: 22,
        Category:        "vegetable",
        Description:     "",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      300,
        MaxPortion:      300,
        Required:        false,
        Frequency:       2,
    },
    "rucola": {
        Name:            "Rughetta o rucola",
        StandardPortion: 100,
        Unit:            "g",
        CaloriesPer100g: 25,
        Category:        "vegetable",
        Description:     "",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      100,
        MaxPortion:      100,
        Required:        false,
        Frequency:       2,
    },

    // FRUTTA E FRUTTA SECCA
    "frutta_fresca": {
        Name:            "Frutta fresca",
        StandardPortion: 150,
        Unit:            "g",
        CaloriesPer100g: 50,
        Category:        "fruit",
        Description:     "media",
        MealTypes:       []string{"colazione", "spuntino", "merenda"},
        MinPortion:      150,
        MaxPortion:      200,
        Required:        true,
        Frequency:       3,
    },
    "frutta_secca": {
        Name:            "Frutta secca e oleosa",
        StandardPortion: 20,
        Unit:            "g",
        CaloriesPer100g: 600,
        Category:        "fat",
        Description:     "media",
        MealTypes:       []string{"spuntino", "merenda"},
        MinPortion:      20,
        MaxPortion:      30,
        Required:        false,
        Frequency:       1,
    },
}


// Helper functions
func calculateCalories(quantity float64, caloriesPer100g float64) float64 {
    return (quantity * caloriesPer100g) / 100
}

func containsCategory(items []Food, category string) bool {
    for _, item := range items {
        for _, rule := range foodRules {
            if rule.Name == item.Name && rule.Category == category {
                return true
            }
        }
    }
    return false
}

func getCategoryCalories(items []Food, category string) float64 {
    var calories float64 = 0
    for _, item := range items {
        for _, rule := range foodRules {
            if rule.Name == item.Name && rule.Category == category {
                calories += item.Calories
            }
        }
    }
    return calories
}

func containsFood(items []Food, foodName string) bool {
    foodRule, exists := foodRules[foodName]
    if !exists {
        return false
    }
    
    for _, item := range items {
        if item.Name == foodRule.Name {
            return true
        }
    }
    return false
}

func isAppropriateForMeal(rule FoodRules, mealType string) bool {
    for _, allowedMeal := range rule.MealTypes {
        if allowedMeal == mealType {
            return true
        }
    }
    return false
}

func addFoodItem(items *[]Food, totalCalories *float64, rule FoodRules, targetCalories float64) {
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

// Funzione per organizzare gli ingredienti
func organizeIngredients() []MealIngredients {
    mealMap := map[string]map[string][]string{
        "colazione": {
            "Bevande":    {},
            "Carboidrati": {},
            "Proteine":    {},
            "Frutta":      {},
            "Extra":       {},
        },
        "spuntino": {
            "Frutta":     {},
            "Snack":      {},
            "Extra":      {},
        },
        "pranzo": {
            "Carboidrati": {},
            "Proteine":    {},
            "Verdure":     {},
            "Extra":       {},
        },
        "merenda": {
            "Frutta":     {},
            "Snack":      {},
            "Proteine":   {},
            "Extra":      {},
        },
        "cena": {
            "Carboidrati": {},
            "Proteine":    {},
            "Verdure":     {},
            "Extra":       {},
        },
    }

    categoryMapping := map[string]string{
        "beverage":  "Bevande",
        "carb":      "Carboidrati",
        "protein":   "Proteine",
        "vegetable": "Verdure",
        "fruit":     "Frutta",
        "fat":       "Extra",
    }

    // Popola la mappa con gli ingredienti
    for key, rule := range foodRules {
        displayCategory := categoryMapping[rule.Category]
        for _, mealType := range rule.MealTypes {
            if categories, exists := mealMap[mealType]; exists {
                if displayCategory == "" {
                    displayCategory = "Extra"
                }
                categories[displayCategory] = append(
                    categories[displayCategory],
                    key,
                )
            }
        }
    }

    // Converti la mappa in slice per il JSON
    var result []MealIngredients
    mealNames := map[string]string{
        "colazione": "Colazione",
        "spuntino":  "Spuntino",
        "pranzo":    "Pranzo",
        "merenda":   "Merenda",
        "cena":      "Cena",
    }

    // Ordina gli ingredienti all'interno di ogni categoria
    for mealType, categories := range mealMap {
        var mealCategories []IngredientCategory
        for catName, ingredients := range categories {
            if len(ingredients) > 0 {
                sort.Strings(ingredients) // Ordina gli ingredienti alfabeticamente
                mealCategories = append(mealCategories, IngredientCategory{
                    Name:        catName,
                    Ingredients: ingredients,
                })
            }
        }
        // Ordina le categorie per nome
        sort.Slice(mealCategories, func(i, j int) bool {
            return mealCategories[i].Name < mealCategories[j].Name
        })
        result = append(result, MealIngredients{
            MealName:   mealNames[mealType],
            Categories: mealCategories,
        })
    }

    // Ordina i pasti nell'ordine corretto della giornata
    sort.Slice(result, func(i, j int) bool {
        mealOrder := map[string]int{
            "Colazione": 1,
            "Spuntino":  2,
            "Pranzo":    3,
            "Merenda":   4,
            "Cena":      5,
        }
        return mealOrder[result[i].MealName] < mealOrder[result[j].MealName]
    })

    return result
}

// Meal generation function
func generateMealWithUserIngredients(mealType string, userIngredients []string, targetCalories float64) Meal {
    rand.Seed(time.Now().UnixNano())
    
    var items []Food
    var totalCalories float64 = 0
    rules := mealRules[mealType]

    // 1. Aggiungi elementi obbligatori per il tipo di pasto
    for _, category := range rules.RequiredCategories {
        if !containsCategory(items, category) {
            // Cerca prima tra gli ingredienti dell'utente
            found := false
            for _, ing := range userIngredients {
                if rule, exists := foodRules[ing]; exists {
                    if rule.Category == category && isAppropriateForMeal(rule, mealType) {
                        addFoodItem(&items, &totalCalories, rule, targetCalories)
                        found = true
                        break
                    }
                }
            }
            
            // Se non trovato tra gli ingredienti dell'utente, aggiungi uno standard
            if !found {
                var standardFoods []string
                for k, v := range foodRules {
                    if v.Category == category && isAppropriateForMeal(v, mealType) {
                        standardFoods = append(standardFoods, k)
                    }
                }
                if len(standardFoods) > 0 {
                    randomIndex := rand.Intn(len(standardFoods))
                    if rule, exists := foodRules[standardFoods[randomIndex]]; exists {
                        addFoodItem(&items, &totalCalories, rule, targetCalories)
                    }
                }
            }
        }
    }

    // 2. Aggiungi altri ingredienti dell'utente se appropriati
    for _, ing := range userIngredients {
        if rule, exists := foodRules[ing]; exists {
            if isAppropriateForMeal(rule, mealType) && totalCalories < targetCalories {
                if categoryLimit, ok := rules.CategoryLimits[rule.Category]; ok {
                    currentCatCals := getCategoryCalories(items, rule.Category)
                    calories := calculateCalories(rule.StandardPortion, rule.CaloriesPer100g)
                    if currentCatCals + calories <= categoryLimit {
                        addFoodItem(&items, &totalCalories, rule, targetCalories)
                    }
                }
            }
        }
    }

    return Meal{
        Items:    items,
        Calories: math.Round(totalCalories),
    }
}

func main() {
    r := gin.Default()

    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // Routes
    r.GET("/api/ingredients", func(c *gin.Context) {
        c.JSON(http.StatusOK, organizeIngredients())
    })

    r.POST("/api/generate-plan", func(c *gin.Context) {
        var request struct {
            Ingredients    []string `json:"ingredients"`
            TargetCalories int      `json:"targetCalories"`
        }

        if err := c.BindJSON(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        plan := MealPlan{
            Colazione: generateMealWithUserIngredients("colazione", request.Ingredients, float64(request.TargetCalories)*0.25),
            Spuntino:  generateMealWithUserIngredients("spuntino", request.Ingredients, float64(request.TargetCalories)*0.10),
            Pranzo:    generateMealWithUserIngredients("pranzo", request.Ingredients, float64(request.TargetCalories)*0.35),
            Merenda:   generateMealWithUserIngredients("merenda", request.Ingredients, float64(request.TargetCalories)*0.10),
            Cena:      generateMealWithUserIngredients("cena", request.Ingredients, float64(request.TargetCalories)*0.20),
        }

        c.JSON(http.StatusOK, plan)
    })

    log.Println("Server starting on :8080")
    r.Run(":8080")
}