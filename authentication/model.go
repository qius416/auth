package authentication

// User model for post from client as well as for db mapping
type User struct {
	Email    string `json:"email" gorethink:"email"`
	Name     string `json:"name" gorethink:"name"`
	Password string `json:"password" gorethink:"password"`
	Role     string `gorethink:"role"`
}

// Auth is a model to store tokens for authentication
type Auth struct {
	Token string `json:"token"`
}
