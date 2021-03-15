package model

type UseCasesInterface interface {
	Register(request LoginRequest) (*LoginResponse, error)
	Login(request LoginRequest) (*LoginResponse, error)
	Logout(token Token) error

	CreateDoc(token Token, doc Doc) (*Doc, error)
	GetDoc(token Token, docId Id) (*Doc, error)
	EditDoc(token Token, newDoc Doc) (*Doc, error)
	DeleteDoc(token Token, docId Id) error
	ChangeDocAccess(token Token, request DocAccessRequest) error
	GetAllDocs(token Token) ([]Doc, error)

	GetUserById(userId Id) (*User, error)
	GetFriends(token Token) ([]User, error)

	CreateGroup(token Token, group Group) (*Group, error)
	DeleteGroup(token Token, groupId Id) error
	AddMember(token Token, groupId Id, newMemberId Id) error
	RemoveMember(token Token, groupId Id, memberId Id) error
	GetMembers(token Token, request GroupMembersChunkRequest) ([]User, error)
}

type Id string

type Token string

type User struct {
	Id    Id `json:"id"`
	Login string `json:"login"`
}

type Doc struct {
	Id       Id `json:"id"`
	AuthorId Id `json:"authorId"`
	Text     string `json:"text"`
	Access   string `json:"access"`
}

type Group struct {
	Id   Id `json:"id"`
	Name string `json:"name"`
}

type DocAccessRequest struct {
	DocId  Id `json:"id"`
	Access string `json:"access"`
}

type GroupMembersChunkRequest struct {
	Id    Id `json:"id"`
	Begin int `json:"begin"`
	Size  int `json:"end"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token Token
	User  User
}



