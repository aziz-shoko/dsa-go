package main

type Database interface {
	GetUser(id int) (string, error)
}

type UserService struct {
	db Database
}

func (s *UserService) GetUserName(id int) (string, error) {
	return s.db.GetUser(id)
}