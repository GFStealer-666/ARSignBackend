package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:khmksAKRJREjvPNBXUTYzbREXokRNkXE@tcp(interchange.proxy.rlwy.net:18891)/railway?charset=utf8mb4&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB: " + err.Error())
	}

	type QuizAnswerDetail struct {
		gorm.Model
		QuestionText     string `json:"questionText"`
		SelectedAnswer   string `json:"selectedAnswer"`
		CorrectAnswer    string `json:"correctAnswer"`
		TimeToAnswer     string `json:"timeToAnswer"` // <- change from float64 to string
		WasCorrect       string `json:"wasCorrect"`   // <- also change if you're passing "ถูกต้อง"
		QuizPlayerDataID uint
	}

	type PlayerQuizData struct {
		gorm.Model
		ProjectID             string             `json:"projectId"`
		APIKey                string             `json:"apiKey"`
		AppVersion            string             `json:"appversion"`
		UniqueID              string             `json:"uniqueId"`
		UserName              string             `json:"userName"`
		Gender                string             `json:"gender"`
		Age                   int                `json:"age"`
		Topic                 string             `json:"topic"`
		Difficulty            string             `json:"difficulty"`
		UserType              string             `json:"userType"`
		Score                 int                `json:"score"`
		MaxScore              int                `json:"maxScore"`
		Accuracy              string             `json:"accuracy"`
		Correct               int                `json:"correct"`
		Wrong                 int                `json:"wrong"`
		Stars                 int                `json:"stars"`
		TimeTaken             string             `json:"timeTaken"`
		Timestamp             string             `json:"timestamp"`
		IsCollectAdvancedData bool               `json:"isCollectAdvancedData"`
		AnswerDetails         []QuizAnswerDetail `json:"answerDetails" gorm:"foreignKey:QuizPlayerDataID"`
	}

	db := database
	db.AutoMigrate(&PlayerQuizData{}, &QuizAnswerDetail{})

	r := gin.Default()
	r.GET("/players/:id", func(c *gin.Context) {
		var p PlayerQuizData
		if err := db.First(&p, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		c.JSON(200, p)
	})

	r.GET("/players", func(c *gin.Context) {
		var players []PlayerQuizData
		if err := db.Find(&players).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, players)
	})

	r.POST("/submit-quiz", func(c *gin.Context) {
		var p PlayerQuizData
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&p)
		c.JSON(201, p)
	})

	r.DELETE("/delete/players/:id", func(c *gin.Context) {
		var p PlayerQuizData
		if err := db.First(&p, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		db.Delete(&p)
		c.JSON(200, gin.H{"message": "Player deleted"})
	})

	r.Run()
}
