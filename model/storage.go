package model

type Storage interface {
	GetUser(userId Id) (*User, error)
	GetUserByLogin(login string) (*User, error)
	GetHashedPassword(userId Id) (*Password, error)
    AddUser(newUser User, password Password) error
	EditUser(newUser User) (*User, error)
	CheckLoginExists(login string) bool

	CheckAccess(userId Id, docId Id) (string, error)
	GetDoc(docId Id) (*Doc, error)
	AddDoc(newDoc Doc) (*Id, error)
	EditDoc(newDoc Doc) (*Doc, error)
	EditDocAccess(docId Id, request DocAccessRequest) error
	DeleteDoc(docId Id) error
	GetAllDocs(userId Id) ([]Doc, error)

	CreateGroup(group Group) (*Group, error)
	DeleteGroup(groupId Id) error
	EditGroup(newGroup Group) (*Group, error)
	GetGroupById(groupId Id) (*Group, error)

	AddMember(groupId Id, newMemberId Id) error
	RemoveMember(groupId Id, memberId Id) error
	GetMembers(request GroupMembersChunkRequest) ([]User, error)

	GenerateNewUserId() Id
	GenerateNewDocId() Id
	GenerateNewGroupId() Id
}
