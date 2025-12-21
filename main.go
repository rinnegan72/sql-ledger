package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type UpdateData struct {
	Amount int `json:"amount"`
}

func main() {
	env_err := godotenv.Load()
	if env_err != nil {
		fmt.Println("There is not .env file, make sure it exists")
	}
	fmt.Println("You get as much as you ask!")
	sql_db_host := os.Getenv("DATABASE_HOST")
	sql_db_user := os.Getenv("DATABASE_USER")
	sql_db_password := os.Getenv("DATABASE_PASSWORD")
	sql_db_port := os.Getenv("DATABASE_PORT")
	sql_db_name := os.Getenv("DATABASE_NAME")
	sql_ssl_mode := os.Getenv("DATABASE_SSL_MODE")
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=%s", sql_db_user, sql_db_password, sql_db_host, sql_db_port, sql_db_name, sql_ssl_mode)
	// dsn form for mysql: username:password@protocol(address)/dbname?param=value
	var err error
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.POST("/user/:id/:action", updateByID)
	router.Run("localhost:8080")
}
func updateByID(c *gin.Context) {
	var user UpdateData
	id := c.Param("id")
	action := c.Param("action")
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if action == "deposit" || action == "withdraw" {
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong with connecting to the DB"})
			fmt.Println(err)
			return
		}
		fmt.Println("performing ", action, " transaction for userid ", id, " amount ", user.Amount)
		var account_balance int
		err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ? FOR UPDATE", id).Scan(&account_balance)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to start transaction"})
			fmt.Println(err)
			return
		}
		transaction_amount := user.Amount
		if action == "withdraw" {
			if transaction_amount > account_balance {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
				fmt.Println(err)
				return
			}
			transaction_amount = transaction_amount * -1
		}
		account_balance = account_balance + transaction_amount
		_, err = tx.Exec("UPDATE accounts SET balance = ? WHERE id = ?", account_balance, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Update failed"})
			fmt.Println(err)
			return
		}
		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Commit failed"})
			fmt.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": "true", "message": "Transaction successful"})
	} else {
		fmt.Println("Invalid Action")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Action"})
	}
}
