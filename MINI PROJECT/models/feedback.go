package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type Feedback struct {
	ID     int    `json:"id"`
	Kritik string `json:"kritik"`
	Saran  string `json:"saran"`
}
