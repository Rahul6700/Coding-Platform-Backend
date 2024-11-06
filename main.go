package main

import (
  "final/auth"
  "final/compiler"
  "final/questions"
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

var db *gorm.DB

func main(){

  db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
  if err != nil {
     panic("failed to connect to the database")
  }


  db.AutoMigrate(&compiler.Dbstruct{}, &auth.User{})

  auth.SetDB(db)

  r := gin.Default()

  r.POST("/signup",auth.SignUp)
  r.POST("/signin",auth.SignIn)
  r.GET("/fetch",auth.Fetch)

  r.POST("/run",compiler.Compile)

  r.POST("/question",questions.CreateQuestion)
  r.GET("/question",questions.FetchQuestions)

  r.Run()

}
