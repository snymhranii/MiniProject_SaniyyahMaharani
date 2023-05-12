package userdb

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

type UserDB struct {
	DB *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{
		DB: db,
	}
}

func (udb *UserDB) GetUserByID(id int) (*User, error) {
	var user User
	err := udb.DB.QueryRow("SELECT id, name, email, password FROM users WHERE id=?", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (udb *UserDB) CreateUser(u *User) error {
	_, err := udb.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", u.Name, u.Email, u.Password)
	return err
}

func (udb *UserDB) UpdateUser(u *User) error {
	_, err := udb.DB.Exec("UPDATE users SET name=?, email=?, password=? WHERE id=?", u.Name, u.Email, u.Password, u.ID)
	return err
}

func (udb *UserDB) DeleteUser(id int) error {
	_, err := udb.DB.Exec("DELETE FROM users WHERE id=?", id)
	return err
}

func (udb *UserDB) ListUsers() ([]*User, error) {
	rows, err := udb.DB.Query("SELECT id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func getUserIDFromParam(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}
	return id, nil
}

func (udb *UserDB) GetUserHandler(c echo.Context) error {
	id, err := getUserIDFromParam(c)
	if err != nil {
		return err
	}
	user, err := udb.GetUserByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, user)
}

func (udb *UserDB) CreateUserHandler(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return err
	}
	if err := udb.CreateUser(&user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}
	return c.JSON(http.StatusOK, user)
}

func (udb *UserDB) UpdateUserHandler(c echo.Context) error {
	id, err := getUserIDFromParam(c)
	if err != nil {
		return err
	}
	user, err := udb.GetUserByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	if err := c.Bind(user); err != nil {
		return err
	}
	if err := udb.UpdateUser(user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}
	return c.JSON(http.StatusOK, user)
}

func (udb *UserDB) DeleteUserHandler(c echo.Context) error {
	id, err := getUserIDFromParam(c)
	if err != nil {
		return err
	}
	if err := udb.DeleteUser(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}
	return c.NoContent(http.StatusOK)
}

func (udb *UserDB) ListUsersHandler(c echo.Context) error {
	users, err := udb.ListUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users")
	}
	return c.JSON(http.StatusOK, users)
}
