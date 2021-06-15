package model

type UseCasesInterface interface {
	Register(request LoginRequest) (*User, error)
	Login(request LoginRequest) (*LoginResponse, error)
	Auth(tokenStr string) (*string, error)
	Logout(token Token) error

	CreateDoc(userId Id, doc Doc) (*Doc, error)
	GetDoc(userId Id, docId Id) (*Doc, error)
	EditDoc(userId Id, newDoc Doc) (*Doc, error)
	DeleteDoc(userId Id, docId Id) error
	ChangeDocAccess(userId Id, request DocAccessRequest) (*Doc, error)

	GetAllDocs(userId Id) ([]Doc, error)

	GetUserById(userId Id) (*User, error)
	EditUser(userId Id, newUser User) (*User, error)

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
	Id      Id `json:"id"`
	Name    string `json:"name"`
	Creator Id `json:"creator_id"`
}

type DocAccessRequest struct {
	DocId  Id `json:"id"`
	Type int `json:"type"`
	ItemId Id `json:"itemId"`
	Access string `json:"access"`
}

type MemberRequest struct {
	GroupId  Id `json:"groupId"`
	MemberId Id `json:"memberId"`
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



