package handler

import (
	"fmt"
	"strconv"

	"github.com/amannlalwani/Notes-app-using-gofr/models"
	"github.com/amannlalwani/Notes-app-using-gofr/store"
	"gofr.dev/pkg/gofr"
)

type handler struct {
	store store.Store
}

func New(input_store store.Store) handler {
	return handler{store: input_store}
}

type response struct {
	Notes []models.Note
}

func (h handler) Get(ctx *gofr.Context) (interface{}, error) {
	resp, err := h.store.Get(ctx)
	if err != nil {
		panic(err)
	}
	return response{Notes: resp}, nil
}

func (h handler) Create(ctx *gofr.Context) (interface{}, error) {

	var new_note models.Note

	err := ctx.Bind(&new_note)

	if err != nil {
		panic(err)
	}
	//checking as there is should be no empty value in title and content
	if new_note.Title == "" || new_note.Content == "" {
		return "Please enter a non-empty value in Title and Content ", nil
	}

	resp, err := h.store.Create(ctx, new_note)

	if err != nil {

		panic(err)
	}

	fmt.Println("Created Success.")

	return resp, nil

}

func (h handler) Update(ctx *gofr.Context) (interface{}, error) {
	var new_note models.Note

	id := ctx.PathParam("id")

	err := ctx.Bind(&new_note)
	if err != nil {
		panic(err)
	}

	int_id, err := strconv.Atoi(id)

	if err != nil {
		panic(err)
	}
	//checking as there is should be no empty value in title and content
	if new_note.Title == "" || new_note.Content == "" {
		return "Please enter a non-empty value in Title and Content ", nil
	}

	resp, err := h.store.Update(ctx, int_id, new_note)

	if err != nil {
		panic(err)
	}

	fmt.Println("Updated Success.")

	return resp, nil

}

func (h handler) Delete(ctx *gofr.Context) (interface{}, error) {

	id := ctx.PathParam("id")

	int_id, err := strconv.Atoi(id)

	if err != nil {
		panic(err)
	}

	resp, err := h.store.Delete(ctx, int_id)

	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted Success.")

	return resp, nil

}
