package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS middleware personalizzato
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

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/ingredients", func(c *gin.Context) {
		ingredients := []string{
			"pesto",
			"pollo",
			"pasta",
			"riso",
			"insalata",
			"uova",
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

		// Piano di esempio
		plan := gin.H{
			"colazione": gin.H{
				"items": []gin.H{
					{
						"name":     "Pane integrale",
						"quantity": 50,
						"unit":     "g",
					},
					{
						"name":     "Caffè",
						"quantity": 1,
						"unit":     "tazza",
					},
				},
			},
			"pranzo": gin.H{
				"items": []gin.H{
					{
						"name":     "Pasta",
						"quantity": 80,
						"unit":     "g",
					},
					{
						"name":     "Pesto",
						"quantity": 30,
						"unit":     "g",
					},
				},
			},
			"cena": gin.H{
				"items": []gin.H{
					{
						"name":     "Pollo alla griglia",
						"quantity": 200,
						"unit":     "g",
					},
					{
						"name":     "Insalata mista",
						"quantity": 100,
						"unit":     "g",
					},
				},
			},
		}

		c.JSON(http.StatusOK, plan)
	})

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
