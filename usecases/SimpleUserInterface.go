package usecases

import (
	"doccer/model"
	"errors"
	"fmt"
)

type SimpleUserInterface struct {
	users      []model.User
	docs       []model.Doc
	lastUserId int
	lastDocId  int
}

func (s *SimpleUserInterface) nextId(isUserId bool) model.Id {
	if isUserId {
		s.lastUserId += 1
		return model.Id(fmt.Sprintf("%d", s.lastUserId))
	}
	s.lastDocId += 1
	return model.Id(fmt.Sprintf("%d", s.lastDocId))
}

func (s *SimpleUserInterface) Register(request model.LoginRequest) (*model.LoginResponse, error) {
	user := model.User{
		Id:    s.nextId(true),
		Login: request.Login,
	}

	s.users = append(s.users, user)
	return &model.LoginResponse{
		Token: model.Token(user.Id),
		User:  user,
	}, nil
}

func (s *SimpleUserInterface) Login(request model.LoginRequest) (*model.LoginResponse, error) {
	for _, u := range s.users {
		if u.Login == request.Login {
			return &model.LoginResponse{
				Token: model.Token(u.Id),
				User:  u,
			}, nil
		}
	}
	return nil, errors.New("invalid username")
}

func (s *SimpleUserInterface) Logout(token model.Token) error {
	return nil
}

func (s *SimpleUserInterface) CreateDoc(token model.Token, doc model.Doc) (*model.Doc, error) {
	newDoc := model.Doc{
		Id:       s.nextId(false),
		AuthorId: model.Id(token),
		Text:     doc.Text,
		Access:   doc.Access,
	}
	s.docs = append(s.docs, newDoc)
	return &newDoc, nil
}

func (s *SimpleUserInterface) GetDoc(token model.Token, docId model.Id) (*model.Doc, error) {
	for _, d := range s.docs {
		if d.Id == docId {
			return &d, nil
		}
	}
	return nil, errors.New("incorrect docId")
}

func (s *SimpleUserInterface) EditDoc(token model.Token, newDoc model.Doc) (*model.Doc, error) {
	for _, d := range s.docs {
		if d.Id == newDoc.Id {
			d = newDoc
		}
	}
	return &newDoc, nil
}

func (s *SimpleUserInterface) DeleteDoc(token model.Token, docId model.Id) error {
	for i, d := range s.docs {
		if d.Id == docId {
			s.docs[i] = s.docs[len(s.docs) - 1]
			s.docs = s.docs[:len(s.docs) - 1]
			return nil
		}
	}
	return errors.New("incorrect docId")
}

func (s *SimpleUserInterface) ChangeDocAccess(token model.Token, request model.DocAccessRequest) error {
	for i, d := range s.docs {
		if d.Id == request.DocId {
			s.docs[i].Access = request.Access
			return nil
		}
	}
	return errors.New("incorrect docId")
}

func (s *SimpleUserInterface) GetAllDocs(token model.Token) ([]model.Doc, error) {
	return s.docs, nil
}

func (s *SimpleUserInterface) GetUserById(userId model.Id) (*model.User, error) {
	for _, u := range s.users {
		if u.Id == userId {
			return &u, nil
		}
	}
	return nil, errors.New("incorrect userId")
}

func (s *SimpleUserInterface) GetFriends(token model.Token) ([]model.User, error) {
	return nil, errors.New("not implemented")
}

func (s *SimpleUserInterface) CreateGroup(token model.Token, group model.Group) (*model.Group, error) {
	return nil, errors.New("not implemented")
}

func (s *SimpleUserInterface) DeleteGroup(token model.Token, groupId model.Id) error {
	return errors.New("not implemented")
}

func (s *SimpleUserInterface) AddMember(token model.Token, groupId model.Id, newMemberId model.Id) error {
	return errors.New("not implemented")
}

func (s *SimpleUserInterface) RemoveMember(token model.Token, groupId model.Id, memberId model.Id) error {
	return errors.New("not implemented")
}

func (s *SimpleUserInterface) GetMembers(token model.Token, request model.GroupMembersChunkRequest) ([]model.User, error) {
	return nil, errors.New("not implemented")
}

func NewSimpleUserInterface() *SimpleUserInterface {
	return &SimpleUserInterface{
		users: make([]model.User, 0),
		docs:  make([]model.Doc, 0),
	}
}



