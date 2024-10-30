package dto

type UserCreate struct {
	Username string
	Password string
	Role     string
}

type PersonCreate struct {
	ID      string   `json:"personId"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Hobbies []string `json:"hobbies"`
}
type PersonUpdate struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Hobbies []string `json:"hobbies"`
}
