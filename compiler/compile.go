package compiler

import (
  "github.com/gin-gonic/gin"
  "os"
  "os/exec"
  "strings"
  "bytes"
  "errors"
  "fmt"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
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

  var details Inp
   if err := c.BindJSON(&details); err != nil {
    c.String(400, "Invalid input: %v", err)
    return
}
  //opening db connection
  db, err := gorm.Open(sqlite.Open("data.db"),&gorm.Config{})
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

  if details.Language == "python"{

    cmd := exec.Command("python3","temp.txt")
    cmd.Stdin = strings.NewReader(details.Input)
    //fmt.Println("running code with",details.Input)
    //cmd.Run()
    cmd.Stdout = &buff
    cmd.Stderr = &buff
    err = cmd.Run()

    if err != nil{
      c.String(400,"error running code")
    }
    fmt.Println("Buffer content:", buff.String())
    c.String(200, buff.String())

  }else if details.Language == "javascript"{

    cmd := exec.Command("node","temp.txt")
    cmd.Stdin = strings.NewReader(details.Input)
    //fmt.Println("running code with",details.Input)
    //cmd.Run()
    cmd.Stdout = &buff
    cmd.Stderr = &buff
    err = cmd.Run()

    if err != nil{
      c.String(400,"error running code")
    }
    fmt.Println("Buffer content:", buff.String())
    c.String(200, buff.String())

  }else if details.Language == "ruby"{

    cmd := exec.Command("ruby","temp.txt")
    cmd.Stdin = strings.NewReader(details.Input)
    //fmt.Println("running code with",details.Input)
    //cmd.Run()
    cmd.Stdout = &buff
    cmd.Stderr = &buff
    err = cmd.Run()

    if err != nil{
      c.String(400,"error running code")
    }
    fmt.Println("Buffer content:", buff.String())
    c.String(200, buff.String())

  }else if details.Language == "php" {

    cmd := exec.Command("php","temp.txt")
    cmd.Stdin = strings.NewReader(details.Input)
    //fmt.Println("running code with",details.Input)
    //cmd.Run()
    cmd.Stdout = &buff
    cmd.Stderr = &buff
    err = cmd.Run()

    if err != nil{
      c.String(400,"error running code")
    }
    fmt.Println("Buffer content:", buff.String())
    c.String(200, buff.String())

  }else if details.Language == "perl" {

    cmd := exec.Command("perl","temp.txt")
    cmd.Stdin = strings.NewReader(details.Input)
    //fmt.Println("running code with",details.Input)
    //cmd.Run()
    cmd.Stdout = &buff
    cmd.Stderr = &buff
    err = cmd.Run()

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
