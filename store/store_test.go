package store

import (
	"context"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/amannlalwani/Notes-app-using-gofr/models"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/datastore"
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
)

func TestCoreLayer(*testing.T) {
	app := gofr.New()

	// initializing the seeder
	seeder := datastore.NewSeeder(&app.DataStore, "../db")
	seeder.ResetCounter = true

	createTable(app)
}

func createTable(app *gofr.Gofr) {
	// drop table to clean previously added id's
	_, err := app.DB().Exec("DROP TABLE IF EXISTS notes;")

	if err != nil {
		return
	}

	_, err = app.DB().Exec("CREATE TABLE IF NOT EXISTS notes " +
		"( note_id INT AUTO_INCREMENT PRIMARY KEY, title VARCHAR(255) NOT NULL UNIQUE, content TEXT NOT NULL);")
	if err != nil {
		return
	}
}

func TestCreateNote(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		ctx.Logger.Error("mock connection failed")
	}

	ctx.DataStore = datastore.DataStore{ORM: db}
	ctx.Context = context.Background()

	tests := []struct {
		desc    string
		note    models.Note
		mockErr error
		err     error
	}{
		{"Valid case", models.Note{ID: 2, Title: "TestTitle123", Content: "Testing COntent 123"}, nil, nil},
		{"DB error", models.Note{ID: 5, Title: "TestTitle123", Content: "Tesing Content 456"}, errors.DB{}, errors.DB{Err: errors.DB{}}},
		{"empty value error", models.Note{ID: 7, Title: "", Content: "Testing Content 789"}, errors.Error("Please enter a non-empty value in Title and Content"), errors.Error("Please enter a non-empty value in Title and Content")},
		{"empty value error", models.Note{ID: 11, Title: "TestTitle456", Content: ""}, errors.Error("Please enter a non-empty value in Title and Content"), errors.Error("Please enter a non-empty value in Title and Content")},
		{"negative id error", models.Note{ID: -12, Title: "TestTitle789", Content: "Testing Content 000"}, errors.Error("Please enter a positive ID"), errors.Error("Please enter a positive ID")},
	}

	for index, testcase := range tests {
		mock.ExpectExec("INSERT INTO notes (note_id,title,content) VALUES (?,?,?)").WithArgs(testcase.note.ID, testcase.note.Title, testcase.note.Content).WillReturnResult(sqlmock.NewResult(2, 1)).WillReturnError(testcase.mockErr)

		data := sqlmock.NewRows([]string{"note_id", "title", "content"}).AddRow(testcase.note.ID, testcase.note.Title, testcase.note.Content)

		mock.ExpectQuery("SELECT note_id,title,content FROM notes WHERE note_id=?").WithArgs(testcase.note.ID).WillReturnRows(data).WillReturnError(testcase.mockErr)

		store := New()

		resp, err := store.Create(ctx, testcase.note)

		ctx.Logger.Log(resp)
		assert.IsType(t, testcase.err, err, "TEST[%d],failed.\n%s", index, testcase)
	}
}

func TestGetNotes(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		ctx.Logger.Error("mock connection failed")
	}

	ctx.DataStore = datastore.DataStore{ORM: db}
	ctx.Context = context.Background()

	tests := []struct {
		desc    string
		notes   []models.Note
		mockErr error
		err     error
	}{
		{"Valid case with notes", []models.Note{
			{ID: 1, Title: "TestTitle123", Content: "Testing Content 123"},
			{ID: 2, Title: "TestTitle456", Content: "Testing Content 456"},
		}, nil, nil},
		{"Valid case with no notes", []models.Note{}, nil, nil},
		{"Error case", nil, errors.Error("database error"), errors.DB{Err: errors.Error("database error")}},
	}

	for index, testcase := range tests {
		data := sqlmock.NewRows([]string{"note_id", "title", "content"})

		for _, note := range testcase.notes {
			data.AddRow(note.ID, note.Title, note.Content)
		}

		mock.ExpectQuery("SELECT note_id,title,content FROM notes").WillReturnRows(data).WillReturnError(testcase.mockErr)

		store := New()
		resp, err := store.Get(ctx)

		assert.Equal(t, testcase.err, err, "TEST[%d],failed.\n%s", index, testcase.desc)
		assert.Equal(t, testcase.notes, resp, "TEST[%d],failed.\n%s", index, testcase.desc)
	}

}

