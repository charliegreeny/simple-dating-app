package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/preference"
	appMock "github.com/charliegreeny/simple-dating-app/mock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	tests := []struct {
		name         string
		id           string
		stubRows     *sqlmock.Rows
		stubSQLErr   error
		stubCacheErr error
		want         *app.UserOutput
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "successfully return cached entity and nil error",
			id:   "id",
			stubRows: sqlmock.NewRows([]string{
				"id", "name", "gender", "date_of_birth", "email", "password"}).
				AddRow("id", "First Last", "FEMALE", dob, "test@email.com", "password"),
			stubSQLErr: nil,
			want: &app.UserOutput{
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
			name:         "successfully return entity from db and nil error",
			id:           "id",
			stubCacheErr: app.ErrNotFound{Message: "not found"},
			stubRows: sqlmock.NewRows([]string{
				"id", "name", "gender", "date_of_birth", "email", "password"}).
				AddRow("id", "First Last", "FEMALE", dob, "test@email.com", "password"),
			stubSQLErr: nil,
			want: &app.UserOutput{
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
			name: "no cache and no db record returns nil output and ErrNotFound",
			id:   "id",
			stubRows: sqlmock.NewRows([]string{
				"id", "name", "gender", "date_of_birth", "email", "password"}),
			stubSQLErr:   gorm.ErrRecordNotFound,
			stubCacheErr: app.ErrNotFound{Message: "not found"},
			want:         nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.As(err, &app.ErrNotFound{})
			},
		},
		{
			name: "general db error returns nil and wrapped error",
			id:   "id",
			stubRows: sqlmock.NewRows([]string{
				"id", "name", "gender", "date_of_birth", "email", "password"}),
			stubCacheErr: app.ErrNotFound{Message: "not found"},
			stubSQLErr:   gorm.ErrInvalidDB,
			want:         nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, fmt.Errorf("db error getting cache: %w", gorm.ErrInvalidDB))
			},
		},
	}
	for _, tt := range tests {
		sqlMock, gormDb, db := appMock.MockDb(t)
		defer db.Close()
		t.Run(tt.name, func(t *testing.T) {
			mockCache := &appMock.UserCache{}
			mockCache.On("Get", context.Background(), mock.Anything).
				Return(&app.User{UserOutput: tt.want}, tt.stubCacheErr).Maybe()
			s := service{
				db:    gormDb,
				log:   zap.NewNop(),
				cache: mockCache,
			}
			m := sqlMock.ExpectQuery("SELECT (.+) FROM `users` WHERE id =(.+)").
				WithArgs(tt.id, 1).WillReturnRows(tt.stubRows)

			if tt.stubSQLErr != nil {
				m.WillReturnError(tt.stubSQLErr)
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
		want       *app.UserOutput
		stubSQLErr error
		wantErr    assert.ErrorAssertionFunc
		cacheCalls int
	}{
		{
			name: "successfully insert and return cache and nil error",
			input: &Input{
				Name:     "First Last",
				Email:    "test@email.com",
				Gender:   "FEMALE",
				Dob:      "2000-01-01",
				Password: "pwd",
			},
			want: &app.UserOutput{
				ID:       mockUUID(),
				Email:    "test@email.com",
				Name:     "First Last",
				Gender:   "FEMALE",
				Age:      24,
				Password: "pwd",
			},
			stubSQLErr: nil,
			cacheCalls: 1,
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
			name: "18 today returns UserOutput and nil error",
			input: &Input{
				Name:     "First Last",
				Email:    "test@email.com",
				Gender:   "FEMALE",
				Dob:      time.Now().AddDate(-18, 0, 0).Format(time.DateOnly),
				Password: "pwd",
			},
			want: &app.UserOutput{
				ID:       mockUUID(),
				Email:    "test@email.com",
				Name:     "First Last",
				Gender:   "FEMALE",
				Age:      18,
				Password: "pwd",
			},
			cacheCalls: 1,
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
					err.Error() == "cache must be at least 18 years old"
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
			stubSQLErr: &mysql.MySQLError{
				Number:  app.ForeignKeyErrCode,
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
			stubSQLErr: &mysql.MySQLError{
				Number:   app.DuplicateKeyErrCode,
				SQLState: [5]byte{},
				Message:  "duplicate key error",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, app.ErrBadRequest{}) && err.Error() == "email already in use"
			},
		},
	}
	for _, tt := range tests {
		sqlMock, gormDb, db := appMock.MockDb(t)
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := &appMock.UserCache{}
			cacheMock.On("Add", context.Background(), mock.MatchedBy(func(u *app.User) bool {
				return assert.Equal(t, &app.User{
					UserOutput:   tt.want,
					Pref:         preference.DefaultPreferences(tt.want.ID, tt.want.Gender),
					Loc:          nil,
					DistanceFrom: 0,
				}, u)
			})).
				Return(nil).Times(tt.cacheCalls)

			s := service{
				db:      gormDb,
				uuidGen: mockUUID,
				log:     zap.NewNop(),
				cache:   cacheMock,
			}
			sqlMock.ExpectBegin()
			m := sqlMock.ExpectExec("INSERT INTO `users` (.+) VALUES (.+)").
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.stubSQLErr != nil {
				sqlMock.ExpectRollback()
				m.WillReturnError(tt.stubSQLErr)
			} else {
				sqlMock.ExpectCommit()
			}

			got, err := s.Create(context.Background(), tt.input)
			tt.wantErr(t, err, fmt.Sprintf("Create(%v)", tt.input))
			assert.Equalf(t, tt.want, got, "Create(%v)", tt.input)
			_ = db.Close()
		})
	}
}
