package model

type Storage interface {
	GetUser(userId Id) (*User, error)
	GetUserByLogin(login string) (*User, error)
	GetDoc(docId Id) (*Doc, error)
	GetHashedPassword(userId Id) (*Password, error)
    AddUser(newUser User, password Password) error
	CheckLoginExists(login string) bool

	GenerateNewUserId() Id
	GenerateNewDocId() Id
}
