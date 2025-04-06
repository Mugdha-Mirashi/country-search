package main

import (
	"country-search/constants"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(constants.ExternalServiceError)
	}
}
