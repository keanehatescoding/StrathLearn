package data

import (
	"encoding/json"
	"net/http"
	"os"
	"strathlearn/backend/models"

	"github.com/gin-gonic/gin"
)

func GetLanguages(c *gin.Context) {

	file, err := os.ReadFile("data/languages.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read languages file"})
		return
	}

	var languages []models.Language
	if err := json.Unmarshal(file, &languages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse languages data"})
		return
	}

	activeLanguages := make([]models.Language, 0)
	for _, lang := range languages {
		if !lang.IsArchived {
			activeLanguages = append(activeLanguages, lang)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"languages": activeLanguages,
		"count":     len(activeLanguages),
	})
}
