package controller

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
	Gender  string `json:"gender"`
	Review  string `json:"review"`
}

func getAllVisitors(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			visitor  Visitor
			visitors []Visitor
		)

		rows, err := db.Query("SELECT * FROM visitors")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&visitor.ID, &visitor.Name, &visitor.Age, &visitor.Address, &visitor.Gender, &visitor.Review)
			if err != nil {
				return err
			}
			visitors = append(visitors, visitor)
		}

		return c.JSON(http.StatusOK, visitors)
	}
}

func getVisitor(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		var visitor Visitor

		err = db.QueryRow("SELECT * FROM visitors WHERE id=?", id).Scan(&visitor.ID, &visitor.Name, &visitor.Age, &visitor.Address, &visitor.Gender, &visitor.Review)
		if err != nil {
			if err == sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound, "Visitor not found")
			}
			return err
		}

		return c.JSON(http.StatusOK, visitor)
	}
}

func createVisitor(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var visitor Visitor

		if err := c.Bind(&visitor); err != nil {
			return err
		}

		result, err := db.Exec("INSERT INTO visitors(name, age, address, gender, review) VALUES (?, ?, ?, ?, ?)", visitor.Name, visitor.Age, visitor.Address, visitor.Gender, visitor.Review)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}

		visitor.ID = int(id)

		return c.JSON(http.StatusCreated, visitor)
	}
}

func updateVisitor(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var visitor Visitor

		if err := c.Bind(&visitor); err != nil {
			return err
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		result, err := db.Exec("UPDATE visitors SET name=?, age=?, address=?, gender=?, review=? WHERE id=?", visitor.Name, visitor.Age, visitor.Address, visitor.Gender, visitor.Review, id)
		if err != nil {
			return err
		}

		rowsUpdated, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsUpdated == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Visitor not found")
		}

		return c.JSON(http.StatusOK, visitor)
	}
}

func deleteVisitor(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		result, err := db.Exec("DELETE FROM visitors WHERE id=?", id)
		if err != nil {
			return err
		}

		rowsDeleted, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsDeleted == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Visitor not found")
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func NewVisitorController(db *sql.DB) *echo.Echo {
	e := echo.New()

	e.GET("/visitors", getAllVisitors(db))
	e.GET("/visitors/:id", getVisitor(db))
	e.POST("/visitors", createVisitor(db))
	e.PUT("/visitors/:id", updateVisitor(db))
	e.DELETE("/visitors/:id", deleteVisitor(db))

	return e
}
