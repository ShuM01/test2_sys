package database

import (
	"database/sql"
	"encoding/json"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/ShuM01/test2/models"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("postgres", "postgres://feedback:test2@localhost/feedback_form?sslmode=disable")
	if err != nil {
		return err
	}

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func GetAllFeedbacks() ([]models.Feedback, error) {
	rows, err := db.Query("SELECT id, data FROM feedbacks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []models.Feedback
	for rows.Next() {
		var id int
		var data string
		err := rows.Scan(&id, &data)
		if err != nil {
			return nil, err
		}
		var f models.Feedback
		err = json.Unmarshal([]byte(data), &f)
		if err != nil {
			return nil, err
		}
		f.ID = id
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}

func GetFeedbackByID(id int) (models.Feedback, error) {
	var data string
	err := db.QueryRow("SELECT data FROM feedbacks WHERE id = ?", id).Scan(&data)
	if err != nil {
		return models.Feedback{}, err
	}
	var f models.Feedback
	err = json.Unmarshal([]byte(data), &f)
	if err != nil {
		return models.Feedback{}, err
	}
	f.ID = id
	return f, nil
}

func InsertFeedback(f models.Feedback) (int, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return 0, err
	}
	result, err := db.Exec("INSERT INTO feedbacks (data) VALUES (?)", string(data))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func UpdateFeedback(id int, f models.Feedback) error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE feedbacks SET data = ? WHERE id = ?", string(data), id)
	return err
}

func DeleteFeedback(id int) error {
	_, err := db.Exec("DELETE FROM feedbacks WHERE id = ?", id)
	return err
}
