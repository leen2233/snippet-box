package mocks

import (
  "time"
  "snippetbox.leen2233.me/internal/models"
)

var mockUser = &models.User{
  ID: 2,
  Name: "user",
  Email: "test@test.com",
  HashedPassword: []byte("password"),
  Created:  time.Now(),
}


type UserModel struct {}


func (m *UserModel) Insert(name, email, password string) error {
  return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
  if email == "test@test.com" {
    return 2, nil
  }
  return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
  if id == 2 {
    return true, nil
  }
  return false, nil
}

