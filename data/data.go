package data

type Id string

type User struct {
	Id    Id     `json:"id"`
	Login string `json:"login"`
}

type Doc struct {
	Id           Id     `json:"id"`
	AuthorId     Id     `json:"authorId"`
	Text         string `json:"text"`
	Access       string `json:"access"`
	Lang         string `json:"lang"`
	LinterStatus string `json:"lstatus"`
}

type Group struct {
	Id      Id     `json:"id"`
	Name    string `json:"name"`
	Creator Id     `json:"creator_id"`
}

