package compiler

import (
  "context"
  "time"
  "github.com/gin-gonic/gin"
  "os"
  "os/exec"
  "strings"
  "bytes"
  "errors"
  "fmt"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  "final/key"
  "github.com/golang-jwt/jwt/v5"
)

type Inp struct{
  //gorm.Model
  Language string `json:"language"`
  Code string  `json:"code"`
  Input string  `json:"input"`
}

type Dbstruct struct {
  gorm.Model
  Language string //`json:"language"`
  Code string  //`json:"code"`
  Input string  //`json:"input"`
  Output string
}

//supported languages for now: python, javascript, php, ruby, perl
// while sending requests use this format for the "language":
// "python", "javascript", "php", "ruby", "perl"

func Compile(c *gin.Context){

  //token validation happens first
  tokenstring := c.GetHeader("Authorization") //token is fetched from header
  if tokenstring == "" {
    c.String(400, "token missing")
    return
  }

  //token signing using the secret key
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte(key.SECRET_KEY), nil 
	})

  //if validation fails, the function execution terminates
	if err != nil || !token.Valid {
		c.String(400, "Token validation failed")
		return
	}

	//if token passes all these, the function runs
  var details Inp
   if err := c.BindJSON(&details); err != nil {
    c.String(400, "Invalid input: %v", err)
    return
}
  //opening db connection
  db, err := gorm.Open(sqlite.Open("code.db"),&gorm.Config{})
  if err != nil{
    panic("initial connection to db unsuccessful")
  }

  db.AutoMigrate(&Dbstruct{})

  //fmt.Println("the struct is",details) 
  file, err := os.Create("temp.txt")
  if err != nil{
    c.String(400, "error in creating temp file")
  }

  defer file.Close()

 /*  file.Write([]byte(details.code)) */
  fmt.Println("the code is",details.Code)
  n, err := file.WriteString(details.Code)
  if err != nil {
    c.String(400, "error while writing code to file %d",n)
  }

  // checking for memorisation
result := db.First(&Dbstruct{}, " Language = ? AND Code = ? AND Input = ?", details.Language, details.Code, details.Input)
var buff bytes.Buffer

//if no memorisation found
if errors.Is(result.Error, gorm.ErrRecordNotFound) {

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  if details.Language == "python"{

    cmd := exec.CommandContext(ctx,"python3","temp.txt")
    cmd.Stdin = strings.NewReader(details.Input)
    //fmt.Println("running code with",details.Input)
    //cmd.Run()
    cmd.Stdout = &buff
    cmd.Stderr = &buff
    err = cmd.Run()

  if ctx.Err() == context.DeadlineExceeded { 
    c.String(400, "Code execution timed out")
    return
  }

    if err != nil{
      c.String(400,"error running code")
    }
    fmt.Println("Buffer content:", buff.String())
    c.String(200, buff.String())

  }else if details.Language == "javascript"{

    cmd := exec.CommandContext(ctx,"node","temp.txt")
    cmd.Stdin = strings.NewReader(details.Input)
    //fmt.Println("running code with",details.Input)
    //cmd.Run()
    cmd.Stdout = &buff
    cmd.Stderr = &buff
    err = cmd.Run()

  if ctx.Err() == context.DeadlineExceeded { 
    c.String(400, "Code execution timed out")
    return
  }
    if err != nil{
      c.String(400,"error running code")
    }
    fmt.Println("Buffer content:", buff.String())
    c.String(200, buff.String())

  }else if details.Language == "ruby"{

    cmd := exec.CommandContext(ctx,"ruby","temp.txt")
    cmd.Stdin = strings.NewReader(details.Input)
    //fmt.Println("running code with",details.Input)
    //cmd.Run()
    cmd.Stdout = &buff
    cmd.Stderr = &buff
    err = cmd.Run()

  if ctx.Err() == context.DeadlineExceeded { 
    c.String(400, "Code execution timed out")
    return
  }

    if err != nil{
      c.String(400,"error running code")
    }
    fmt.Println("Buffer content:", buff.String())
    c.String(200, buff.String())

  }else if details.Language == "php" {

    cmd := exec.CommandContext(ctx,"php","temp.txt")
    cmd.Stdin = strings.NewReader(details.Input)
    //fmt.Println("running code with",details.Input)
    //cmd.Run()
    cmd.Stdout = &buff
    cmd.Stderr = &buff
    err = cmd.Run()

  if ctx.Err() == context.DeadlineExceeded { 
    c.String(400, "Code execution timed out")
    return
  }

    if err != nil{
      c.String(400,"error running code")
    }
    fmt.Println("Buffer content:", buff.String())
    c.String(200, buff.String())

  }else if details.Language == "perl" {

    cmd := exec.CommandContext(ctx,"perl","temp.txt")
    cmd.Stdin = strings.NewReader(details.Input)
    //fmt.Println("running code with",details.Input)
    //cmd.Run()
    cmd.Stdout = &buff
    cmd.Stderr = &buff
    err = cmd.Run()

  if ctx.Err() == context.DeadlineExceeded { 
    c.String(400, "Code execution timed out")
    return
  }

    if err != nil{
      c.String(400,"error running code")
    }
    fmt.Println("Buffer content:", buff.String())
    c.String(200, buff.String())

  }
  return

}
// fetching output from db using memorisation
if result.RowsAffected > 0 {
    var myInput Dbstruct
    db.First(&myInput, "Language = ? AND Code = ? AND Input = ?",details.Language, details.Code, details.Input)
    myOutput := myInput.Output
    c.String(200, myOutput)
    return
}
    
      myInput := Dbstruct{
      Language: details.Language,
      Code: details.Code,
      Input: details.Input,
      Output: buff.String(),
    }
    db.Create(&myInput)


 os.Remove("temp.txt")
}
