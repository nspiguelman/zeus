package domain

type KahootGame struct {
	trivias      []Trivia
	users		[]User
	pin 		string
}

func NewKahootGame() KahootGame {
	trivias := make([]Trivia, 0)
	users := make([]User, 0)
	pin := "123534" //TODO: Use random generator
	return KahootGame{trivias, users, pin}
}

func (kg *KahootGame) PublishTrivia(trivia Trivia) (int, error) {
	// TODO: Validate trivia
	kg.trivias = append(kg.trivias, trivia)
	return trivia.GetId(), nil
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