package service

import (
	"github.com/charliegreeny/simple-dating-app/app"
	"time"
)

type Input struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Dob      string `json:"dateOfBirth" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Entity struct {
	ID          string    `gorm:"column:id; primaryKey"`
	Name        string    `gorm:"column:name"`
	Gender      string    `gorm:"column:gender"`
	DateOfBirth time.Time `gorm:"column:date_of_birth"`
	Email       string    `gorm:"column:email"`
	Password    string    `gorm:"column:password"`
}

func (e Entity) TableName() string {
	return "users"
}

func (e Entity) ToOutput() *app.UserOutput {
	age := calcAge(e.DateOfBirth)
	return &app.UserOutput{
		ID:       e.ID,
		Name:     e.Name,
		Gender:   e.Gender,
		Email:    e.Email,
		Password: e.Password,
		Age:      age,
	}
}
