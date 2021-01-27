package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nspiguelman/zeus/data"
	"github.com/nspiguelman/zeus/domain"
	"github.com/nspiguelman/zeus/services"
	"gopkg.in/olahol/melody.v1"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

type KahootController struct {
	rm *data.RepositoryManager
	KahootGames *services.KahootGame
}

func NewKahootController() KahootController {
	dbData := data.New()
	sqlDB, _ := dbData.DB.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Panic(err.Error())
	}
	repositoryManager := data.NewRepositoryManager(dbData)

	return KahootController{
		rm:          repositoryManager,
		KahootGames: services.NewKahootGame(repositoryManager),
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

		kahoot := domain.NewKahoot(kahootInput)
		if err := kc.rm.KahootRepository.Create(kahoot); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
			return
		}
		c.JSON(http.StatusCreated, kahoot)
	}
}

func (kc *KahootController) CreateQuestion() gin.HandlerFunc {
	return func (c *gin.Context) {
		pin := c.Param("pin")

		kahoot, err := kc.rm.KahootRepository.GetByPin(pin);
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
		if err := kc.rm.QuestionRepository.Create(question); err != nil {
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
			if err := kc.rm.AnswerRepository.Create(answer); err != nil {
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

func (kc *KahootController) Login() gin.HandlerFunc {
	return func (c *gin.Context) {
		// TODO: persistir en la db
		name := c.Param("name")
		var token = kc.KahootGames.GenerateToken(name)
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
	// TODO: 1. Hacer un channel donde reciba todas las respuesta del cliente
	// TODO: 2. Hacer una goroutine que este todo el tiempo escuchando al chanel para procesar la respuesta y guardar en bdd (siempre y cuando el canal este abierto , lo cual es determinado por el timeOut de la maquina de estado
	// TODO: 3. Cerrar el channel cuando haya timeout.
	// TODO: 4. Analizar si mover el timeout de la maquina de estados y dejarlo como una función timer
	socket.HandleMessage(func(x *melody.Session, msg []byte) {
		// pin := x.Request.Header.Get("pin")
		kc.KahootGames.ArrivalOrder += 1
		token := x.Request.Header.Get("token")

		answer := domain.AnswerMessage{}
		answer.Token = token

		err := json.Unmarshal([]byte(msg), &answer)
		if err != nil {
			panic(err.Error())
		}
		answer.Token = token

		if kc.KahootGames.IsTimeout {
			answer.IsTimeout = true
		}
		kc.KahootGames.ProcessAnswer(answer)
	})
}


// TODO: Analizar si la goroutine muere antes. Porque tenemos la duda si la función termina e interrumpe la goroutine
func (kc *KahootController) SendQuestion(socket *melody.Melody) gin.HandlerFunc {
	return func(c *gin.Context){
		pin := c.Param("pin")
		if !kc.KahootGames.IsStarted {
			if err := kc.KahootGames.Start(pin); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			if !kc.KahootGames.IsScoreSent {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send next question. Score must be sent."})
				return
			}

			if err := kc.KahootGames.NextQuestion(pin); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// setear timeout en false
		kc.KahootGames.BroadCastQuestion(socket, domain.QuestionMessage{
			QuestionId: kc.KahootGames.CurrentQuestion,
			AnswerIds: kc.KahootGames.GetCurrentAnswerIds(),
			TypeMessage: "question",
		})
	}
}