package main

import (
	"log"
	"net/http"
	"math"
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


// Regole per gli alimenti
type FoodRules struct {
    Name              string
    StandardPortion   float64
    Unit              string
    CaloriesPer100g   float64
    AlternativePortion *float64
    Category          string
    Description       string
    MealTypes         []string
    MinPortion        float64
    MaxPortion        float64

}

// Strutture per i pasti tipici
type MealTemplate struct {
	MainOptions       []string
	SideOptions       []string
	RequiredTypes     []string
	StandardStructure []string
}

// Helper per creare puntatori a float64
func pointer(v float64) *float64 {
	return &v
}

// Database delle regole alimentari basato sulle schede esempio
var foodRules = map[string]FoodRules{
    "caffe": {
        Name:            "Caffè",
        StandardPortion: 30,
        Unit:            "g",
        CaloriesPer100g: 1,
        Description:     "1 Tazzina",
        MealTypes:       []string{"colazione"},
        MinPortion:      30,
        MaxPortion:      30,
        Category:        "beverage",
    },
    "panbauletto": {
        Name:            "Panbauletto Integrale",
        StandardPortion: 48,
        Unit:            "g",
        CaloriesPer100g: 270,
        Description:     "Mulino Bianco",
        MealTypes:       []string{"colazione", "merenda"},
        MinPortion:      48,
        MaxPortion:      48,
        Category:        "carb",
    },
    "prosciutto_cotto": {
        Name:            "Prosciutto cotto",
        StandardPortion: 50,
        Unit:            "g",
        CaloriesPer100g: 145,
        Description:     "alta qualità - sgrassato",
        MealTypes:       []string{"colazione"},
        MinPortion:      50,
        MaxPortion:      50,
        Category:        "protein",
    },
    "spremuta_arancia": {
        Name:            "Spremuta di arancia",
        StandardPortion: 200,
        Unit:            "g",
        CaloriesPer100g: 45,
        Description:     "1 Bicchiere",
        MealTypes:       []string{"colazione"},
        MinPortion:      200,
        MaxPortion:      200,
        Category:        "beverage",
    },

    "frutta_fresca": {
        Name:            "Frutta fresca",
        StandardPortion: 150,
        Unit:            "g",
        CaloriesPer100g: 50,
        Description:     "media",
        MealTypes:       []string{"colazione", "spuntino", "merenda"},
        MinPortion:      150,
        MaxPortion:      200,
        Category:        "fruit",
    },
    "crackers_integrali": {
        Name:            "Crackers integrali",
        StandardPortion: 30,
        Unit:            "g",
        CaloriesPer100g: 430,
        Description:     "1 Pacchetto",
        MealTypes:       []string{"spuntino"},
        MinPortion:      30,
        MaxPortion:      30,
        Category:        "carb",
    },
    "lenticchie": {
        Name:            "Lenticchie secche",
        StandardPortion: 70,
        Unit:            "g",
        CaloriesPer100g: 325,
        Description:     "secche",
        MealTypes:       []string{"pranzo"},
        MinPortion:      70,
        MaxPortion:      70,
        Category:        "protein",
        AlternativePortion: pointer(200.0), // alternativa: lenticchie bollite
    },
    "riso_venere": {
        Name:            "Riso venere",
        StandardPortion: 100,
        Unit:            "g",
        CaloriesPer100g: 340,
        MealTypes:       []string{"pranzo"},
        MinPortion:      100,
        MaxPortion:      100,
        Category:        "carb",
        AlternativePortion: pointer(100.0), // alternative: riso basmati o integrale
    },
    "yogurt_greco": {
        Name:            "Yogurt greco magro alla frutta",
        StandardPortion: 150,
        Unit:            "g",
        CaloriesPer100g: 97,
        MealTypes:       []string{"merenda"},
        MinPortion:      150,
        MaxPortion:      150,
        Category:        "protein",
    },
    "frutta_secca": {
        Name:            "Frutta secca e oleosa",
        StandardPortion: 20,
        Unit:            "g",
        CaloriesPer100g: 600,
        Description:     "media",
        MealTypes:       []string{"spuntino", "merenda"},
        MinPortion:      20,
        MaxPortion:      30,
        Category:        "fat",
    },
    "pane_integrale": {
        Name:            "Pane integrale",
        StandardPortion: 120,
        Unit:            "g",
        CaloriesPer100g: 250,
        Description:     "4 Fette",
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      120,
        MaxPortion:      120,
        Category:        "carb",
    },
    "broccoli": {
        Name:            "Broccolo",
        StandardPortion: 300,
        Unit:            "g",
        CaloriesPer100g: 34,
        Description:     "a testa",
        MealTypes:       []string{"cena"},
        MinPortion:      300,
        MaxPortion:      300,
        Category:        "vegetable",
    },
    "pesce_spada": {
        Name:            "Pesce spada",
        StandardPortion: 250,
        Unit:            "g",
        CaloriesPer100g: 144,
        MealTypes:       []string{"cena"},
        MinPortion:      250,
        MaxPortion:      250,
        Category:        "protein",
    },
    "salmone_affumicato": {
        Name:            "Salmone affumicato",
        StandardPortion: 50,
        Unit:            "g",
        CaloriesPer100g: 217,
        MealTypes:       []string{"colazione", "pranzo"},
        MinPortion:      50,
        MaxPortion:      100,
        Category:        "protein",
    },
    "avocado": {
        Name:            "Avocado",
        StandardPortion: 50,
        Unit:            "g",
        CaloriesPer100g: 160,
        Description:     "1/4 di Frutto",
        MealTypes:       []string{"colazione", "pranzo"},
        MinPortion:      30,
        MaxPortion:      50,
        Category:        "fat",
    },
    "carote": {
        Name:            "Carote",
        StandardPortion: 150,
        Unit:            "g",
        CaloriesPer100g: 41,
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      150,
        MaxPortion:      200,
        Category:        "vegetable",
    },
    "mozzarella_light": {
        Name:            "Mozzarella Light",
        StandardPortion: 125,
        Unit:            "g",
        CaloriesPer100g: 206,
        Description:     "Santa Lucia",
        MealTypes:       []string{"pranzo"},
        MinPortion:      125,
        MaxPortion:      125,
        Category:        "protein",
    },
    "lattuga": {
        Name:            "Lattuga",
        StandardPortion: 80,
        Unit:            "g",
        CaloriesPer100g: 15,
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      80,
        MaxPortion:      160,
        Category:        "vegetable",
    },
    "pomodori": {
        Name:            "Pomodori da insalata",
        StandardPortion: 200,
        Unit:            "g",
        CaloriesPer100g: 18,
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      200,
        MaxPortion:      200,
        Category:        "vegetable",
    },
    "tacchino": {
        Name:            "Tacchino petto",
        StandardPortion: 250,
        Unit:            "g",
        CaloriesPer100g: 104,
        MealTypes:       []string{"cena"},
        MinPortion:      250,
        MaxPortion:      250,
        Category:        "protein",
    },
    "zucchine": {
        Name:            "Zucchine",
        StandardPortion: 300,
        Unit:            "g",
        CaloriesPer100g: 17,
        MealTypes:       []string{"pranzo", "cena"},
        MinPortion:      200,
        MaxPortion:      300,
        Category:        "vegetable",
    },
}

var mealTemplates = map[string]MealTemplate{
	"colazione": {
		MainOptions: []string{"panbauletto", "pane_integrale"},
		SideOptions: []string{"prosciutto_cotto", "salmone_affumicato", "frutta_fresca"},
		RequiredTypes: []string{"carb", "protein"},
		StandardStructure: []string{
			"Panbauletto Integrale - 48g",
			"Caffè - 30g",
			"Proteine (prosciutto/salmone) - 50g",
			"Frutta fresca - 150g",
		},
	},
	"spuntino": {
		MainOptions: []string{"frutta_fresca", "yogurt_greco"},
		SideOptions: []string{"frutta_secca"},
		StandardStructure: []string{
			"Frutta fresca - 150g",
			"Frutta secca - 20g",
		},
	},
	"pranzo": {
		MainOptions: []string{"pasta", "riso", "pane_integrale"},
		SideOptions: []string{"verdure", "legumi"},
		RequiredTypes: []string{"carb", "protein", "vegetable"},
		StandardStructure: []string{
			"Carboidrati (pasta/riso) - 80g",
			"Proteine (tonno/pollo/legumi) - 150g",
			"Verdure - 200g",
		},
	},
	"merenda": {
		MainOptions: []string{"yogurt_greco", "frutta_fresca"},
		SideOptions: []string{"frutta_secca"},
		StandardStructure: []string{
			"Yogurt greco - 150g",
			"Frutta secca - 20g",
		},
	},
	"cena": {
		MainOptions: []string{"pesce", "carne", "uova"},
		SideOptions: []string{"verdure", "patate"},
		RequiredTypes: []string{"protein", "vegetable"},
		StandardStructure: []string{
			"Proteine (pesce/carne) - 200g",
			"Verdure - 300g",
			"Pane integrale - 120g",
		},
	},
}

func calculateCalories(quantity float64, caloriesPer100g float64) float64 {
    return (quantity * caloriesPer100g) / 100
}

func generateMealWithUserIngredients(mealType string, userIngredients []string, targetCalories float64) Meal {

	items := []Food{}
    totalCalories := float64(0)


	
	// Prima, cerca di utilizzare gli ingredienti dell'utente se appropriati per questo pasto
	usedUserIngredients = make(map[string]bool)
	template = mealTemplates[mealType]

	// 1. Cerca di inserire gli ingredienti dell'utente nei posti più appropriati
	for _, ing := range userIngredients {
		rule, exists := foodRules[ing]
		if exists {
			// Verifica se l'ingrediente è appropriato per questo tipo di pasto
			isAppropriate := false
			for _, allowedMeal := range rule.MealTypes {
				if allowedMeal == mealType {
					isAppropriate = true
					break
				}
			}

			if isAppropriate {
				calories := calculateCalories(rule.StandardPortion, rule.CaloriesPer100g)
				if totalCalories + calories <= targetCalories {
					items = append(items, Food{
						Name:     rule.Name,
						Quantity: rule.StandardPortion,
						Unit:     rule.Unit,
						Calories: math.Round(calories),
					})
					totalCalories += calories
					usedUserIngredients[ing] = true
				}
			}
		}
	}

// 2. Completa il pasto con alimenti standard dalla template
for _, option := range template.MainOptions {
    if totalCalories >= targetCalories {
        break
    }

    if !usedUserIngredients[option] {
        if rule, exists := foodRules[option]; exists {
            calories := calculateCalories(rule.StandardPortion, rule.CaloriesPer100g)
            if totalCalories + calories <= targetCalories {
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

// Aggiungi alimenti dai SideOptions se necessario
for _, option := range template.SideOptions {
    if totalCalories >= targetCalories {
        break
    }

    if !usedUserIngredients[option] {
        if rule, exists := foodRules[option]; exists {
            calories := calculateCalories(rule.StandardPortion, rule.CaloriesPer100g)
            if totalCalories + calories <= targetCalories {
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

	// 3. Aggiungi sempre verdure a pranzo e cena se non ci sono
	if (mealType == "pranzo" || mealType == "cena") && totalCalories < targetCalories {
		vegetables := Food{
			Name:     "Verdure miste",
			Quantity: 200,
			Unit:     "g",
			Calories: 50,
		}
		items = append(items, vegetables)
		totalCalories += vegetables.Calories
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

	 // Nuovo endpoint per ottenere le alternative di un ingrediente
	 r.GET("/api/ingredients/:name/alternatives", func(c *gin.Context) {
        name := c.Param("name")
        if rule, exists := foodRules[name]; exists {
            c.JSON(http.StatusOK, gin.H{
                "ingredient": rule.Name,
            })
        } else {
            c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
        }
    })

    // Modifichiamo l'endpoint degli ingredienti per includere info sulle alternative
    r.GET("/api/ingredients", func(c *gin.Context) {
        result := make([]gin.H, 0)
        for k, rule := range foodRules {
            result = append(result, gin.H{
                "id":           k,
                "name":         rule.Name,
                "description":  rule.Description,
            })
        }
        c.JSON(http.StatusOK, result)
    })

	r.GET("/api/ingredients", func(c *gin.Context) {
		ingredients := make([]string, 0, len(foodRules))
		for k := range foodRules {
			ingredients = append(ingredients, k)
		}
		c.JSON(http.StatusOK, ingredients)
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