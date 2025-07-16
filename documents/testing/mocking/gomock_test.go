package main

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGetUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockDatabase(ctrl)

	// Set expectations
	mockDB.EXPECT().
		GetUser(42).
		Return("Bob", nil)

	service := &UserService{
		db: mockDB,
	}

	name, err := service.GetUserName(42)
	if err != nil {
		t.Fatal(err)
	}

	if name != "Bob" {
		t.Errorf("expected Bob, got %q", name)
	}
}
