package users

import (
	"bytes"
	"context"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"gofr-crud/models"
	"gofr-crud/services"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_Create(t *testing.T) {
	app := gofr.New()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockService := services.NewMockUsers(ctrl)
	handler := New(mockService)
	ctx := gofr.NewContext(nil, nil, app)
	ctx.Context = context.Background()

	user := models.User{ID: 1, Name: "David", Age: 37}
	testCases := []struct {
		desc string
		body models.User
		err  error
		mock *gomock.Call
	}{
		{desc: "success case",
			body: user,
			err:  nil,
			mock: mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil),
		},

		{
			desc: "failure case",
			body: user,
			err:  errors.Error("DB error"),
			mock: mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.Error("DB error")),
		},
		{
			desc: "failure case",
			body: models.User{ID: 1},
			err:  errors.InvalidParam{Param: []string{"body"}},
			mock: nil,
		},
	}
	for i, tc := range testCases {
		pr, _ := json.Marshal(tc.body)
		req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(pr))
		res := httptest.NewRecorder()

		r := request.NewHTTPRequest(req)
		w := responder.NewContextualResponder(res, req)
		ctx := gofr.NewContext(w, r, app)
		_, err := handler.Create(ctx)
		if err != nil && !reflect.DeepEqual(err.Error(), tc.err.Error()) {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

	}

}
