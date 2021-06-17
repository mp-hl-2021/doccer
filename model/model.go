package model

import "doccer/data"

type UseCasesInterface interface {
	Register(request LoginRequest) (*data.User, error)
	Login(request LoginRequest) (*LoginResponse, error)
	Auth(tokenStr string) (*string, error)
	Logout(token Token) error

	CreateDoc(userId data.Id, doc data.Doc) (*data.Doc, error)
	GetDoc(userId data.Id, docId data.Id) (*data.Doc, error)
	EditDoc(userId data.Id, newDoc data.Doc) (*data.Doc, error)
	DeleteDoc(userId data.Id, docId data.Id) error
	ChangeDocAccess(userId data.Id, request DocAccessRequest) (*data.Doc, error)

	GetAllDocs(userId data.Id) ([]data.Doc, error)

	GetUserById(userId data.Id) (*data.User, error)
	EditUser(userId data.Id, newUser data.User) (*data.User, error)

	CreateGroup(userId data.Id, group data.Group) (*data.Group, error)
	DeleteGroup(userId data.Id, groupId data.Id) error
	EditGroup(userId data.Id, newGroup data.Group) (*data.Group, error)

	AddMember(userId data.Id, groupId data.Id, MemberId data.Id) error
	RemoveMember(userId data.Id, groupId data.Id, memberId data.Id) error
	GetMembers(userId data.Id, request GroupMembersChunkRequest) ([]data.User, error)
}

type Token string

type Password []byte

type DocAccessRequest struct {
	DocId  data.Id `json:"id"`
	Type   int     `json:"type"`
	ItemId data.Id `json:"itemId"`
	Access string  `json:"access"`
}

type MemberRequest struct {
	GroupId  data.Id `json:"groupId"`
	MemberId data.Id `json:"memberId"`
}

type GroupMembersChunkRequest struct {
	Id    data.Id `json:"id"`
	Begin int     `json:"begin"`
	Size  int     `json:"end"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token Token
	User  data.User
}



