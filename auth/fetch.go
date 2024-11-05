package auth

import (
  //"gorm.io/gorm"
  //"gorm.io/driver/sqlite"
  "github.com/gin-gonic/gin"
  //"final/auth"
  "fmt"
)

func Fetch(c *gin.Context){
 
  db = GetDB() //get the db from db.go file

  var users []User
  result := db.Find(&users)

  if result.Error != nil {
    fmt.Println(result.Error)
    c.String(500, "failed to fetch data in fetch")
      return
  }

  if len(users) == 0{
    c.String(400,"no users found")
  }else{
  c.JSON(200,gin.H{"users":users})
  }
}
