package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var tablePrefix = os.Getenv("TABLE_PREFIX")

type PlayerQuizData struct {
	gorm.Model
	AppVersion            string             `json:"appversion"`
	UniqueID              string             `json:"uniqueId"` // Unique ID for devices
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

func (PlayerQuizData) TableName() string {
	return tablePrefix + "_player_data"
}

type QuizAnswerDetail struct {
	gorm.Model
	QuestionText     string `json:"questionText"`
	SelectedAnswer   string `json:"selectedAnswer"`
	CorrectAnswer    string `json:"correctAnswer"`
	TimeToAnswer     string `json:"timeToAnswer"`
	WasCorrect       string `json:"wasCorrect"`
	QuizPlayerDataID uint
}

func (QuizAnswerDetail) TableName() string {
	return tablePrefix + "_answer_details"
}

func main() {

	var dsn = os.Getenv("DATABASE_URL")

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB: " + err.Error())
	}

	db := database
	db.AutoMigrate(&PlayerQuizData{}, &QuizAnswerDetail{})

	r := gin.Default()
	r.GET("/players/:id", func(c *gin.Context) {
		var p PlayerQuizData
		if err := db.Preload("AnswerDetails").First(&p, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		c.JSON(200, p)
	})

	r.GET("/leaderboard", func(c *gin.Context) {
		var players []PlayerQuizData
		if err := db.Preload("AnswerDetails").Find(&players).Error; err != nil {
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local
	}
	r.Run(":" + port)
}
