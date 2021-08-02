package heroku_store

type UserStore interface {
	Create(user *User) (*User, error)
	List() ([]User, error)
}
