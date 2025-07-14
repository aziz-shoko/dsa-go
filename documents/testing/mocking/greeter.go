package main

type Greeter interface {
	Greet(name string) string
}

type Service struct {
	G Greeter
}

func (s *Service) Hello(name string) string {
	return s.G.Greet(name)
}