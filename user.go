package ddd

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserService interface {
	User(id int) (*User, error)
}
