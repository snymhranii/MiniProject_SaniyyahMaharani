package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

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
}
