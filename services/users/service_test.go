package users

import (
	"context"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/golang/mock/gomock"
	"gofr-crud/models"
	"gofr-crud/stores"
	"reflect"
	"strconv"
	"testing"
)

func TestService_Create(t *testing.T) {
	app := gofr.New()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockStore := stores.NewMockUsers(ctrl)
	mock := New(mockStore)
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	user := models.User{ID: 1, Name: "David", Age: 37}
	testCases := []struct {
		desc  string
		input models.User
		err   error
		mock  *gomock.Call
	}{
		{desc: "success case",
			input: user,
			err:   nil,
			mock:  mockStore.EXPECT().Create(ctx, &user).Return(nil),
		},

		{
			desc:  "failure case",
			input: user,
			err:   errors.Error("DB error"),
			mock:  mockStore.EXPECT().Create(ctx, &user).Return(errors.Error("DB error")),
		},
	}

	for i, tc := range testCases {
		err := mock.Create(ctx, &user)
		if err != nil && !reflect.DeepEqual(err.Error(), tc.err.Error()) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestService_GetAll(t *testing.T) {
	app := gofr.New()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockStore := stores.NewMockUsers(ctrl)
	mock := New(mockStore)
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	user := []models.User{
		{ID: 1, Name: "David", Age: 37},
		{ID: 2, Name: "Jai", Age: 47},
	}

	testCases := []struct {
		desc   string
		output []models.User
		err    error
		mock   *gomock.Call
	}{
		{
			desc:   "success case",
			output: user,
			err:    nil,
			mock:   mockStore.EXPECT().GetAll(ctx).Return(user, nil),
		},
		{
			desc:   "failure case",
			output: nil,
			err:    errors.Error("error"),
			mock:   mockStore.EXPECT().GetAll(ctx).Return(nil, errors.Error("error")),
		},
	}
	for i, tc := range testCases {
		user, err := mock.GetAll(ctx)

		if err != nil && !reflect.DeepEqual(err.Error(), tc.err.Error()) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if !reflect.DeepEqual(user, tc.output) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestService_GetByID(t *testing.T) {
	app := gofr.New()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	mockStore := stores.NewMockUsers(ctrl)
	mock := New(mockStore)
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()
	user := models.User{
		ID:   1,
		Name: "alex",
		Age:  37,
	}
	testCases := []struct {
		desc   string
		id     int
		output models.User
		err    error
		mock   *gomock.Call
	}{
		{
			desc:   "success case",
			id:     1,
			output: user,
			err:    nil,
			mock:   mockStore.EXPECT().GetByID(ctx, 1).Return(user, nil),
		},
		{
			desc:   "invalid id",
			id:     -3,
			output: models.User{},
			err:    errors.InvalidParam{Param: []string{"id"}},
			mock:   nil,
		},
		{
			desc:   "failure case",
			id:     1,
			output: models.User{},
			err:    errors.Error("error"),
			mock:   mockStore.EXPECT().GetByID(ctx, 1).Return(models.User{}, errors.Error("error")),
		},
	}

	for i, tc := range testCases {
		user, err := mock.GetByID(ctx, tc.id)

		if err != nil && !reflect.DeepEqual(err.Error(), tc.err.Error()) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if !reflect.DeepEqual(user, tc.output) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestService_Update(t *testing.T) {
	app := gofr.New()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockStore := stores.NewMockUsers(ctrl)
	mock := New(mockStore)
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()
	user := models.User{
		ID:   1,
		Name: "alex",
		Age:  37,
	}
	user1 := models.User{
		ID:   -1,
		Name: "alex",
		Age:  37,
	}

	var testCases = []struct {
		desc string
		id   int
		user models.User
		err  error
		mock []*gomock.Call
	}{
		{desc: "success case",
			id:   user.ID,
			user: user,
			err:  nil,
			mock: []*gomock.Call{
				mockStore.EXPECT().GetByID(ctx, user.ID).Return(models.User{}, nil),
				mockStore.EXPECT().Update(ctx, &user).Return(nil),
			},
		},
		{desc: "id < 0",
			id:   -1,
			user: user1,
			err:  errors.InvalidParam{Param: []string{"id"}},
			mock: nil,
		},
		{desc: "failure case",
			id:   user.ID,
			user: user,
			err: errors.EntityNotFound{
				Entity: "user",
				ID:     strconv.Itoa(user.ID),
			},
			mock: []*gomock.Call{
				mockStore.EXPECT().GetByID(ctx, user.ID).Return(models.User{},
					errors.EntityNotFound{
						Entity: "user",
						ID:     strconv.Itoa(user.ID),
					})},
		},
		{desc: "failure case",
			id:   user.ID,
			user: user,
			err:  errors.Error("error"),
			mock: []*gomock.Call{
				mockStore.EXPECT().GetByID(ctx, user.ID).Return(models.User{}, nil),
				mockStore.EXPECT().Update(ctx, &user).Return(errors.Error("error")),
			},
		},
	}

	for i, tc := range testCases {
		err := mock.Update(ctx, &tc.user)
		if err != nil && !reflect.DeepEqual(err.Error(), tc.err.Error()) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestService_Delete(t *testing.T) {
	app := gofr.New()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockStore := stores.NewMockUsers(ctrl)
	mock := New(mockStore)
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()
	ID := 2
	testCases := []struct {
		desc string
		id   int
		err  error
		mock []*gomock.Call
	}{
		{
			desc: "success case",
			id:   1,
			err:  errors.Error("error"),
			mock: []*gomock.Call{
				mockStore.EXPECT().GetByID(ctx, 1).Return(models.User{}, nil),
				mockStore.EXPECT().Delete(ctx, 1).Return(errors.Error("error")),
			},
		},
		{desc: "id < 0",
			id:   -1,
			err:  errors.InvalidParam{Param: []string{"id"}},
			mock: nil,
		},
		{desc: "failure case",
			id: ID,
			err: errors.EntityNotFound{
				Entity: "user",
				ID:     strconv.Itoa(ID),
			},
			mock: []*gomock.Call{
				mockStore.EXPECT().GetByID(ctx, ID).Return(models.User{},
					errors.EntityNotFound{
						Entity: "user",
						ID:     strconv.Itoa(ID),
					})},
		},
		{desc: "failure case",
			id:  ID,
			err: nil,
			mock: []*gomock.Call{
				mockStore.EXPECT().GetByID(ctx, ID).Return(models.User{}, nil),
				mockStore.EXPECT().Delete(ctx, ID).Return(nil),
			},
		},
	}

	for i, tc := range testCases {
		err := mock.Delete(ctx, tc.id)
		if err != nil && !reflect.DeepEqual(err.Error(), tc.err.Error()) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
