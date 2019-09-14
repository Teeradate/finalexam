package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:status`
}

func authMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Error!!! Unauthorization"})
		c.Abort()
		return
	}

	c.Next()

}

var cusglobal []Customer

func postTodo(c *gin.Context) {
	url := os.Getenv("DATABASE_URL")
	fmt.Println(url)
	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	defer db.Close()
	t := Customer{}

	err1 := c.ShouldBindJSON(&t)
	if err1 != nil {
		log.Println(err1)
		c.JSON(http.StatusBadRequest, gin.H{"status": "JSON parsing on insertion error!!! " + err1.Error()})
		return
	}

	fmt.Println("Name ", t.Name)
	fmt.Println("Email ", t.Email)
	fmt.Println("Status ", t.Status)
	row := db.QueryRow("INSERT INTO customer (name, email, status) values ($1,$2,$3) RETURNING id", t.Name, t.Email, t.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Insertion error!!! " + err.Error()})
		return
	}

	t.ID = id

	c.JSON(http.StatusCreated, gin.H{
		"id":     t.ID,
		"name":   t.Name,
		"email":  t.Email,
		"status": t.Status,
	})
}

func getOneTodo(c *gin.Context) {
	url := os.Getenv("DATABASE_URL")
	fmt.Println(url)
	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	defer db.Close()

	id := c.Param("id")

	var cusLocal Customer

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customer WHERE id = $1")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL select error!!! " + err.Error()})
		return
	}

	idnum, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Convert id to num error!!!" + err.Error()})
		return
	}

	row := stmt.QueryRow(idnum)

	err = row.Scan(&cusLocal.ID, &cusLocal.Name, &cusLocal.Email, &cusLocal.Status)
	if err != nil {
		log.Println("Select id = " + id)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Select id %d error!!!: %s", idnum, err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     cusLocal.ID,
		"name":   cusLocal.Name,
		"email":  cusLocal.Email,
		"status": cusLocal.Status,
	})
}

func updateTodo(c *gin.Context) {
	url := os.Getenv("DATABASE_URL")
	fmt.Println(url)
	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	defer db.Close()

	id := c.Param("id")

	//Pre-select.
	var cusLocal Customer

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customer WHERE id = $1")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL select error!!! " + err.Error()})
		return
	}

	idnum, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Convert id to num error!!! " + err.Error()})
		return
	}

	row := stmt.QueryRow(idnum)

	err = row.Scan(&cusLocal.ID, &cusLocal.Name, &cusLocal.Email, &cusLocal.Status)
	if err != nil {
		log.Println("Select id = " + id)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Select id %d error!!!: %s", idnum, err.Error())})
		return
	}

	//Actual update.
	t := Customer{}

	err = c.ShouldBindJSON(&t)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "JSON parsing Error!!! " + err.Error()})
		return
	}

	stmt, err = db.Prepare("UPDATE customer SET name=$2,email=$3,status=$4 WHERE id=$1")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL for update error!!! " + err.Error()})
		return
	}

	_, err = stmt.Exec(id, t.Name, t.Email, t.Status)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Execute update error!!! " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     t.ID,
		"name":   t.Name,
		"email":  t.Email,
		"status": t.Status,
	})
}

func deleteTodo(c *gin.Context) {
	url := os.Getenv("DATABASE_URL")
	fmt.Println(url)
	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	defer db.Close()

	id := c.Param("id")
	//Pre-select.
	var cusLocal Customer

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customer WHERE id = $1")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL select error!!! " + err.Error()})
		return
	}

	idnum, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Convert id to num error!!!" + err.Error()})
		return
	}

	row := stmt.QueryRow(idnum)

	err = row.Scan(&cusLocal.ID, &cusLocal.Name, &cusLocal.Email, &cusLocal.Status)
	if err != nil {
		log.Println("Select id = " + id)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("Select id %d error!!!: %s", idnum, err.Error())})
		return
	}

	//Actual delete.
	stmt, err = db.Prepare("DELETE FROM customer WHERE id = $1")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Prepare SQL for delete row error!!! " + err.Error()})
		return
	}

	_, err = stmt.Exec(idnum)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Execute deletion error!!! " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "customer deleted",
	})
}

func getAllTodoHandler(c *gin.Context) {
	cusglobal = nil
	url := os.Getenv("DATABASE_URL")
	fmt.Println(url)
	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	defer db.Close()

	stmt, err := db.Prepare("select id, name, email, status from customer")

	if err != nil {
		log.Fatal("can't prepare query", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("can't query", err)
	}

	for rows.Next() {
		//var id int
		//var title, status string
		cust := Customer{}
		err := rows.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
		if err != nil {
			log.Fatal("can't scab", err)
		}
		cusglobal = append(cusglobal, cust)

	}

	c.JSON(http.StatusOK, cusglobal)

}

func main() {
	r := gin.Default()
	r.Use(authMiddleware)
	r.GET("/customers", getAllTodoHandler)
	r.POST("/customers", postTodo)
	r.GET("/customers/:id", getOneTodo)
	r.PUT("/customers/:id", updateTodo)
	r.DELETE("/customers/:id", deleteTodo)
	r.Run(":2019")
}
