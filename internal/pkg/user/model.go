package user

import "time"

type Input struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Dob      string `json:"dateOfBirth" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Output struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Password string `json:"password"`
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

func (e Entity) toOutput() *Output {
	age := calcAge(e.DateOfBirth)
	return &Output{
		ID:       e.ID,
		Name:     e.Name,
		Gender:   e.Gender,
		Email:    e.Email,
		Password: e.Password,
		Age:      age,
	}
}
