package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"userApiGo/model" //import user model

	"log"
	"net/http"
	"os"
	"strconv" // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected successfully")
	// return the connection
	return db
}

// Create new user
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	// set the header to content type

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty user of type model.User
	var user model.User

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	insertID := insertUser(user)

	// format for response
	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetUser will return a single user by its id
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert id to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// get one user
	user, err := getUser(int64(id))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(user)
}

// Get all users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get all users from db
	allUsers, err := getAllUsers()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	json.NewEncoder(w).Encode(allUsers)
}

// Update one user
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the id from the request (id like a key)
	params := mux.Vars(r)

	// convert string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// Create user. Use type model.User
	var user model.User

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// update user
	updatedRows := updateUser(int64(id), user)

	// format the messge string
	msg := fmt.Sprintf("Updated successfully. Not have %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send respnse
	json.NewEncoder(w).Encode(res)
}

// Delete user from db
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	// convert id to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// format in string
	msg := fmt.Sprintf("Updated successfully")

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func insertUser(user model.User) int64 {

	// Open connect
	db := createConnection()

	// Close connect
	defer db.Close()

	var id int64

	err := db.QueryRow(`INSERT INTO users (name, lastname, age, birthdate) VALUES ($1, $2, $3, $4) RETURNING id`, user.Name, user.Lastname, user.Age, user.Birthdate).Scan(&id)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Put user %v", id)

	// return new id
	return id
}

// get one user from db
func getUser(id int64) (model.User, error) {
	// Open connect
	db := createConnection()

	//Close connect
	defer db.Close()

	// create user
	var user model.User

	// execute sql
	row := db.QueryRow(`SELECT * FROM users WHERE id=$1`, id)

	err := row.Scan(&user.ID, &user.Name, &user.Lastname, &user.Age, &user.Birthdate)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("Have no raws")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Error %v", err)
	}

	// return empty user on error
	return user, err
}

//get all users
func getAllUsers() ([]model.User, error) {
	// Open connect
	db := createConnection()

	// Close connect
	defer db.Close()

	var users []model.User

	// Create the select sql query
	sqlStatement := `SELECT * FROM users`

	// Execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Close the statement
	defer rows.Close()

	// Iterate rows
	for rows.Next() {
		var user model.User

		err = rows.Scan(&user.ID, &user.Name, &user.Lastname, &user.Age, &user.Birthdate)

		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// append one user in our users
		users = append(users, user)

	}

	// return empty user on error
	return users, err
}

// update user in db
func updateUser(id int64, user model.User) int64 {

	// Open connect
	db := createConnection()

	// Close connect
	defer db.Close()

	fmt.Println(user.Name)
	fmt.Println(user.Lastname)
	fmt.Println(user.Age)
	fmt.Println(user.Birthdate)

	// execute sql
	res, err := db.Exec(`UPDATE users SET name=$2, lastname=$3, age=$4, birthdate =$5 WHERE id=$1`, id, user.Name, user.Lastname, user.Age, user.Birthdate)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Total user %v", rowsAffected)

	return rowsAffected
}

// delete user
func deleteUser(id int64) int64 {

	// Open connect
	db := createConnection()

	// Close connect
	defer db.Close()

	// execute the sql statement
	res, err := db.Exec(`DELETE FROM users WHERE id=$1`, id)

	if err != nil {
		log.Fatalf("Cant delete. Error: %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Error: %v", rowsAffected)

	return rowsAffected
}
