package user

type Usecase interface {
	Create(user *User) error
	Update(user *User) error
}

type Repository interface {
	Create(user *User) error
	Update(user *User) error
}
