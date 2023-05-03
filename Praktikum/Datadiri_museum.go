package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Gender  string `json:"gender"`
}

var people []Person

// GET - get all people
func getAllPeople(c echo.Context) error {
	return c.JSON(http.StatusOK, people)
}

// GET - get a person by name
func getPerson(c echo.Context) error {
	name := c.Param("name")

	for _, p := range people {
		if p.Name == name {
			return c.JSON(http.StatusOK, p)
		}
	}

	return c.JSON(http.StatusNotFound, "Person not found")
}

// POST - add a person
func addPerson(c echo.Context) error {
	var p Person
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	people = append(people, p)
	return c.JSON(http.StatusCreated, "Person added successfully")
}

// PUT - update a person
// PUT - update a person
func updatePerson(c echo.Context) error {
	name := c.Param("name")
	var updatedPerson Person

	err := json.NewDecoder(c.Request().Body).Decode(&updatedPerson)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	for i, p := range people {
		if p.Name == name {
			people[i].Age = updatedPerson.Age
			people[i].Address = updatedPerson.Address
			people[i].Name = updatedPerson.Name
			people[i].Gender = updatedPerson.Gender
			return c.JSON(http.StatusOK, "Person updated successfully")
		}
	}

	return c.JSON(http.StatusNotFound, "Person not found")
}

// DELETE - delete a person
func deletePerson(c echo.Context) error {
	name := c.Param("name")

	for i, p := range people {
		if p.Name == name {
			people = append(people[:i], people[i+1:]...)
			return c.JSON(http.StatusOK, "Person deleted successfully")
		}
	}

	return c.JSON(http.StatusNotFound, "Person not found")
}

func main() {
	// Initialize Echo instance
	e := echo.New()

	// Seed initial data
	people = append(people, Person{Name: "sany", Age: 10, Address: "jalan irigasi", Gender: "Wanita"})

	// Define routes
	e.GET("/people", getAllPeople)
	e.GET("/people/:name", getPerson)
	e.POST("/people", addPerson)
	e.PUT("/people/:name", updatePerson)
	e.DELETE("/people/:name", deletePerson)

	// Start server
	e.Logger.Fatal(e.Start(":1045"))
}
