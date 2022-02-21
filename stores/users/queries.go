package users

const (
	insertUser  = "INSERT INTO users(id,name,age) VALUES(?,?,?)"
	getUsers    = "SELECT * FROM users;"
	getUserByID = "SELECT * FROM users WHERE id=?;"
	updateUser  = "UPDATE users SET name=?,age=? WHERE id=?;"
	deleteUser  = "DELETE FROM users WHERE id=?;"
)
