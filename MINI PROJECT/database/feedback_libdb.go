package database

import (
	"database/sql"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type Feedback struct {
	ID     int    `json:"id"`
	Kritik string `json:"kritik"`
	Saran  string `json:"saran"`
}

func GetAllFeedback(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		rows, err := db.Query("SELECT id, kritik, saran FROM feedback")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching feedback"})
		}
		defer rows.Close()

		feedbacks := make([]Feedback, 0)
		for rows.Next() {
			feedback := Feedback{}
			err := rows.Scan(&feedback.ID, &feedback.Kritik, &feedback.Saran)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching feedback"})
			}
			feedbacks = append(feedbacks, feedback)
		}
		return c.JSON(http.StatusOK, feedbacks)
	}
}

func GetFeedback(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid feedback ID"})
		}

		row := db.QueryRow("SELECT id, kritik, saran FROM feedback WHERE id=?", id)
		feedback := Feedback{}
		err = row.Scan(&feedback.ID, &feedback.Kritik, &feedback.Saran)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "Feedback not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching feedback"})
		}

		return c.JSON(http.StatusOK, feedback)
	}
}

func AddFeedback(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		feedback := Feedback{}
		err := c.Bind(&feedback)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		}

		result, err := db.Exec("INSERT INTO feedback (kritik, saran) VALUES (?, ?)", feedback.Kritik, feedback.Saran)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating feedback"})
		}

		id, err := result.LastInsertId()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating feedback"})
		}
		feedback.ID = int(id)

		return c.JSON(http.StatusCreated, feedback)
	}
}

func UpdateFeedback(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid feedback ID"})
		}

		feedback := Feedback{}
		err = c.Bind(&feedback)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		}

		_, err = db.Exec("UPDATE feedback SET kritik=?, saran=? WHERE id=?", feedback.Kritik, feedback.Saran, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error updating feedback"})
		}

		feedback.ID = id

		return c.JSON(http.StatusOK, feedback)
	}
}

func DeleteFeedback(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid feedback ID"})
		}
		_, err = db.Exec("DELETE FROM feedback WHERE id=?", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error deleting feedback"})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
