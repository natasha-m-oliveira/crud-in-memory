package api

type UserBody struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Biography string `json:"bio"`
}
