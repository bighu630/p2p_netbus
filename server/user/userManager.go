package user

type UserManager struct {
	Users map[string]*User // userID => User
}

func NewUserManager() *UserManager {
	return &UserManager{
		Users: map[string]*User{},
	}
}