func TestUpdateNote(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		ctx.Logger.Error("mock connection failed")
	}

	ctx.DataStore = datastore.DataStore{ORM: db}
	ctx.Context = context.Background()

	tests := []struct {
		desc    string
		note    models.Note
		id      int
		mockErr error
		err     error
	}{
		{"Valid case", models.Note{ID: 2, Title: "UpdatedTestTitle123", Content: " Updated Testing COntent 123"}, 2, nil, nil},
		{"DB error", models.Note{ID: 2, Title: "UpdatedTestTitle123", Content: "Updated Tesing Content 456"}, 3, errors.DB{}, errors.DB{Err: errors.DB{}}},
		{"empty value error", models.Note{ID: 7, Title: "", Content: "Testing Content 789"}, 7, errors.Error("Please enter a non-empty value in Title and Content"), errors.Error("Please enter a non-empty value in Title and Content")},
		{"empty value error", models.Note{ID: 11, Title: "TestTitle456", Content: ""}, 11, errors.Error("Please enter a non-empty value in Title and Content"), errors.Error("Please enter a non-empty value in Title and Content")},
		{"negative id error", models.Note{ID: -12, Title: "TestTitle789", Content: "Testing Content 000"}, 12, errors.Error("Please enter a positive ID"), errors.Error("Please enter a positive ID")},
	}

	for index, testcase := range tests {
		mock.ExpectExec("UPDATE notes SET note_id=?,title = ?, content = ? WHERE note_id = ?").WithArgs(testcase.note.ID, testcase.note.Title, testcase.note.Content, testcase.id).WillReturnResult(sqlmock.NewResult(int64(testcase.id), 1)).WillReturnError(testcase.mockErr)

		data := sqlmock.NewRows([]string{"note_id", "title", "content"}).AddRow(testcase.note.ID, testcase.note.Title, testcase.note.Content)

		mock.ExpectQuery("SELECT note_id,title,content FROM notes WHERE note_id=?").WithArgs(testcase.note.ID).WillReturnRows(data).WillReturnError(testcase.mockErr)

		store := New()

		resp, err := store.Update(ctx, testcase.id, testcase.note)

		ctx.Logger.Log(resp)
		assert.IsType(t, testcase.err, err, "TEST[%d],failed.\n%s", index, testcase)
	}
}

func TestDeleteNote(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		ctx.Logger.Error("mock connection failed")
	}

	ctx.DataStore = datastore.DataStore{ORM: db}
	ctx.Context = context.Background()

	tests := []struct {
		desc    string
		note    models.Note
		id      int
		mockErr error
		err     error
	}{
		{"Valid case", models.Note{ID: 2, Title: "Deleted TestTitle123", Content: "Deleted Testing Content 123"}, 2, nil, nil},
		{"DB error", models.Note{ID: 2, Title: "Deleted TestTitle123", Content: "Deleted  Tesing Content 456"}, 3, errors.DB{}, errors.DB{Err: errors.DB{}}},
	}

	for index, testcase := range tests {

		data := sqlmock.NewRows([]string{"note_id", "title", "content"}).AddRow(testcase.note.ID, testcase.note.Title, testcase.note.Content)

		mock.ExpectQuery("SELECT note_id,title,content FROM notes WHERE note_id=?").WithArgs(testcase.note.ID).WillReturnRows(data).WillReturnError(testcase.mockErr)

		mock.ExpectExec("DELETE FROM notes  WHERE note_id = ?").WithArgs(testcase.id).WillReturnResult(sqlmock.NewResult(int64(testcase.id), 1)).WillReturnError(testcase.mockErr)

		store := New()

		resp, err := store.Delete(ctx, testcase.id)

		ctx.Logger.Log(resp)
		assert.IsType(t, testcase.err, err, "TEST[%d],failed.\n%s", index, testcase)
	}
}
