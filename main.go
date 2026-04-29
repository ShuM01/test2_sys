package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Feedback struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

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

func GetAllFeedbacks() ([]Feedback, error) {
	rows, err := db.Query("SELECT id, data FROM feedbacks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []Feedback
	for rows.Next() {
		var id int
		var data string
		err := rows.Scan(&id, &data)
		if err != nil {
			return nil, err
		}
		var f Feedback
		err = json.Unmarshal([]byte(data), &f)
		if err != nil {
			return nil, err
		}
		f.ID = id
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}

func GetFeedbackByID(id int) (Feedback, error) {
	var data string
	err := db.QueryRow("SELECT data FROM feedbacks WHERE id = $1", id).Scan(&data)
	if err != nil {
		return Feedback{}, err
	}
	var f Feedback
	err = json.Unmarshal([]byte(data), &f)
	if err != nil {
		return Feedback{}, err
	}
	f.ID = id
	return f, nil
}

func InsertFeedback(f Feedback) (int, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return 0, err
	}
	var id int
	err = db.QueryRow("INSERT INTO feedbacks (data) VALUES ($1) RETURNING id", string(data)).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateFeedback(id int, f Feedback) error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE feedbacks SET data = $1 WHERE id = $2", string(data), id)
	return err
}

func DeleteFeedback(id int) error {
	_, err := db.Exec("DELETE FROM feedbacks WHERE id = $1", id)
	return err
}

func HandleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		feedbacks, err := GetAllFeedbacks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(feedbacks)
	case http.MethodPost:
		var f Feedback
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := InsertFeedback(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		f.ID = id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(f)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleItem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/feedback/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		f, err := GetFeedbackByID(id)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				http.NotFound(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(f)
	case http.MethodPut:
		var f Feedback
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.ID = id
		err := UpdateFeedback(id, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(f)
	case http.MethodDelete:
		err := DeleteFeedback(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleNames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	feedbacks, err := GetAllFeedbacks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	names := make([]string, len(feedbacks))
	for i, f := range feedbacks {
		names[i] = f.Name
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

func HandleEmails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	feedbacks, err := GetAllFeedbacks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	emails := make([]string, len(feedbacks))
	for i, f := range feedbacks {
		emails[i] = f.Email
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emails)
}

func HandleSubjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	feedbacks, err := GetAllFeedbacks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	subjects := make([]string, len(feedbacks))
	for i, f := range feedbacks {
		subjects[i] = f.Subject
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}

func HandleMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	feedbacks, err := GetAllFeedbacks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	messages := make([]string, len(feedbacks))
	for i, f := range feedbacks {
		messages[i] = f.Message
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func main() {
	err := InitDB()
	if err != nil {
		fmt.Printf("Error initializing DB: %v\n", err)
		return
	}

	// Serve static files from the "frontend" folder
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	// API Routes
	http.HandleFunc("/api/feedback", HandleCollection)
	http.HandleFunc("/api/feedback/", HandleItem)

	// Specific List Endpoints
	http.HandleFunc("/api/feedback/names", HandleNames)
	http.HandleFunc("/api/feedback/emails", HandleEmails)
	http.HandleFunc("/api/feedback/subjects", HandleSubjects)
	http.HandleFunc("/api/feedback/messages", HandleMessages)

	fmt.Println("system running at http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
