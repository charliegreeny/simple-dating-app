package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	m "github.com/charliegreeny/simple-dating-app/mock"
	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"testing"
	"time"
)

var mockUUID = func() string {
	return "uuid"
}

func Test_service_Get(t *testing.T) {
	dob, _ := time.Parse("2006-01-02", "1999-12-31")
	mock, gormDb, db := m.MockDb(t)
	defer db.Close()
	tests := []struct {
		name     string
		id       string
		stubRows *sqlmock.Rows
		stubErr  error
		want     *Output
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "successfully return user entity and nil error",
			id:   "id",
			stubRows: sqlmock.NewRows([]string{
				"id", "name", "gender", "date_of_birth", "email", "password"}).
				AddRow("id", "First Last", "FEMALE", dob, "test@email.com", "password"),
			stubErr: nil,
			want: &Output{
				ID:       "id",
				Name:     "First Last",
				Gender:   "FEMALE",
				Age:      24,
				Email:    "test@email.com",
				Password: "password",
			},
			wantErr: assert.NoError,
		},
		{
			name: "no user found returns nil output and ErrNotFound",
			id:   "id",
			stubRows: sqlmock.NewRows([]string{
				"id", "name", "gender", "date_of_birth", "email", "password"}),
			stubErr: gorm.ErrRecordNotFound,
			want:    nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.As(err, &app.ErrNotFound{})
			},
		},
		{
			name: "general db error returns nil and wrapped error",
			id:   "id",
			stubRows: sqlmock.NewRows([]string{
				"id", "name", "gender", "date_of_birth", "email", "password"}),
			stubErr: gorm.ErrInvalidDB,
			want:    nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, fmt.Errorf("db error getting user: %w", gorm.ErrInvalidDB))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				db:  gormDb,
				log: zap.NewNop(),
			}
			m := mock.ExpectQuery("SELECT (.+) FROM `users` WHERE id =(.+)").
				WithArgs(tt.id, 1).WillReturnRows(tt.stubRows)

			if tt.stubErr != nil {
				m.WillReturnError(tt.stubErr)
			}

			got, err := s.Get(context.Background(), tt.id)
			if tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_service_Create(t *testing.T) {
	tests := []struct {
		name       string
		input      *Input
		want       *Output
		stubSQLErr error
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "successfully insert and return user and nil error",
			input: &Input{
				Name:     "First Last",
				Email:    "test@email.com",
				Gender:   "FEMALE",
				Dob:      "2000-01-01",
				Password: "pwd",
			},
			want: &Output{
				ID:       mockUUID(),
				Email:    "test@email.com",
				Name:     "First Last",
				Gender:   "FEMALE",
				Age:      24,
				Password: "pwd",
			},
			stubSQLErr: nil,
			wantErr:    assert.NoError,
		},
		{
			name: "invalid date format returns nil and ErrBadRequest error",
			input: &Input{
				Name:     "First Last",
				Email:    "test@email.com",
				Gender:   "FEMALE",
				Dob:      "01-01-2000",
				Password: "pwd",
			},
			want:       nil,
			stubSQLErr: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, app.ErrBadRequest{}) &&
					err.Error() == "invalid date format provided - needs to be yyyy-mm-dd"
			},
		},
		{
			name: "18 today returns Output and nil error",
			input: &Input{
				Name:     "First Last",
				Email:    "test@email.com",
				Gender:   "FEMALE",
				Dob:      time.Now().AddDate(-18, 0, 0).Format(time.DateOnly),
				Password: "pwd",
			},
			want: &Output{
				ID:       mockUUID(),
				Email:    "test@email.com",
				Name:     "First Last",
				Gender:   "FEMALE",
				Age:      18,
				Password: "pwd",
			},
			stubSQLErr: nil,
			wantErr:    assert.NoError,
		},
		{
			name: "under 18 returns nil and ErrForbidden error",
			input: &Input{
				Name:     "First Last",
				Email:    "test@email.com",
				Gender:   "FEMALE",
				Dob:      time.Now().AddDate(-18, 00, 1).Format(time.DateOnly),
				Password: "pwd",
			},
			want:       nil,
			stubSQLErr: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, app.ErrForbidden{}) &&
					err.Error() == "user must be at least 18 years old"
			},
		},
		{
			name: "invalid gender returns nil and ErrBadRequest error",
			input: &Input{
				Name:     "First Last",
				Email:    "test@email.com",
				Gender:   "INVALID_GENDER",
				Dob:      "2000-01-01",
				Password: "pwd",
			},
			want: nil,
			stubSQLErr: &mysql2.MySQLError{
				Number:  foreignKeyErrCode,
				Message: "foreign key constraint failed",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, app.ErrBadRequest{}) && err.Error() == "invalid gender provided"
			},
		},
		{
			name: "non unique email returns nil and ErrBadRequest error",
			input: &Input{
				Name:     "First Last",
				Email:    "test@email.com",
				Gender:   "INVALID_GENDER",
				Dob:      "2000-01-01",
				Password: "pwd",
			},
			want: nil,
			stubSQLErr: &mysql2.MySQLError{
				Number:   duplicateKeyErrCode,
				SQLState: [5]byte{},
				Message:  "duplicate key error",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, app.ErrBadRequest{}) && err.Error() == "email already in use"
			},
		},
	}
	for _, tt := range tests {
		mock, gormDb, db := m.MockDb(t)
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				db:      gormDb,
				uuidGen: mockUUID,
				log:     zap.NewNop(),
				cache:   NewCache(gormDb),
			}
			mock.ExpectBegin()
			m := mock.ExpectExec("INSERT INTO `users` (.+) VALUES (.+)").
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.stubSQLErr != nil {
				mock.ExpectRollback()
				m.WillReturnError(tt.stubSQLErr)
			} else {
				mock.ExpectCommit()
			}

			got, err := s.Create(context.Background(), tt.input)
			tt.wantErr(t, err, fmt.Sprintf("Create(%v)", tt.input))
			assert.Equalf(t, tt.want, got, "Create(%v)", tt.input)
			_ = db.Close()
		})
	}
}
