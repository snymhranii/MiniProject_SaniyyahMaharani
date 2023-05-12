package database

import (
	"database/sql"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type Visitor struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Review  string `json:"review"`
}

func GetVisitors(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var visitors []Visitor

		rows, err := db.Query("SELECT * FROM visitors")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		defer rows.Close()

		for rows.Next() {
			var visitor Visitor
			err := rows.Scan(&visitor.ID, &visitor.Name, &visitor.Age, &visitor.Address, &visitor.Review)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			visitors = append(visitors, visitor)
		}

		return c.JSON(http.StatusOK, visitors)
	}
}

func GetVisitor(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid visitor ID"})
		}

		var visitor Visitor
		err = db.QueryRow("SELECT * FROM visitors WHERE id = ?", id).Scan(&visitor.ID, &visitor.Name, &visitor.Age, &visitor.Address, &visitor.Review)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "visitor not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, visitor)
	}
}

func CreateVisitor(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var visitor Visitor
		if err := c.Bind(&visitor); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		}

		result, err := db.Exec("INSERT INTO visitors (name, age, address, review) VALUES (?, ?, ?, ?)", visitor.Name, visitor.Age, visitor.Address, visitor.Review)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		id, _ := result.LastInsertId()
		visitor.ID = int(id)

		return c.JSON(http.StatusCreated, visitor)
	}
}

func UpdateVisitor(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid visitor ID"})
		}

		var visitor Visitor
		if err := c.Bind(&visitor); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		}
		visitor.ID = id

		_, err = db.Exec("UPDATE visitors SET name = ?, age = ?, address = ?, review = ? WHERE id = ?", visitor.Name, visitor.Age, visitor.Address, visitor.Review, visitor.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, visitor)
	}
}

func DeleteVisitor(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid visitor ID"})
		}
		_, err = db.Exec("DELETE FROM visitors WHERE id = ?", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
