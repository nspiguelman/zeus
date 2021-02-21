package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nspiguelman/zeus/data"
	"github.com/nspiguelman/zeus/domain"
	"github.com/nspiguelman/zeus/services"
	"gopkg.in/olahol/melody.v1"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"strconv"
)

type KahootController struct {
	rm         *data.RepositoryManager
	KahootGame *services.KahootGame
}

func NewKahootController() KahootController {
	dbData := data.New()
	sqlDB, _ := dbData.DB.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Panic(err.Error())
	}
	repositoryManager := data.NewRepositoryManager(dbData)

	return KahootController{
		rm:         repositoryManager,
		KahootGame: services.NewKahootGame(repositoryManager),
	}
}

func (kc *KahootController) Ping() gin.HandlerFunc {
	return func (c *gin.Context) {
		c.JSON(200, gin.H{ "message": "pong" })
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

		kahoot, err := kc.KahootGame.CreateKahoot(kahootInput)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
			return
		}
		c.JSON(http.StatusCreated, kahoot)
	}
}

func (kc *KahootController) CreateQuestion() gin.HandlerFunc {
	return func (c *gin.Context) {
		var questionInput domain.QuestionInput
		if err := c.ShouldBindJSON(&questionInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing input: " + err.Error()})
			return
		}
		if err := validator.Validate(questionInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error validating input: " + err.Error()})
			return
		}

		pin := c.Param("pin")
		question, statusCode, err := kc.KahootGame.CreateQuestion(questionInput, pin)
		if statusCode == http.StatusCreated {
			c.JSON(statusCode, question)
		} else {
			c.JSON(statusCode, gin.H{"error": err})
		}
 	}
}

func (kc *KahootController) CreateAnswer () gin.HandlerFunc {
	return func (c *gin.Context) {
		fmt.Println("Entro a createAnswer")
		var answerInput []domain.AnswerInput
		if err := c.ShouldBindJSON(&answerInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing input: " + err.Error()})
			return
		}

		questionId, err := strconv.Atoi(c.Param("questionId"))
		fmt.Println(questionId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
			return
		}

		answers, statusCode, errMessage := kc.KahootGame.CreateAnswers(answerInput, questionId)
		if statusCode == http.StatusCreated {
			c.JSON(statusCode, answers)
		} else {
			c.JSON(statusCode, gin.H{ "error": errMessage })
		}
	}
}

// TODO: ordenar las llamadas. Las llamadas al rm no deben estar aca.
// TODO: Si hay logica tiene que estar desarrollada mediante kahootGames,
func (kc *KahootController) Login() gin.HandlerFunc {
	return func (c *gin.Context) {

		name := c.Param("name")
		var token = kc.KahootGame.GenerateToken(name)

		pin := c.Param("pin")
		kahoot, err := kc.rm.KahootRepository.GetByPin(pin)
		if kahoot == nil && err == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room " + pin + " not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		KahootUser := domain.NewUser(name,token,kahoot.ID)

		// TODO: arrancar transaction hasta crear en user Postgres y el key value en redis
		if err := kc.rm.UserRepository.Create(KahootUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": "Error saving user: " + err.Error() })
			return
		}
		if err := kc.KahootGame.InitScore(token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": "Error saving user score: " + err.Error() })
			return
		}
		// TODO: terminar transaction
		c.JSON(200, gin.H{
			"token": token,
		})
	}
}

func (kc *KahootController) HandShake(socket *melody.Melody) gin.HandlerFunc {
	return func (c *gin.Context) {
		socket.HandleRequest(c.Writer, c.Request)
	}
}

func (kc *KahootController) HandleMessage(socket *melody.Melody) {
	socket.HandleMessage(func(x *melody.Session, msg []byte) {
		// pin := x.Request.Header.Get("pin")
		token := x.Request.Header.Get("token")

		answer := domain.AnswerMessage{}

		err := json.Unmarshal([]byte(msg), &answer)
		if err != nil {
			panic(err.Error())
		}
		answer.Token = token

		kc.KahootGame.Answer(answer)
	})
}


func (kc *KahootController) SendQuestion(socket *melody.Melody) gin.HandlerFunc {
	return func(c *gin.Context){
		pin := c.Param("pin")
		if !kc.KahootGame.IsStarted {
			if err := kc.KahootGame.Start(pin); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			if !kc.KahootGame.IsScoreSent {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send next question. Score must be sent."})
				return
			}

			if err := kc.KahootGame.NextQuestion(pin); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		kc.KahootGame.BroadCastQuestion(socket, domain.QuestionMessage{
			QuestionId: kc.KahootGame.CurrentQuestion,
			AnswerIds: kc.KahootGame.GetCurrentAnswerIds(),
			TypeMessage: "question",
		})
	}
}