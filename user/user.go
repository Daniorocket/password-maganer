package user

type User struct {
	Login    string
	IsLogged bool
	Key      []byte
}

func InitUser() *User {
	return &User{
		Login:    "",
		IsLogged: false,
		Key:      nil,
	}
}
