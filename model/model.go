package model

// User schema
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Lastname  string `json:"lastname"`
	Age       int64  `json:"age"`
	Birthdate string `json: "birthdate"`
}
