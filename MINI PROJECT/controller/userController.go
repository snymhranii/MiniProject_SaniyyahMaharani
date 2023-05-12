package controller

import (
	"database/sql"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserController struct {
	DB *sql.DB
}

func (c *UserController) CreateUser(ctx echo.Context) error {
	u := new(User)
	if err := ctx.Bind(u); err != nil {
		return err
	}
	result, err := c.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", u.Name, u.Email, u.Password)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	u.ID = int(id)
	return ctx.JSON(http.StatusCreated, u)
}

func (c *UserController) GetUser(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}
	row := c.DB.QueryRow("SELECT id, name, email, password FROM users WHERE id = ?", id)
	u := new(User)
	err = row.Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, u)
}

func (c *UserController) UpdateUser(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}
	u := new(User)
	if err := ctx.Bind(u); err != nil {
		return err
	}
	_, err = c.DB.Exec("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?", u.Name, u.Email, u.Password, id)
	if err != nil {
		return err
	}
	u.ID = id
	return ctx.JSON(http.StatusOK, u)
}

func (c *UserController) DeleteUser(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}
	_, err = c.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
