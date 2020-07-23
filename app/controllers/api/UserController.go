package api

import (
	"fmt"

	"github.com/gofiber/fiber"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/mikeychowy/fiber-crayplate/database"
)

type userData struct {
	UserId int
	Name   string
}

type responseStruct struct {
	Success bool
	Status  int
	Message string
	Data    []userData
}

// GetAllUsers Return all users as JSON
func GetAllUsers(c *fiber.Ctx) {

	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)

	dataSlice := make([]userData, 0, 1000)

	db := database.Instance()
	rows, err := db.Query(c.Context(), "SELECT * FROM users")

	if err != nil {
		panic(fmt.Sprintf("Error returning all users from database: %s", err))
	}

	defer rows.Close()

	if rows.Err() != nil {
		panic(fmt.Sprintf("Error reading all users rows from database: %s", rows.Err()))
	}

	rs := responseStruct{
		Success: true,
		Status:  200,
		Message: "Here is all the users we have",
		Data:    dataSlice,
	}
	ud := userData{
		UserId: 0,
		Name:   "",
	}
	var id int
	var name string

	c.Type("json")

	for rows.Next() {
		errR := rows.Scan(&id, &name)

		if errR != nil {
			panic(fmt.Sprintf("Error Scanning result set of users: %s", errR))
		}

		*&ud.UserId, *&ud.Name = id, name

		*&dataSlice = append(*&dataSlice, *&ud)
	}

	if len(*&dataSlice) <= 0 {

		*&rs.Success, *&rs.Status, *&rs.Message, *&rs.Data = false, 404, "We can't find any users, create some first", *&dataSlice
		c.Status(*&rs.Status)

		output, errJ := jsoniter.Marshal(*&rs)
		if errJ != nil {
			panic(fmt.Sprintf("Error converting all users to json: %s", errJ))
		}
		c.Send(output)
	} else {

		*&rs.Success, *&rs.Status, *&rs.Message, *&rs.Data = true, 200, "Here are all the users", *&dataSlice
		c.Status(*&rs.Status)

		output, errJ := jsoniter.Marshal(*&rs)
		if errJ != nil {
			panic(fmt.Sprintf("Error converting all users to json: %s", errJ))
		}
		c.Send(output)
	}

}

// GetUser Return a single user as JSON
func GetUser(c *fiber.Ctx) {
	// db := database.Instance()
	c.JSON("a")
	panic("a")
	// if err != nil {
	// 	panic(fmt.Sprintf("Error returning a user as JSON: %s", err))
	// }
}

// AddUser Add a single user to the database
func AddUser(c *fiber.Ctx) {
	// db := database.Instance()
	// User := new(models.User)
	// if err := c.BodyParser(User); err != nil {
	// 	c.Send("An error occurred when parsing the new user", err)
	// }
	// if res := db.Create(&User); res.Error != nil {
	// 	c.Send("An error occurred when storing the new user", res.Error)
	// }
	// // Match role to user
	// if User.RoleID != 0 {
	// 	Role := new(models.Role)
	// 	if res := db.Find(&Role, User.RoleID); res.Error != nil {
	// 		c.Send("An error occurred when retrieving the role", res.Error)
	// 	}
	// 	if Role.ID != 0 {
	// 		User.Role = *Role
	// 	}
	// }
	err := c.JSON("b")
	if err != nil {
		panic("Error occurred when returning JSON of a user")
	}
}

// EditUser Edit a single user
func EditUser(c *fiber.Ctx) {
	// db := database.Instance()
	// id := c.Params("id")
	// EditUser := new(models.User)
	// User := new(models.User)
	// if err := c.BodyParser(EditUser); err != nil {
	// 	c.Send("An error occurred when parsing the edited user", err)
	// }
	// if res := db.Find(&User, id); res.Error != nil {
	// 	c.Send("An error occurred when retrieving the existing user", res.Error)
	// }
	// User does not exist
	// if User.ID == 0 {
	// 	c.SendStatus(404)
	// 	err := c.JSON(fiber.Map{
	// 		"ID": id,
	// 	})
	// 	if err != nil {
	// 		panic("Error occurred when returning JSON of a user")
	// 	}
	// 	return
	// }
	// User.Name = EditUser.Name
	// User.Email = EditUser.Email
	// User.RoleID = EditUser.RoleID
	// Match role to user
	// if User.RoleID != 0 {
	// 	Role := new(models.Role)
	// 	if res := db.Find(&Role, User.RoleID); res.Error != nil {
	// 		c.Send("An error occurred when retrieving the role", res.Error)
	// 	}
	// 	if Role.ID != 0 {
	// 		User.Role = *Role
	// 	}
	// }
	// // Save user
	// db.Save(&User)

	err := c.JSON("c")
	if err != nil {
		panic("Error occurred when returning JSON of a user")
	}
}

// DeleteUser Delete a single user
func DeleteUser(c *fiber.Ctx) {
	// id := c.Params("id")
	// db := database.Instance()

	// var User models.User
	// db.Find(&User, id)
	// if res := db.Find(&User); res.Error != nil {
	// 	c.Send("An error occurred when finding the user to be deleted", res.Error)
	// }
	// db.Delete(&User)

	err := c.JSON(fiber.Map{
		"ID":      1,
		"Deleted": true,
	})
	if err != nil {
		panic("Error occurred when returning JSON of a user")
	}
}
