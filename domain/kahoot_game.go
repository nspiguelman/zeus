package domain

type KahootGame struct {
	trivias      []Question
	host		User
	users		[]User
	pin 		string
}

func NewKahootGame(host User) KahootGame {
	trivias := make([]Question, 0)
	users := make([]User, 0)
	pin := "123534" //TODO: Use random generator
	return KahootGame{trivias,host,  users, pin}
}

func NewKahootGameTrivias(host User, trivias []Question) KahootGame {
	users := make([]User, 0)
	pin := "123534" //TODO: Use random generator
	return KahootGame{trivias, host, users, pin}
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