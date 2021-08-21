package resources

type User struct {
	Email string `schema:"email"`
	Pass  string `schema:"pass"`
}

type AllUsers struct {
	Users []User
}
