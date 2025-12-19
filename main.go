package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateData struct {
	Amount int `json:"amount"`
}

func main() {
	fmt.Println("You get as much as you ask!")
	router := gin.Default()
	router.POST("/user/:id/:action", updateByID)
	router.Run("localhost:8080")
}
func updateByID(c *gin.Context) {
	var user UpdateData
	id := c.Param("id")
	action := c.Param("action")
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if action == "deposit" || action == "withdraw" {
		fmt.Println("performing ", action, " transaction for userid ", id, " amount ", user.Amount)
		c.JSON(http.StatusOK, gin.H{"sucess": "true", "message": "Transaction successful"})
	} else {
		fmt.Println("Invalid Action")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Action"})
	}
}
