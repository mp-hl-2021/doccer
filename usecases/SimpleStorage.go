package usecases

import (
	"doccer/model"
	"strconv"
)

type SimpleStorage struct {
	users      map[model.Id]model.User
	docs       map[model.Id]model.Doc
	passwords  map[model.Id]model.Password

	loginToUser map[string]model.User

	lastUserId int
	lastDocId  int
}

func (s *SimpleStorage) AddUser(newUser model.User, password model.Password) error {
	s.users[newUser.Id] = newUser
	s.loginToUser[newUser.Login] = newUser
	s.passwords[newUser.Id] = password
	return nil
}

func (s *SimpleStorage) GetUserByLogin(login string) (*model.User, error) {
	user, exists := s.loginToUser[login]
	if !exists {
		return nil, model.ErrNotFound
	}

	return &user, nil
}

func (s *SimpleStorage) GetUser(userId model.Id) (*model.User, error) {
	user, exists := s.users[userId]
	if !exists {
		return nil, model.ErrNotFound
	}

	return &user, nil
}

func (s *SimpleStorage) GetDoc(docId model.Id) (*model.Doc, error) {
	doc, exists := s.docs[docId]
	if !exists {
		return nil, model.ErrNotFound
	}

	return &doc, nil
}

func (s *SimpleStorage) CheckLoginExists(login string) bool {
	_, exists := s.loginToUser[login]
	return exists
}

func (s *SimpleStorage) GetHashedPassword(userId model.Id) (*model.Password, error) {
	password, exists := s.passwords[userId]
	if !exists {
		return nil, model.ErrNotFound
	}

	return &password, nil
}

func (s *SimpleStorage) GenerateNewUserId() model.Id {
	s.lastUserId += 1
	id := strconv.Itoa(s.lastUserId)

	return model.Id(id)
}

func (s *SimpleStorage) GenerateNewDocId() model.Id {
	s.lastDocId += 1
	id := strconv.Itoa(s.lastDocId)

	return model.Id(id)
}
