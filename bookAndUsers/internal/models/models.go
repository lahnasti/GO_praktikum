package models

type User struct {
	UID      int    `json:"uId"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Book struct {
	BID    int    `json:"bId"`
	Title  string `json:"title"`
	Author string `json:"author"`
	UID    int    `json:"uId"`
}
