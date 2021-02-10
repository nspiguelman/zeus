package domain

type Kahootuser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Token string 	`json:"token"`
	KahootId int	`json:"kahootId"`
}

func NewUser(username string, token string, kahootID int) *Kahootuser{
	return &Kahootuser{
		Username: username,
		Token:    token,
		KahootId: kahootID,
	}
}