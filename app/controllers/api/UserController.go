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

type requestBodyStruct struct {
	Name string `json:"name" xml:"name" form:"name"`
}

// GetAllUsers : Respond all users as JSON
func GetAllUsers(c *fiber.Ctx) {

	// naming strategy for jsoniter so we don't have to add json tags individually
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)

	// this is the slice to hold "data:[]" part of the response
	dataSlice := make([]userData, 0, 1000)

	// not exactly the db instance, it's the pool instance
	// but for all intents and purposes works exactly like db connection
	// here we query all rows from the query string and sort them by user id
	db := database.Instance()
	rows, err := db.Query(c.Context(), "SELECT * FROM users ORDER BY user_id ASC")

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
	// create the user data holder struct, to describe how i would like the json data insides to be like
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
	// create the user data holder struct, to describe how i would like the json data insides to be like
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

// AddUser : Add a single user to the database
func AddUser(c *fiber.Ctx) {
	rbod := new(requestBodyStruct)

	// parse the body and pass values to the struct
	if err := c.BodyParser(rbod); err != nil {
		panic(fmt.Sprintf("error parsing body in create new user: %s", err))
	}

	// this is the slice to hold "data:[]" part of the response
	dataSlice := make([]userData, 0, 1)

	// my standard struct for a response
	rs := responseStruct{
		Success: true,
		Status:  201,
		Message: "Here is all the users we have",
		Data:    dataSlice,
	}
	// create the user data holder struct, to describe how i would like the json data insides to be like
	ud := userData{
		UserId: 0,
		Name:   "",
	}

	// not exactly the db instance, it's the pool instance
	// but for all intents and purposes works exactly like db connection
	// here we execute the insert
	db := database.Instance()
	if _, err := db.Exec(c.Context(), "INSERT INTO users(name) VALUES($1)", rbod.Name); err != nil {
		panic(fmt.Sprintf("error inserting new user into database: %s", err))
	}

	// query the new user to double check
	if err := db.QueryRow(c.Context(), "SELECT * FROM users WHERE name=$1", rbod.Name).Scan(&ud.UserId, &ud.Name); err != nil {
		panic(fmt.Sprintf("error in getting the new user from database: %s", err))
	}

	// append the user data to the slice
	dataSlice = append(dataSlice, ud)

	// set the status, and header of content-type
	// also set the response message and data
	c.Status(rs.Status)
	c.Type("json")
	rs.Message, rs.Data = "Here is the new user", dataSlice

	// naming strategy for jsoniter so we don't have to add json tags individually
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)

	output, errJ := jsoniter.Marshal(rs)
	if errJ != nil {
		panic(fmt.Sprintf("Error converting the new user to json: %s", errJ))
	}

	c.Send(output)
}

// EditUser : Edit a single user
func EditUser(c *fiber.Ctx) {
	rbod := new(requestBodyStruct)

	// parse the body and pass values to the struct
	if err := c.BodyParser(rbod); err != nil {
		panic(fmt.Sprintf("error parsing body in edit user: %s", err))
	}

	// this is the slice to hold "data:[]" part of the response
	dataSlice := make([]userData, 0, 1)

	// my standard struct for a response
	rs := responseStruct{
		Success: true,
		Status:  200,
		Message: "Here is all the users we have",
		Data:    dataSlice,
	}
	// create the user data holder struct, to describe how i would like the json data insides to be like
	ud := userData{
		UserId: 0,
		Name:   "",
	}

	// get the request parameter of user id
	queryID := c.Params("id")

	// not exactly the db instance, it's the pool instance
	// but for all intents and purposes works exactly like db connection
	// here we execute the update
	db := database.Instance()
	if _, err := db.Exec(c.Context(), "UPDATE users SET name=$1 WHERE user_id=$2", rbod.Name, queryID); err != nil {
		panic(fmt.Sprintf("error updating specified user into database: %s", err))
	}

	// query the updated user to double check
	if err := db.QueryRow(c.Context(), "SELECT * FROM users WHERE user_id=$1", queryID).Scan(&ud.UserId, &ud.Name); err != nil {
		panic(fmt.Sprintf("error in getting the updated user from database: %s", err))
	}

	// append the user data to the slice
	dataSlice = append(dataSlice, ud)

	// set the status, and header of content-type
	// also set the response message and data
	c.Status(rs.Status)
	c.Type("json")
	rs.Message, rs.Data = "Here is the updated user", dataSlice

	// naming strategy for jsoniter so we don't have to add json tags individually
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)

	output, errJ := jsoniter.Marshal(rs)
	if errJ != nil {
		panic(fmt.Sprintf("Error converting the updated user to json: %s", errJ))
	}

	c.Send(output)
}

// DeleteUser : Delete a single user
func DeleteUser(c *fiber.Ctx) {

	// this is the slice to hold "data:[]" part of the response
	dataSlice := make([]userData, 0, 1)

	// my standard struct for a response
	rs := responseStruct{
		Success: true,
		Status:  202,
		Message: "Here is all the users we have",
		Data:    dataSlice,
	}

	// get the request parameter of user id
	queryID := c.Params("id")

	// not exactly the db instance, it's the pool instance
	// but for all intents and purposes works exactly like db connection
	// here we execute the update
	db := database.Instance()
	if _, err := db.Exec(c.Context(), "DELETE FROM users WHERE user_id=$1", queryID); err != nil {
		panic(fmt.Sprintf("error deleting specified user into database: %s", err))
	}

	// set the status, and header of content-type
	// also set the response message and data
	c.Status(rs.Status)
	c.Type("json")
	rs.Message = fmt.Sprintf("User %s successfuly deleted.", queryID)

	// naming strategy for jsoniter so we don't have to add json tags individually
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)

	output, errJ := jsoniter.Marshal(rs)
	if errJ != nil {
		panic(fmt.Sprintf("Error converting the deleted user to json: %s", errJ))
	}

	c.Send(output)
}
