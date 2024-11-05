package main

import (
  "final/auth"
  "final/compiler"
  "github.com/gin-gonic/gin"
)

func main(){
  
  r := gin.Default()

  r.POST("/signup",auth.SignUp)
  r.POST("/signin",auth.SignIn)

  r.POST("/run",compiler.Compile)

  r.Run()

}
