package store

import (
	"fmt"

	"github.com/amannlalwani/Notes-app-using-gofr/models"
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
)

type note struct{}

func New() Store {
	return note{}
}

func (n note) Get(ctx *gofr.Context) ([]models.Note, error) {

	//retrieving notes from db
	data, err := ctx.DB().QueryContext(ctx, "SELECT note_id,title,content FROM notes ")
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	fmt.Println("Data Retrieved Successfully")

	//making sure connection closes at end of this block.
	defer data.Close()

	notes := make([]models.Note, 0)

	//making slice of notes data retrieved from DB.
	for data.Next() {
		var n models.Note

		err = data.Scan(&n.ID, &n.Title, &n.Content)

		if err != nil {
			return nil, errors.DB{Err: err}
		}

		notes = append(notes, n)
	}

	err = data.Err()

	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return notes, nil
}

func (n note) Create(ctx *gofr.Context, inp models.Note) (models.Note, error) {
	var res models.Note

	//checking as ID should be positive
	if inp.ID < 0 {
		return models.Note{}, errors.Error("Please enter a positive ID")
	}

	//checking as there is should be no empty value in title and content
	if inp.Title == "" || inp.Content == "" {
		return models.Note{}, errors.Error("Please enter a non-empty value in Title and Content")
	}

	// inserting input note into db
	queryInsert := "INSERT INTO notes (note_id,title,content) VALUES (?,?,?)"

	//Used to execute insert query
	result, err := ctx.DB().ExecContext(ctx, queryInsert, inp.ID, inp.Title, inp.Content)

	if err != nil {
		return models.Note{}, errors.DB{Err: err}
	}

	fmt.Println("Data Inserted Successfully")

	lastInsertId, err := result.LastInsertId()

	if err != nil {
		return models.Note{}, errors.DB{Err: err}
	}

	// Now using SELECT query to find the inserted row
	queryFind := "SELECT note_id,title,content FROM notes WHERE note_id=?"

	_ = ctx.DB().QueryRowContext(ctx, queryFind, lastInsertId).Scan(&res.ID, &res.Title, &res.Content)

	// Handle the error if any
	if err != nil {
		return models.Note{}, errors.DB{Err: err}
	}

	fmt.Println("Response Generated Successfully")

	return res, nil

}

func (n note) Update(ctx *gofr.Context, id int, inp models.Note) (models.Note, error) {
	var res models.Note
	//checking as ID should be positive
	if inp.ID < 0 {
		return models.Note{}, errors.Error("Please enter a positive ID")
	}

	//checking as there is should be no empty value in title and content
	if inp.Title == "" || inp.Content == "" {
		return models.Note{}, errors.Error("Please enter a non-empty value in Title and Content")
	}

	//query to update the note with the id provided in the path
	queryUpdate := "UPDATE notes SET note_id=?,title = ?, content = ? WHERE note_id = ?"

	//used to run update query

	_, err := ctx.DB().ExecContext(ctx, queryUpdate, inp.ID, inp.Title, inp.Content, id)

	if err != nil {
		return models.Note{}, errors.DB{Err: err}
	}

	fmt.Println("Data Updated Successfully")

	// Now using SELECT query to find the updated row

	queryFind := "SELECT note_id,title,content FROM notes WHERE note_id=?"

	err = ctx.DB().QueryRowContext(ctx, queryFind, inp.ID).Scan(&res.ID, &res.Title, &res.Content)

	// Handle the error if any
	if err != nil {
		return models.Note{}, errors.DB{Err: err}
	}

	fmt.Println("Response Generated Successfully")

	return res, nil
}

func (n note) Delete(ctx *gofr.Context, id int) (models.Note, error) {
	var res models.Note

	// first using SELECT query to find the row to be deleted
	queryFind := "SELECT note_id,title,content FROM notes WHERE note_id=?"

	var error = ctx.DB().QueryRowContext(ctx, queryFind, id).Scan(&res.ID, &res.Title, &res.Content)

	if error != nil {
		return models.Note{}, errors.DB{Err: error}
	}

	fmt.Println("Response Generated Successfully")

	//query to update the note with the id provided in the path

	queryDelete := "DELETE FROM notes  WHERE note_id = ?"

	//used to run delete query

	_, err := ctx.DB().ExecContext(ctx, queryDelete, id)

	if err != nil {
		return models.Note{}, errors.DB{Err: err}
	}

	fmt.Println("Data Deleted Successfully")

	return res, nil
}
