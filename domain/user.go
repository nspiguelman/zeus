package domain

type User struct {
	Nickname string
	Password string
}

func NewUser(nickname string, password string) *User {
	return &User{nickname, password}
}

func (user *User) GetNickname() string {
	return user.Nickname
}

func (user *User) SetNickname(nickname string) {
	user.Nickname = nickname
}

func (user *User) SetPassword(psw string) {
	user.Password = psw
}
