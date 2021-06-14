package model

type Storage interface {
	GetUser(userId Id) (*User, error)
	GetUserByLogin(login string) (*User, error)
	GetHashedPassword(userId Id) (*Password, error)
    AddUser(newUser User, password Password) error
	EditUser(userId Id, newUser User) (*User, error)
	GetUserById(userId Id) (*User, error)
	CheckLoginExists(login string) bool

	CheckAccess(userId Id, docId Id) (string, error)
	GetDoc(userId Id, docId Id) (*Doc, error)
	AddDoc(userId Id, newDoc Doc) error
	EditDoc(userId Id, newDoc Doc) (*Doc, error)
	EditDocAccess(userId Id, docId Id, request DocAccessRequest) (*Doc, error)
	DeleteDoc(userId Id, docId Id) error
	GetAllDocs(userId Id) ([]Doc, error)

	CreateGroup(userId Id, group Group) (*Group, error)
	DeleteGroup(userId Id, groupId Id) error
	EditGroup(userId Id, groupId Id, newGroup Group) (*Group, error)
	GetGroupById(groupId Id) (*Group, error)

	AddMember(userId Id, groupId Id, newMemberId Id) error
	RemoveMember(userId Id, groupId Id, memberId Id) error
	GetMembers(request GroupMembersChunkRequest) ([]User, error)

	GenerateNewUserId() Id
	GenerateNewDocId() Id
	GenerateNewGroupId() Id
}
