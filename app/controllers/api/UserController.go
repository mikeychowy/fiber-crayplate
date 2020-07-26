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

// GetAllUsers : Respond all users as JSON
func GetAllUsers(c *fiber.Ctx) {

	// naming strategy for jsoniter so we don't have to add json tags individually
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)

	// this is the slice to hold "data:[]" part of the response
	dataSlice := make([]userData, 0, 1000)

	// not exactly the db instance, it's the pool instance
	// but for all intents and purposes works exactly like db connection
	// here we query all rows from the query
	db := database.Instance()
	rows, err := db.Query(c.Context(), "SELECT * FROM users")

	// pool error handling
	if err != nil {
		panic(fmt.Sprintf("Error returning all users from database: %s", err))
	}

	// so we don't forget to close the rows downstairs
	defer rows.Close()

	// rows error handler
	if rows.Err() != nil {
		panic(fmt.Sprintf("Error reading all users rows from database: %s", rows.Err()))
	}

	// my standard struct for a response
	rs := responseStruct{
		Success: true,
		Status:  200,
		Message: "Here is all the users we have",
		Data:    dataSlice,
	}
	// create the user data holder struct, to describe how i would like the json to be like
	ud := userData{
		UserId: 0,
		Name:   "",
	}
	// holders for ids and names
	var id int
	var name string

	// change response Content-Type header to application/json
	c.Type("json")

	// iterate through all rows returned from db
	for rows.Next() {
		// like i said, holders
		errR := rows.Scan(&id, &name)

		// scanner error handling
		if errR != nil {
			panic(fmt.Sprintf("Error Scanning result set of users: %s", errR))
		}

		// pass to holder struct
		ud.UserId, ud.Name = id, name

		// append to data slice
		dataSlice = append(dataSlice, ud)
	}

	if len(dataSlice) <= 0 {
		// check for 404

		rs.Success, rs.Status, rs.Message, rs.Data = false, 404, "We can't find any users, create some first", dataSlice
		c.Status(rs.Status)

		// marshal response struct to json
		output, errJ := jsoniter.Marshal(rs)
		if errJ != nil {
			panic(fmt.Sprintf("Error converting all users to json: %s", errJ))
		}
		c.Send(output)
	} else {
		// success is here

		rs.Success, rs.Status, rs.Message, rs.Data = true, 200, "Here are all the users", dataSlice
		c.Status(rs.Status)

		// marshal response struct to json
		output, errJ := jsoniter.Marshal(rs)
		if errJ != nil {
			panic(fmt.Sprintf("Error converting all users to json: %s", errJ))
		}
		c.Send(output)
	}
}

// GetUser : Respond a single user by id as JSON
func GetUser(c *fiber.Ctx) {

	// naming strategy for jsoniter so we don't have to add json tags individually
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)

	// this is the slice to hold "data:[]" part of the response
	dataSlice := make([]userData, 0, 10)

	// my standard struct for a response
	rs := responseStruct{
		Success: true,
		Status:  200,
		Message: "Here is all the users we have",
		Data:    dataSlice,
	}
	// create the user data holder struct, to describe how i would like the json to be like
	ud := userData{
		UserId: 0,
		Name:   "",
	}

	// get the request parameter of user id
	queryID := c.Params("id")

	// not exactly the db instance, it's the pool instance
	// but for all intents and purposes works exactly like db connection
	// here we query a row from the query
	db := database.Instance()
	err := db.QueryRow(c.Context(), "SELECT * FROM users WHERE user_id=$1", queryID).Scan(&ud.UserId, &ud.Name)

	// scanner error handling
	if err != nil {
		panic(fmt.Sprintf("Error returning user with id %s from the database: %s", queryID, err))
	}

	// append the user data to the slice
	dataSlice = append(dataSlice, ud)

	// set the status, and header of content-type
	// also set the response message and data
	c.Status(rs.Status)
	c.Type("json")
	rs.Message, rs.Data = "Here is the specified user", dataSlice

	output, errJ := jsoniter.Marshal(rs)
	if errJ != nil {
		panic(fmt.Sprintf("Error converting the specified user to json: %s", errJ))
	}

	c.Send(output)
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
