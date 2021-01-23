package domain

type KahootGame struct {
	trivias      []Question
	host		User
	users		[]User
	pin 		string
	CurrentQuestion int // pregunta actual
	TotalQuestions int // total de preguntas
	IsTimeout bool // timeout de la pregunta
	IsStarted bool // ya empezo el kahoot, no se puede suscribir nadie
	IsScoreSent bool // la siguiente pregunta solo puede ser enviada cuando el score sea notificado
}

func NewKahootGame(host User) KahootGame {
	trivias := make([]Question, 0)
	users := make([]User, 0)
	pin := "123534" //TODO: Use random generator
	CurrentQuestion := 0
	TotalQuestions := 0
	IsTimeout := false
	IsStarted := false
	IsScoreSent := false
	return KahootGame{
		trivias,
		host,
		users,
		pin,
		CurrentQuestion,
		TotalQuestions,
		IsTimeout,
		IsStarted,
		IsScoreSent,
	}
}

func NewKahootGameTrivias(host User, trivias []Question) KahootGame {
	users := make([]User, 0)
	pin := "123534" //TODO: Use random generator
	CurrentQuestion := 0
	TotalQuestions := 0
	IsTimeout := false
	IsStarted := false
	IsScoreSent := false
	return KahootGame{
		trivias,
		host,
		users,
		pin,
		CurrentQuestion,
		TotalQuestions,
		IsTimeout,
		IsStarted,
		IsScoreSent,
	}
}

func (kg *KahootGame) PublishTrivia(trivia Question) (int, error) {
	// TODO: Validate trivia
	kg.trivias = append(kg.trivias, trivia)
	return trivia.ID, nil
}

func (kg *KahootGame) AddUser(user User) error {
	// TODO: Validate user
	kg.users = append(kg.users, user)
	return nil
}

func (kg *KahootGame) GetUsers() []User {
	return kg.users
}

func (kg *KahootGame) GetPin() string {
	return kg.pin
}