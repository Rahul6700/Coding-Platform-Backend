package main

import (
	"final/auth"
	"final/compiler"
	"final/questions"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {

	var err error
	db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected successfully")

	err = db.AutoMigrate(&compiler.Dbstruct{}, &auth.User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("Database migration complete")
	auth.SetDB(db)

	r := gin.Default()
	fmt.Println("Router initialized")

	r.POST("/signup", auth.SignUp)
	r.POST("/signin", auth.SignIn)
	r.GET("/fetch", auth.Fetch)
	r.POST("/run", compiler.Compile)
	r.POST("/question", questions.CreateQuestion)
	r.GET("/question", questions.FetchQuestions)

	port := ":8080"
	fmt.Println("Server listening on port", port)
	err = r.Run(port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

