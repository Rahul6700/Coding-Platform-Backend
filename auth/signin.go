package auth

import (
  "gorm.io/gorm"
  //"gorm.io/driver/sqlite"
  "github.com/gin-gonic/gin"
  "fmt"
  "final/key"
  "time"
  "github.com/golang-jwt/jwt/v5"
)
//
// type SignInStruct struct {
//   Username string  `json:"username"`
//   Password string  `json:"password"`
// }

//var SECRET_KEY = "I am the secret key"

func SignIn(c *gin.Context){

    var dets User;
    if err := c.BindJSON(&dets); err != nil { 
        c.JSON(400, gin.H{"error": "Invalid input"})
        return
    }

    // db, err := gorm.Open(sqlite.Open("code.db"),&gorm.Config{})
    // if err != nil {
    //   c.JSON(500,"failed to open DB connection in Singin")
    //   return 
    // }
    //
    // db.AutoMigrate(&SignInStruct{})
    // 
    // UserDets := SignInStruct{
    //   username : dets.Username,
    //   password : dets.Password,
    // }

    db := GetDB() //get the db from the db.go file

    var Person User
    result := db.Where("username = ? AND password = ?", dets.Username, dets.Password).First(&Person)

    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            c.JSON(400, gin.H{"Error": "Invalid username or password"})
        } else {
            fmt.Println(result.Error)
            c.JSON(500, gin.H{"error": "Error checking database in signin"})
        }
        return
    }
    
    //upon successful signin a JWT is generated for that username
	ttl := time.Second * 300 //token's time to live is 5 mins
	claims := jwt.MapClaims{
		"username": Person.Username,
		"exp":      time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key.SECRET_KEY))
	if err != nil {
		c.String(500,"error signing token")
	}

    c.JSON(200, gin.H{"Signed in successfully, here is your JWT": tokenString})    
}
