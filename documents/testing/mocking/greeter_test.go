package main

import (
    "testing"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

// 1) Embed mockMock
type mockGreeter struct {
	mock.Mock
}

// 2) Implement Greet by delegating to Called(...)
func (m *mockGreeter) Greet(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func TestService_Hello(t *testing.T) {
	m := new(mockGreeter)	

	// 3) Script the behavior
	m.On("Greet", "Alice").Return("Hello, Alice!").Once() // expect exactly one call

	// setup
	svc := Service{G: m}
	res := svc.Hello("Alice")

	require.Equal(t, "Hello, Alice!", res)

	// 4) Verify
	m.AssertExpectations(t)	
}