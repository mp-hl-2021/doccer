package model

type UseCasesInterface interface {
	Register(request LoginRequest) error
	Login(request LoginRequest) (*LoginResponse, error)
	Auth(tokenStr string) (*string, error)
	Logout(token Token) error

	CreateDoc(userId Id, doc Doc) (*Doc, error)
	GetDoc(userId Id, docId Id) (*Doc, error)
	EditDoc(userId Id, newDoc Doc) (*Doc, error)
	DeleteDoc(userId Id, docId Id) error

	GetAllDocs(userId Id) ([]Doc, error)

	GetUserById(userId Id) (*User, error)
	EditUser(userId Id, newUser User) (*User, error)

	GetFriends(userId Id) ([]User, error)
	AddFriend(userId Id, friendId Id) error
	RemoveFriend(userId Id, friendId Id) error

	CreateGroup(userId Id, group Group) (*Group, error)
	DeleteGroup(userId Id, groupId Id) error
	EditGroup(userId Id, groupId Id, newGroup Group) (*Group, error)

	AddMember(userId Id, groupId Id, MemberId Id) error
	RemoveMember(userId Id, groupId Id, memberId Id) error
	GetMembers(userId Id, request GroupMembersChunkRequest) ([]User, error)
}

type Id string


type User struct {
	Id    Id `json:"id"`
	Login string `json:"login"`
}
type Token string

type Password []byte

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



