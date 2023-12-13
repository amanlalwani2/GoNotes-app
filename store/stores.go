package store

import (
	"fmt"

	"github.com/amannlalwani/Notes-app-using-gofr/models"
	"gofr.dev/pkg/gofr"
)

type note struct{}

func New() Store {
	return note{}
}

func (n note) Get(ctx *gofr.Context) ([]models.Note, error) {
	data, err := ctx.DB().QueryContext(ctx, "SELECT note_id,title,content FROM notes ")
	if err != nil {
		panic(err)
	}

	defer data.Close()

	notes := make([]models.Note, 0)

	//making slice of notes data retrieved from DB.
	for data.Next() {
		var n models.Note

		err = data.Scan(&n.ID, &n.Title, &n.Content)

		if err != nil {
			panic(err)
		}

		notes = append(notes, n)
	}

	err = data.Err()

	if err != nil {
		panic(err)
	}

	return notes, nil
}

func (n note) Create(ctx *gofr.Context, inp models.Note) (models.Note, error) {
	var res models.Note

	queryInsert := "INSERT INTO notes (title,content) VALUES (?,?)"

	//Used to execute insert query
	result, err := ctx.DB().ExecContext(ctx, queryInsert, inp.Title, inp.Content)

	if err != nil {

		panic(err)
	}

	fmt.Println("Data Inserted Successfully")

	lastInsertId, err := result.LastInsertId()

	if err != nil {

		panic(err)
	}

	queryFind := "SELECT note_id,title,content,created_at,updated_at FROM notes WHERE note_id=?"

	_ = ctx.DB().QueryRowContext(ctx, queryFind, lastInsertId).Scan(&res.ID, &res.Title, &res.Content)

	return res, nil

}

func (n note) Update(ctx *gofr.Context, id int, inp models.Note) (models.Note, error) {
	var res models.Note

	queryUpdate := "UPDATE notes SET title = ?, content = ? WHERE note_id = ?"

	//used to run update query

	_, err := ctx.DB().ExecContext(ctx, queryUpdate, inp.Title, inp.Content, id)

	if err != nil {
		panic(err)
	}

	fmt.Println("Updated Successfully")

	queryFind := "SELECT note_id,title,content,created_at,updated_at FROM notes WHERE note_id=?"

	_ = ctx.DB().QueryRowContext(ctx, queryFind, id).Scan(&res.ID, &res.Title, &res.Content)

	fmt.Println(res)

	return res, nil
}

func (n note) Delete(ctx *gofr.Context, id int) (models.Note, error) {
	var res models.Note

	queryFind := "SELECT note_id,title,content,created_at,updated_at FROM notes WHERE note_id=?"

	_ = ctx.DB().QueryRowContext(ctx, queryFind, id).Scan(&res.ID, &res.Title, &res.Content)

	queryDelete := "DELETE FROM notes  WHERE note_id = ?"

	//used to run delete query

	_, err := ctx.DB().ExecContext(ctx, queryDelete, id)

	if err != nil {

		panic(err)
	}

	fmt.Println("Deleted Successfully")

	fmt.Println(res)

	return res, nil
}
