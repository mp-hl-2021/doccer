package model

import "doccer/data"

type Storage interface {
	GetUser(userId data.Id) (*data.User, error)
	GetUserByLogin(login string) (*data.User, error)
	GetHashedPassword(userId data.Id) (*Password, error)
    AddUser(newUser data.User, password Password) error
	EditUser(newUser data.User) (*data.User, error)
	CheckLoginExists(login string) bool

	CheckAccess(userId data.Id, docId data.Id) (string, error)
	GetDoc(docId data.Id) (*data.Doc, error)
	AddDoc(newDoc data.Doc) (*data.Id, error)
	EditDoc(newDoc data.Doc) (*data.Doc, error)
	EditDocAccess(docId data.Id, request DocAccessRequest) error
	DeleteDoc(docId data.Id) error
	GetAllDocs(userId data.Id) ([]data.Doc, error)

	CreateGroup(group data.Group) (*data.Group, error)
	DeleteGroup(groupId data.Id) error
	EditGroup(newGroup data.Group) (*data.Group, error)
	GetGroupById(groupId data.Id) (*data.Group, error)

	AddMember(groupId data.Id, newMemberId data.Id) error
	RemoveMember(groupId data.Id, memberId data.Id) error
	GetMembers(request GroupMembersChunkRequest) ([]data.User, error)

	GenerateNewUserId() data.Id
	GenerateNewDocId() data.Id
	GenerateNewGroupId() data.Id
}
