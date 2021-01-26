# userProjectGo

Postman service was used to send requests

Create new user (method POST)
URL: http://localhost:8080/api/newuser
Body: raw / json
Examples: 
{
 "Name":"myName", 
 "Lastname": "myLastname", 
 "Age": 20,
 "Birthdate: "01.10.2001"
} 

Get all users (method GET)
URL: http://localhost:8080/api/users

Get one user (method GET)
URL: http://localhost:8080/api/user/{id}


Delete a user (method DELETE)
URL: http://localhost:8080/api/deleteuser/{id}

Update a user (method PUT)
URL: http://localhost:8080/api/user/1
Body: raw/json
