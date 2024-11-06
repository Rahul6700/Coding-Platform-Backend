package questions

import (
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  "errors"
)

type Question struct {
    gorm.Model
    Question  string    `json:"question"`
    Level     string       `json:"level"`
    //TestCases []TestCase `json:"testcases"`
}

// gorm.Model makes sure that there is a unique ID for each element written to db
//
// type TestCase struct {
//     gorm.Model
//     Input      string `json:"input"`
//     Expected   string `json:"expected"`
// }

func CreateQuestion(c *gin.Context) {

db, err := gorm.Open(sqlite.Open("data.db"),&gorm.Config{})
if err != nil {
  c.String(500, "error opening DB in questions.go")
}

db.AutoMigrate(&Question{})

  var Quest Question
  c.BindJSON(&Quest)

  NewReq := Question {
    Question : Quest.Question,
    Level : Quest.Level,
    //TestCases : Quest.TestCases,
  }

result := db.First(&Question{}, " Question = ? AND Level = ? ", NewReq.Question, NewReq.Level)

if result.RowsAffected > 0 {
  c.String(400, "Question already exists in DB")  
  return
}
if errors.Is(result.Error, gorm.ErrRecordNotFound) {
  db.Create(&NewReq) 
}

}

