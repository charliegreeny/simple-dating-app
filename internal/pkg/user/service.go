package user

import (
	"errors"
	"fmt"
	"github.com/charliegreeny/simple-dating-app/internal/model"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type service struct {
	db  *gorm.DB
	log *zap.Logger
}

const duplicateKeyErrCode = 1062
const foreignKeyErrCode = 1452

func NewService(db *gorm.DB, l *zap.Logger) model.GetterCreator[*Input, *Output] {
	return &service{db: db, log: l}
}

func (s service) Get(id string) (*Output, error) {
	s.log.Debug("Getting user", zap.String("id", id))
	var e *Entity
	r := s.db.First(&e, "id = ?", id)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound{Message: fmt.Sprintf("no user found for id %s", id)}
		}
	}
	return e.toOutput(), nil
}

func (s service) Create(input *Input) (*Output, error) {
	dob, err := time.Parse(time.DateOnly, input.Dob)
	if err != nil {
		return nil, model.ErrBadRequest{Message: "invalid date format provided - needs to be yyyy-mm-dd"}
	}
	e := &Entity{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Gender:      input.Gender,
		DateOfBirth: dob,
		Email:       input.Email,
		Password:    input.Password,
	}
	s.log.Debug("Creating user", zap.String("id", e.ID))
	err = s.db.Create(e).Error
	if err != nil {
		var sqlErr *mysql.MySQLError
		ok := errors.As(err, &sqlErr)
		if ok && sqlErr.Number == duplicateKeyErrCode {
			return nil, model.ErrBadRequest{Message: "email already in use"}
		}
		if ok && sqlErr.Number == foreignKeyErrCode {
			return nil, model.ErrBadRequest{Message: "gender provided not supported"}
		}
		return nil, err
	}
	return e.toOutput(), nil
}
