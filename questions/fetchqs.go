package questions

import (
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

func FetchQuestions(c *gin.Context){

db, err := gorm.Open(sqlite.Open("data.db"),&gorm.Config{})
if err != nil {
  c.String(500, "error opening DB in questions.go")
}
db.AutoMigrate(&Question{})

  var qsn []Question
  result := db.Find(&qsn)
  if result.Error != nil {
    c.String(500, "error fetching qs's from db")
  }

  c.JSON(200, gin.H{"questions" : qsn})

}
