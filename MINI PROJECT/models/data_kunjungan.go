package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type Visitor struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Address string `json:"address"`
	Review  string `json:"review"`
}
