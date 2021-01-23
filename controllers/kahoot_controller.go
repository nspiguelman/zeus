package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nspiguelman/zeus/data"
	"github.com/nspiguelman/zeus/domain"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

type KahootController struct {
	KahootGames *domain.KahootGame
	kahootRepository *data.KahootRepository
	userRepository *data.UserRepository
	questionRepository *data.QuestionRepository
	answerRepository *data.AnswerRepository
}

func NewKahootController() KahootController {
	dbData := data.New()
	sqlDB, _ := dbData.DB.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Panic(err.Error())
	}

	return KahootController{
		KahootGames: &domain.KahootGame{},
		kahootRepository: &data.KahootRepository{
			Data: dbData,
		},
		userRepository: &data.UserRepository{
			Data: dbData,
		},
		questionRepository: &data.QuestionRepository{
			Data: dbData,
		},
		answerRepository: &data.AnswerRepository{
			Data: dbData,
		},
	}
}

func (kc *KahootController) CreateKahoot() func(c *gin.Context) {
	return func (c *gin.Context) {
		var kahootInput domain.KahootInput
		if err := c.ShouldBindJSON(&kahootInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
			return
		}

		if err := validator.Validate(kahootInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
			return
		}

		kahoot := domain.NewKahoot(kahootInput)
		if err := kc.kahootRepository.Create(kahoot); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
			return
		}

		c.JSON(http.StatusCreated, kahoot)
	}
}

func (kc *KahootController) CreateQuestion() gin.HandlerFunc {
	return func (c *gin.Context) {
		pin := c.Param("pin")
		log.Print(pin)
		kahoot, err := kc.kahootRepository.GetByPin(pin);
		if kahoot == nil && err == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room " + pin + " not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var questionInput domain.QuestionInput
		if err := c.ShouldBindJSON(&questionInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing input: " + err.Error()})
			return
		}

		if err := validator.Validate(questionInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error validating input: " + err.Error()})
			return
		}

		question := domain.NewQuestion(kahoot.ID, questionInput)
		if err := kc.questionRepository.Create(question); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": "Error saving question: " + err.Error() })
			return
		}

		answerNo := 0
		trueFound := false
		answers := make([]domain.Answer, 0)

		for _ ,answerInput := range questionInput.Answers {
			if trueFound && answerInput.IsTrue {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Only one answer can be true"})
				return
			}

			answer := domain.NewAnswer(question.ID, answerNo, answerInput)
			if err := kc.answerRepository.Create(answer); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{ "error": "Error creating answer: " + err.Error() })
				return
			}
			answers = append(answers, *answer)

			answerNo++
			trueFound = answerInput.IsTrue
		}

		c.JSON(http.StatusCreated, gin.H{
			"question": question,
			"answers": answers,
		})
 	}
}