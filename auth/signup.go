package auth

import (
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
  //"gorm.io/driver/sqlite"
)

type User struct {
  Username string  `json:"username"`
  Password string  `json:"password"`
}

func SignUp(c *gin.Context){

  var dets User
  if err := c.BindJSON(&dets); err != nil {
    c.JSON(400, gin.H{"error": "Invalid input"})
    return
  }

  // db, err := gorm.Open(sqlite.Open("code.db"),&gorm.Config{})
  // if err != nil {
  //   c.String(500, "error while opening db in signup")
  //   return 
  // }
  //
  // db.AutoMigrate(&SignUpStruct{})
  //

  db := GetDB() //get access to the db from db.go

  var existingUser User
  result := db.Where("username = ?", dets.Username).First(&existingUser)  //this checks if the user alr exists

  if result.Error == nil { 
    c.JSON(400, gin.H{"error": "User already exists"})
    return
  } else if result.Error != gorm.ErrRecordNotFound { //is error is anything else except 'user not found'
    c.String(500, "error while checking for existing user")
    return
  }

  //'user not found' case, the content is written to db
  err := db.Create(&dets).Error 
  if err != nil {
    c.String(500, "error while writing new credentials to db in signup")
    return
  }

  c.String(200, "account created successfully")
}
