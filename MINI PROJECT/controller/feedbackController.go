package controller

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

func GetFeedback(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		feedbacks := []Feedback{}
		rows, err := db.Query("SELECT id, kritik, saran FROM feedback")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var feedback Feedback
			if err := rows.Scan(&feedback.ID, &feedback.Kritik, &feedback.Saran); err != nil {
				return err
			}
			feedbacks = append(feedbacks, feedback)
		}

		return c.JSON(http.StatusOK, feedbacks)
	}
}

func AddFeedback(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var feedback Feedback
		err := c.Bind(&feedback)
		if err != nil {
			return err
		}

		result, err := db.Exec("INSERT INTO feedback (kritik, saran) VALUES (?, ?)", feedback.Kritik, feedback.Saran)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		feedback.ID = int(id)

		return c.JSON(http.StatusCreated, feedback)
	}
}

func UpdateFeedback(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var feedback Feedback
		err := c.Bind(&feedback)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		result, err := db.Exec("UPDATE feedback SET kritik = ?, saran = ? WHERE id = ?", feedback.Kritik, feedback.Saran, id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Feedback not found")
		}

		return c.JSON(http.StatusOK, feedback)
	}
}

func DeleteFeedback(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		result, err := db.Exec("DELETE FROM feedback WHERE id = ?", id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Feedback not found")
		}

		return c.NoContent(http.StatusNoContent)
	}
}
