package utils

import (
	"encoding/json"
	"os"

	"github.com/ShuM01/test2/models"
)

func ReadJSON(filename string) ([]models.Feedback, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var feedbacks []models.Feedback
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&feedbacks)
	return feedbacks, err
}

func WriteJSON(filename string, feedbacks []models.Feedback) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(feedbacks)
}
