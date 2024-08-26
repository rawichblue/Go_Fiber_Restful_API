package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var users = []User{
	{ID: 1, Name: "John"},
	{ID: 2, Name: "Alice"},
	{ID: 3, Name: "Bob"},
}

// GetList
func getUsers(c *fiber.Ctx) error {
	name := c.Query("name")
	var filteredUsers []User

	if name != "" {
		for _, user := range users {
			if user.Name == name {
				filteredUsers = append(filteredUsers, user)
			}
		}
		if len(filteredUsers) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"code":    fiber.StatusNotFound,
				"message": "No users found with the given name",
				"data":    nil,
			})
		}
		return c.JSON(fiber.Map{
			"status": "success",
			"code":   fiber.StatusOK,
			"data":   filteredUsers,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"code":   fiber.StatusOK,
		"data":   users,
	})
}

// GetByID
func getUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"code":    fiber.StatusBadRequest,
			"message": "Invalid user ID",
		})
	}

	for _, user := range users {
		if user.ID == id {
			return c.JSON(fiber.Map{
				"status": "success",
				"code":   fiber.StatusOK,
				"data":   user,
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"code":    fiber.StatusNotFound,
		"message": "User not found",
	})
}

// Create
func createUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"code":    fiber.StatusBadRequest,
			"message": "Cannot parse JSON",
		})
	}

	user.ID = len(users) + 1
	users = append(users, *user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"code":   fiber.StatusCreated,
		"data":   user,
	})
}

// Update
func updateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"code":    fiber.StatusBadRequest,
			"message": "Invalid user ID",
		})
	}

	for i, user := range users {
		if user.ID == id {
			if err := c.BodyParser(&users[i]); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"code":    fiber.StatusBadRequest,
					"message": "Cannot parse JSON",
				})
			}
			users[i].ID = id
			return c.JSON(fiber.Map{
				"status": "success",
				"code":   fiber.StatusOK,
				"data":   users[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"code":    fiber.StatusNotFound,
		"message": "User not found",
	})
}

// Delete
func deleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"code":    fiber.StatusBadRequest,
			"message": "Invalid user ID",
		})
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status":  "success",
				"code":    fiber.StatusNoContent,
				"message": "User deleted",
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"code":    fiber.StatusNotFound,
		"message": "User not found",
	})
}

func main() {
	app := fiber.New()

	app.Get("/users", getUsers)
	app.Get("/users/:id", getUser)
	app.Post("/users", createUser)
	app.Put("/users/:id", updateUser)
	app.Delete("/users/:id", deleteUser)

	app.Listen(":3000")
}
