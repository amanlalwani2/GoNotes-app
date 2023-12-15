package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/amannlalwani/Notes-app-using-gofr/models"
	"github.com/amannlalwani/Notes-app-using-gofr/store"
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/request"
)

func initializeHandlerTest(t *testing.T) (*store.MockStore, handler, *gofr.Gofr) {
	ctrl := gomock.NewController(t)

	mockStore := store.NewMockStore(ctrl)
	h := New(mockStore)
	app := gofr.New()

	return mockStore, h, app
}

func TestGet(t *testing.T) {

	tests := []struct {
		desc string
		resp []models.Note
		err  error
	}{
		{"success case", []models.Note{{ID: 1, Title: "sample title", Content: "sample content"}}, nil},
		{"error case", nil, errors.Error("error fetching notes")},
	}

	mockstore, h, app := initializeHandlerTest(t)

	for _, testcase := range tests {
		req := httptest.NewRequest(http.MethodGet, "/notes", nil)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, app)

		mockstore.EXPECT().Get(ctx).Return(testcase.resp, testcase.err)
		result, err := h.Get(ctx)

		if testcase.err == nil {
			assert.Nil(t, err)
			assert.NotNil(t, result)
			res, ok := result.(response)
			assert.True(t, ok)
			assert.Equal(t, testcase.resp, res.Notes)
		} else {
			assert.NotNil(t, err)
			assert.Equal(t, testcase.err, err)
			assert.Nil(t, result)
		}
	}
}

func TestCreate(t *testing.T) {
	mockstore, h, app := initializeHandlerTest(t)

	input := `{"note_id":5,"title":"Grocery List","content":"Milk,bread,and bananas."}`
	expRes := models.Note{ID: 5, Title: "Grocery List", Content: "Milk,bread,and bananas."}

	inp := strings.NewReader(input)
	req := httptest.NewRequest(http.MethodPost, "/notes", inp)
	r := request.NewHTTPRequest(req)
	ctx := gofr.NewContext(nil, r, app)

	var note models.Note

	_ = ctx.Bind(&note)

	mockstore.EXPECT().Get(ctx).Return(nil, nil).MaxTimes(2)
	mockstore.EXPECT().Create(ctx, note).Return(expRes, nil).MaxTimes(1)

	res, err := h.Create(ctx)

	assert.Nil(t, err, "Test:failed case")

	assert.Equal(t, expRes, res, "Test:success case")

}

func TestCreate_Error(t *testing.T) {
	mockstore, h, app := initializeHandlerTest(t)

	tests := []struct {
		desc   string
		inp    string
		expRes interface{}
		err    error
	}{{"trying to create invalid body", `{"note_id":2 ""Stilte":"Grocery List","content":"Milk,bread,and bananas."}`, models.Note{}, errors.InvalidParam{Param: []string{"body"}}},
		{"trying to create invalid body", `{{"}`, models.Note{}, errors.InvalidParam{Param: []string{"body"}}},
	}

	for index, testcase := range tests {
		input := strings.NewReader(testcase.inp)
		req := httptest.NewRequest(http.MethodPost, "/notes", input)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, app)

		var note models.Note
		_ = ctx.Bind(&note)

		mockstore.EXPECT().Create(ctx, note).Return(testcase.expRes.(models.Note), nil).MaxTimes(1)

		res, err := h.Create(ctx)

		assert.Equal(t, testcase.err, err, "Test[%d],failed.%s\n", index, testcase.desc)
		assert.Nil(t, res, "Test[%d],failed.%s\n", index, testcase.desc)

	}
}

func TestUpdate(t *testing.T) {
	mockstore, h, app := initializeHandlerTest(t)

	input := `{"note_id":5,"title":"Grocery List","content":"Milk,bread,and bananas."}`
	expRes := models.Note{ID: 5, Title: "Grocery List", Content: "Milk,bread,and bananas."}

	inp := strings.NewReader(input)
	req := httptest.NewRequest(http.MethodPut, "/notes/5", inp)
	r := request.NewHTTPRequest(req)
	ctx := gofr.NewContext(nil, r, app)

	var note models.Note

	_ = ctx.Bind(&note)

	mockstore.EXPECT().Get(ctx).Return(nil, nil).MaxTimes(2)
	mockstore.EXPECT().Update(ctx, 5, note).Return(expRes, nil).MaxTimes(1)

	ctx.SetPathParams(map[string]string{"id": "5"})

	res, err := h.Update(ctx)

	assert.Nil(t, err, "Test:failed case")

	assert.Equal(t, expRes, res, "Test:success case")

}

func TestUpdate_Error(t *testing.T) {
	mockstore, h, app := initializeHandlerTest(t)

	tests := []struct {
		desc   string
		inp    string
		expRes interface{}
		err    error
	}{{"trying to update with invalid body", `{"note_id":2 ""Stilte":"Grocery List","content":"Milk,bread,and bananas."}`, models.Note{}, errors.InvalidParam{Param: []string{"body"}}},
		{"trying to update with invalid body", `{{"}`, models.Note{}, errors.InvalidParam{Param: []string{"body"}}},
	}

	for index, testcase := range tests {
		input := strings.NewReader(testcase.inp)
		req := httptest.NewRequest(http.MethodPost, "/notes/3", input)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, app)

		var note models.Note
		_ = ctx.Bind(&note)

		mockstore.EXPECT().Update(ctx, 3, note).Return(testcase.expRes.(models.Note), nil).MaxTimes(1)

		res, err := h.Update(ctx)

		assert.Equal(t, testcase.err, err, "Test[%d],failed.%s\n", index, testcase.desc)
		assert.Nil(t, res, "Test[%d],failed.%s\n", index, testcase.desc)

	}
}

func TestDelete(t *testing.T) {
	mockstore, h, app := initializeHandlerTest(t)

	input := `{}`
	expRes := models.Note{ID: 5, Title: "Grocery List", Content: "Milk,bread,and bananas."}

	inp := strings.NewReader(input)
	req := httptest.NewRequest(http.MethodDelete, "/notes/5", inp)
	r := request.NewHTTPRequest(req)
	ctx := gofr.NewContext(nil, r, app)

	mockstore.EXPECT().Get(ctx).Return(nil, nil).MaxTimes(2)
	mockstore.EXPECT().Delete(ctx, 5).Return(expRes, nil).MaxTimes(1)

	ctx.SetPathParams(map[string]string{"id": "5"})

	res, err := h.Delete(ctx)

	assert.Nil(t, err, "Test:failed case")

	assert.Equal(t, expRes, res, "Test:success case")

}

func TestDelete_Error(t *testing.T) {
	mockstore, h, app := initializeHandlerTest(t)

	tests := []struct {
		desc   string
		inp    string
		expRes interface{}
		err    error
	}{{"trying to update with invalid body", `{}`, models.Note{}, errors.InvalidParam{Param: []string{"body"}}},
		{"trying to update with invalid body", `{{"}`, models.Note{}, errors.InvalidParam{Param: []string{"body"}}},
	}

	for index, testcase := range tests {
		input := strings.NewReader(testcase.inp)
		req := httptest.NewRequest(http.MethodDelete, "/notes/3", input)
		r := request.NewHTTPRequest(req)
		ctx := gofr.NewContext(nil, r, app)

		mockstore.EXPECT().Delete(ctx, 3).Return(testcase.expRes.(models.Note), nil).MaxTimes(1)

		res, err := h.Delete(ctx)

		assert.Equal(t, testcase.err, err, "Test[%d],failed.%s\n", index, testcase.desc)
		assert.Nil(t, res, "Test[%d],failed.%s\n", index, testcase.desc)

	}
}
