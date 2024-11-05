package auth

import (
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

type SignUpStruct struct {
  Username string  `json:"username"`
  Password string  `json:"password"`
}

func SignIn(c *gin.Context){

  var dets SignUpStruct
  c.BindJSON(&dets)

  UserDets := SignUpStruct{
    Username : dets.Username,
    Password : dets.Password,
  }

  db, err := gorm.Open(sqlite.Open("code.db"),gorm.Config{})
  if err != nil {
    c.String(500,"error while opening db in signup")
  }

  db.AutoMigrate(&SignUpStruct)

  UserDets := SignUpStruct {
    Username : dets.Username,
    Password : dets.Password,
  }
  
  result := db.Where("username = ? AND password = ?", dets.Username, dets.Password)

  if result.Error != nil {
    if result.Error == gorm.ErrRecordNotFound {
      err := db.Create(&UserDets)
      if err != nil{
        c.String(500, "error while writing new credentials to db in signup")
      }
      c.String(200,"account created Successfully")
    }else {
      c.String(400, "User already exists")
    }
  }
}
