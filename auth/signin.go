package auth

import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  "github.com/gin-gonic/gin"
)

type SignInStruct struct {
  Username string  `json:"username"`
  Password string  `json:"password"`
}

func SignIn(c *gin.Context){

    var dets SignInStruct;
    c.BindJSON(&dets)
    
    db, err := gorm.Open(sqlite.Open("code.db"),&gorm.Config{})
    if err != nil {
      c.JSON(500,"failed to open DB connection in Singin")
    }

    db.AutoMigrate(&SignInStruct)

    UserDets := SignInStruct{
      username : dets.Username,
      password : dets.Password,
    }

    result := db.Where("username = ? AND password = ?", dets.Username, dets.Password)

    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            c.JSON(400, gin.H{"Error": "Invalid username or password"})
        } else {
            c.JSON(500, gin.H{"error": "Error checking database in signin"})
        }
        return
    }else{
        c.JSON(200, gin.H{"signin Sucessful": dets.Username})
    }
        
  }
