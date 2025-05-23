// In backend/data/languages.go
package data

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strathlearn/backend/models" // Assuming this path is correct based on your module

	"github.com/gin-gonic/gin"
)

func GetLanguages(c *gin.Context) {
	// Path relative to the CWD (which is likely the 'backend' directory)
	filePath := "data/languages.json"

	// Optional: Log the CWD and the path being tried for debugging
	wd, _ := os.Getwd()
	log.Printf("Current working directory: %s", wd)
	log.Printf("Attempting to read file from: %s", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
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
