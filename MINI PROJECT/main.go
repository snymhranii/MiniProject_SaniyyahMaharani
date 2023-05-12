package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type Visitor struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Address string `json:"address"`
	Review  string `json:"review"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Feedback struct {
	ID     int    `json:"id"`
	Kritik string `json:"kritik"`
	Saran  string `json:"saran"`
}

type JwtClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func initMigrate() error {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}

func main() {
	e := echo.New()

	err := initMigrate()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/visitors", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to connect to database",
			})
		}
		defer db.Close()

		rows, err := db.Query("SELECT id, name, age, gender, address, review FROM data_kunjungan")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to query database",
			})
		}
		defer rows.Close()

		var visitors []Visitor
		for rows.Next() {
			var v Visitor
			err = rows.Scan(&v.Id, &v.Name, &v.Age, &v.Gender, &v.Address, &v.Review)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to parse visitors",
				})
			}
			visitors = append(visitors, v)
		}

		return c.JSON(http.StatusOK, visitors)
	})

	e.POST("/visitors", func(c echo.Context) error {
		var v Visitor
		err := c.Bind(&v)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			})
		}

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to connect to database",
			})
		}
		defer db.Close()

		result, err := db.Exec("INSERT INTO data_kunjungan (name, age, gender, address, review) VALUES (?, ?, ?, ?, ?)", v.Name, v.Age, v.Gender, v.Address, v.Review)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to insert visitor",
			})
		}

		lastId, _ := result.LastInsertId()
		v.Id = int(lastId)

		return c.JSON(http.StatusCreated, v)
	})

	e.PUT("/visitors/:id", func(c echo.Context) error {
		id := c.Param("id")

		var v Visitor
		err := c.Bind(&v)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			})
		}

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to connect to database",
			})
		}
		defer db.Close()

		result, err := db.Exec("UPDATE data_kunjungan SET name = ?, age = ?, gender = ?, address = ?, review = ? WHERE id = ?", v.Name, v.Age, v.Gender, v.Address, v.Review, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update visitor",
			})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Visitor not found",
			})
		}

		v.Id, _ = strconv.Atoi(id)

		return c.JSON(http.StatusOK, v)
	})

	e.DELETE("/visitors/:id", func(c echo.Context) error {
		id := c.Param("id")

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to connect to database",
			})
		}
		defer db.Close()

		result, err := db.Exec("DELETE FROM data_kunjungan WHERE id = ?", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to delete visitor",
			})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Visitor not found",
			})
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.POST("/users", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		user := new(User)
		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request payload",
			})
		}

		// insert user into database
		result, err := db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create user",
			})
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create user",
			})
		}

		return c.JSON(http.StatusCreated, map[string]int64{
			"id": lastInsertID,
		})
	})

	// get user by ID
	e.GET("/users/:id", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid user ID",
			})
		}

		// get user from database
		row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
		user := new(User)
		err = row.Scan(&user.Id, &user.Email, &user.Password)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get user",
			})
		}

		return c.JSON(http.StatusOK, user)
	})

	// update user by ID
	e.PUT("/users/:id", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid user ID",
			})
		}

		user := new(User)
		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request payload",
			})
		}

		// update user in database
		result, err := db.Exec("UPDATE users SET email = ?, password = ? WHERE id = ?", user.Email, user.Password, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update user",
			})
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update user",
			})
		} else if rowsAffected == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}

		return c.NoContent(http.StatusNoContent)
	})

	// delete user by ID
	e.DELETE("/users/:id", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid user ID",
			})
		}

		// delete user from database
		result, err := db.Exec("DELETE FROM users WHERE id = ?", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to delete user",
			})
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to delete user",
			})
		} else if rowsAffected == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}

		return c.NoContent(http.StatusNoContent)
	})

	// user login
	e.POST("/login", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		user := new(User)
		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request payload",
			})
		}

		// check if user exists in database
		row := db.QueryRow("SELECT * FROM users WHERE email = ? AND password = ?", user.Email, user.Password)
		err = row.Scan(&user.Id, &user.Email, &user.Password)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid email or password",
			})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get user",
			})
		}

		// create JWT token
		claims := &JwtClaims{
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // token expires in 24 hours
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create JWT token",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": tokenString,
		})
	})

	e.GET("/feedback", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		// Mengeksekusi query untuk mengambil semua feedback
		rows, err := db.Query("SELECT * FROM feedback")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get feedback",
			})
		}
		defer rows.Close()

		// Membuat slice untuk menyimpan feedback
		feedbacks := make([]Feedback, 0)

		// Looping untuk setiap baris hasil query
		for rows.Next() {
			// Membuat variabel untuk menyimpan data feedback dari baris tersebut
			var feedback Feedback

			// Mengambil data dari baris saat ini dan memasukkannya ke variabel feedback
			err := rows.Scan(&feedback.ID, &feedback.Kritik, &feedback.Saran)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to get feedback",
				})
			}

			// Menambahkan feedback ke slice feedbacks
			feedbacks = append(feedbacks, feedback)
		}

		// Mengembalikan response dengan slice feedbacks
		return c.JSON(http.StatusOK, feedbacks)
	})

	// Mengambil feedback berdasarkan ID
	e.GET("/feedback/:id", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		// Mengambil ID dari parameter URL
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid ID",
			})
		}

		// Mengeksekusi query untuk mengambil feedback dengan ID tersebut
		rows, err := db.Query("SELECT * FROM feedback WHERE id=?", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get feedback",
			})
		}
		defer rows.Close()

		// Mengambil data feedback dari baris hasil query
		if rows.Next() {
			var feedback Feedback
			err := rows.Scan(&feedback.ID, &feedback.Kritik, &feedback.Saran)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to get feedback",
				})
			}
			return c.JSON(http.StatusOK, feedback)
		}

		// Jika feedback dengan ID tersebut tidak ditemukan
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Feedback not found",
		})
	})

	// Menambahkan feedback baru
	e.POST("/feedback", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		// Membuat objek Feedback
		feedback := new(Feedback)

		// Membaca data JSON dari request body dan memasukkannya ke objek feedback
		if err := c.Bind(feedback); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			})
		}

		// Mengeksekusi query untuk menambahkan feedback baru
		result, err := db.Exec("INSERT INTO feedback (kritik, saran) VALUES (?, ?)", feedback.Kritik, feedback.Saran)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to add feedback",
			})
		}

		// Mengambil ID dari feedback yang baru ditambahkan
		id, err := result.LastInsertId()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to add feedback",
			})
		}

		// Memasukkan ID ke objek feedback dan mengembalikan response
		feedback.ID = int(id)
		return c.JSON(http.StatusCreated, feedback)
	})

	// Menghapus feedback berdasarkan ID
	e.DELETE("/feedback/:id", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kunjungan_museum")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		// Mengambil ID dari parameter URL
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid ID",
			})
		}
		// Mengeksekusi query untuk menghapus feedback dengan ID tersebut
		result, err := db.Exec("DELETE FROM feedback WHERE id=?", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to delete feedback",
			})
		}

		// Mengambil jumlah baris yang terpengaruh oleh query
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to delete feedback",
			})
		}

		// Jika feedback dengan ID tersebut ditemukan dan berhasil dihapus
		if rowsAffected > 0 {
			return c.NoContent(http.StatusNoContent)
		}

		// Jika feedback dengan ID tersebut tidak ditemukan
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Feedback not found",
		})
	})

	e.Logger.Fatal(e.Start(":8000"))

}
