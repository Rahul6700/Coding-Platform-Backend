package compiler

import (
  //"context"
  //"time"
  "github.com/gin-gonic/gin"
  "os"
  //"strings"
  //"errors"
  //"fmt"
  "log"
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

func Compile(c *gin.Context){
  // Token validation (existing code)
  tokenstring := c.GetHeader("Authorization")
  if tokenstring == "" {
    c.String(400, "token missing")
    return
  }

  token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
    return []byte(key.SECRET_KEY), nil 
  })

  if err != nil || !token.Valid {
    c.String(400, "Token validation failed")
    return
  }

  // Parse input details
  var details Inp
  if err := c.BindJSON(&details); err != nil {
    c.String(400, "Invalid input: %v", err)
    return
  }

  // Open database connection
  db, err := gorm.Open(sqlite.Open("code.db"), &gorm.Config{})
  if err != nil {
    log.Println("Database connection failed")
    c.String(500, "Database connection failed")
    return
  }

  db.AutoMigrate(&Dbstruct{})

  // Create temporary file
  tempFile, err := os.CreateTemp("", "usercode-*")
  if err != nil {
    log.Println("Failed to create temporary file")
    c.String(500, "Failed to create temporary file")
    return
  }
  defer os.Remove(tempFile.Name())
  defer tempFile.Close()

  // Write code to temporary file
  if _, err := tempFile.WriteString(details.Code); err != nil {
    log.Println("Failed to write code to file")
    c.String(500, "Failed to write code to file")
    return
  }

  // Check for memoization
  result := db.First(&Dbstruct{}, "Language = ? AND Code = ? AND Input = ?", details.Language, details.Code, details.Input)

  // Language mapping for Docker
  languageMap := map[string]string{
    "python": "py",
    "javascript": "js",
    "ruby": "rb",
    "php": "php",
    "perl": "pl",
  }

  dockerLang, ok := languageMap[details.Language]
  if !ok {
    log.Printf("Unsupported language: %s", details.Language)
    c.String(400, "Unsupported language")
    return
  }

  // If found in cache, return cached output
  if result.RowsAffected > 0 {
    var cachedResult Dbstruct
    db.First(&cachedResult, "Language = ? AND Code = ? AND Input = ?", details.Language, details.Code, details.Input)
    c.String(200, cachedResult.Output)
    return
  }

  // Run code in Docker sandbox
  log.Printf("Executing code in Docker. Language: %s, File: %s", dockerLang, tempFile.Name())
  output, errOutput := runDocker(tempFile, dockerLang, details.Input)

  // Check for execution errors
  if errOutput != "" {
    log.Printf("Execution error: %s", errOutput)
    c.String(400, "Execution error: %s", errOutput)
    return
  }

  // Save to database cache
  dbEntry := Dbstruct{
    Language: details.Language,
    Code:     details.Code,
    Input:    details.Input,
    Output:   output,
  }
  db.Create(&dbEntry)

  // Return output
  log.Printf("Code execution completed. Output: %s", output)
  c.String(200, output)
}
