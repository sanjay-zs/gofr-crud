package users

import (
	"context"
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gofr-crud/models"
	"reflect"
	"testing"
)

func getTestData() (*gofr.Gofr, sqlmock.Sqlmock, *sql.DB) {
	app := gofr.New()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		fmt.Printf("error %s was not expected when opening a stub database connection\n", err)
	}

	return app, mock, db
}

func TestCreate(t *testing.T) {
	app, mock, db := getTestData()
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()
	store := New()
	ctx.DB().DB = db

	user := models.User{
		ID:   1,
		Name: "David",
		Age:  37,
	}

	queryErr := errors.Error("Query error")

	mock.ExpectExec(insertUser).
		WithArgs(user.ID, user.Name, user.Age).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(insertUser).
		WithArgs(user.ID, user.Name, user.Age).
		WillReturnError(queryErr)

	cases := []struct {
		desc string
		err  error
	}{
		{"success case", nil},
		{"failure case", errors.DB{Err: queryErr}},
	}
	for i, tc := range cases {
		err := store.Create(ctx, &user)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestGetAll(t *testing.T) {
	app, mock, db := getTestData()
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()
	store := New()
	ctx.DB().DB = db

	users := []models.User{
		{
			ID:   1,
			Name: "David",
			Age:  37,
		},
	}

	row1 := sqlmock.NewRows([]string{"ID", "Name", "Age"}).AddRow(1, "David", 37)
	row2 := sqlmock.NewRows([]string{"Name", "Age"}).AddRow("sanjay", 24)
	queryErr := errors.Error("Query error")
	scanErr := errors.Error("Scan error")

	testCases := []struct {
		desc   string
		output []models.User
		err    error
		mock   *sqlmock.ExpectedQuery
	}{
		{desc: "success case", output: users, err: nil,
			mock: mock.ExpectQuery(getUsers).WillReturnRows(row1).WillReturnError(nil)},
		{desc: "failure case", output: nil, err: errors.DB{Err: queryErr},
			mock: mock.ExpectQuery(getUsers).WillReturnError(queryErr)},
		{desc: "failure case", output: nil, err: errors.DB{Err: scanErr},
			mock: mock.ExpectQuery(getUsers).WithArgs(row2).WillReturnError(errors.DB{Err: scanErr})},
	}

	for i, tc := range testCases {
		users, err := store.GetAll(ctx)
		if err != nil && !reflect.DeepEqual(err.Error(), tc.err.Error()) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if !reflect.DeepEqual(users, tc.output) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestGetByID(t *testing.T) {
	app, mock, db := getTestData()
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()
	store := New()
	ctx.DB().DB = db

	user := models.User{ID: 1, Name: "David", Age: 37}
	row1 := sqlmock.NewRows([]string{"ID", "Name", "Age"}).AddRow(1, "David", 37)
	testCases := []struct {
		desc   string
		id     int
		output models.User
		err    error
		mock   *sqlmock.ExpectedQuery
	}{
		{
			desc:   "success case",
			id:     1,
			output: user,
			err:    nil,
			mock:   mock.ExpectQuery(getUserByID).WithArgs(1).WillReturnRows(row1).WillReturnError(nil),
		},
		{
			desc:   "failure case",
			id:     5,
			output: models.User{},
			err:    errors.Error("DB error"),
			mock:   mock.ExpectQuery(getUserByID).WithArgs(5).WillReturnRows().WillReturnError(errors.Error("DB error")),
		},
	}

	for i, tc := range testCases {
		users, err := store.GetByID(ctx, tc.id)
		if err != nil && !reflect.DeepEqual(err.Error(), tc.err.Error()) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if !reflect.DeepEqual(users, tc.output) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestUpdate(t *testing.T) {
	app, mock, db := getTestData()
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()
	store := New()
	ctx.DB().DB = db

	user := models.User{
		ID:   2,
		Name: "jai",
		Age:  40,
	}
	testCases := []struct {
		desc string
		err  error
		mock *sqlmock.ExpectedExec
	}{
		{
			desc: "success case",
			err:  nil,
			mock: mock.ExpectExec(updateUser).WithArgs(user.Name, user.Age, user.ID).WillReturnResult(sqlmock.NewResult(1, 0)),
		},
		{
			desc: "failure case",
			err:  errors.Error("DB error"),
			mock: mock.ExpectExec(updateUser).WithArgs(user.Name, user.Age, user.ID).WillReturnError(errors.Error("DB error")),
		},
	}

	for i, tc := range testCases {
		err := store.Update(ctx, &user)
		if err != nil && !reflect.DeepEqual(err.Error(), tc.err.Error()) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestDelete(t *testing.T) {
	app, mock, db := getTestData()
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()
	store := New()
	ctx.DB().DB = db

	testCases := []struct {
		desc string
		id   int
		err  error
		mock *sqlmock.ExpectedExec
	}{
		{
			desc: "success case",
			id:   1,
			err:  nil,
			mock: mock.ExpectExec(deleteUser).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1)),
		},
		{
			desc: "failure case",
			id:   1,
			err:  errors.Error("DB error"),
			mock: mock.ExpectExec(deleteUser).WithArgs(1).WillReturnError(errors.Error("DB error")),
		},
	}
	for i, tc := range testCases {
		err := store.Delete(ctx, tc.id)
		if err != nil && !reflect.DeepEqual(err.Error(), tc.err.Error()) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
