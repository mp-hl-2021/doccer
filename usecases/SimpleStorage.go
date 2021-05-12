package usecases

import (
	"doccer/model"
	"strconv"
)

type SimpleStorage struct {
	users      map[model.Id]model.User
	docs       map[model.Id]model.Doc
	passwords  map[model.Id]model.Password
	friends    map[model.Id][]model.User
	groups     map[model.Id]model.Group
	members    map[model.Id][]model.User

	loginToUser map[string]model.User

	lastUserId  int
	lastDocId   int
	lastGroupId int
}

func (s *SimpleStorage) AddUser(newUser model.User, password model.Password) error {
	newUser = model.User{
		Id: s.GenerateNewUserId(),
		Login: newUser.Login,
	}
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


func (s *SimpleStorage) EditUser(userId model.Id, newUser model.User) (*model.User, error) {
	if userId != newUser.Id {
		return nil, model.ErrNoAccess
	}
	s.users[userId] = newUser
	return &newUser, nil
}

func (s *SimpleStorage) GetUserById(userId model.Id) (*model.User, error) {
	user, exists := s.users[userId]
	if !exists {
		return nil, model.ErrNotFound
	}
	return &user, nil
}

func (s *SimpleStorage) CheckAccess(userId model.Id, docId model.Id) (string, error) {
	return "absolute", nil
}

func (s *SimpleStorage) GetDoc(userId model.Id, docId model.Id) (*model.Doc, error) {
	checkAccess, err := s.CheckAccess(userId, docId)
	if err != nil || checkAccess == "none" {
		return nil, model.ErrNoAccess
	}
	doc, exists := s.docs[docId]
	if !exists {
		return nil, model.ErrNotFound
	}
	return &doc, nil
}

func (s *SimpleStorage) AddDoc(userId model.Id, doc model.Doc) error {
	doc = model.Doc{
		Id: s.GenerateNewDocId(),
		AuthorId: userId,
		Text: doc.Text,
		Access: doc.Access,
	}
	s.docs[doc.Id] = doc
	return nil
}

func (s *SimpleStorage) EditDoc(userId model.Id, newDoc model.Doc) (*model.Doc, error) {
	checkAccess, err := s.CheckAccess(userId, newDoc.Id)
	if err != nil || checkAccess == "none" || checkAccess == "read" {
		return nil, model.ErrNoAccess
	}
	s.docs[newDoc.Id] = newDoc
	return &newDoc, nil
}

func (s *SimpleStorage) GetAllDocs(userId model.Id) ([]model.Doc, error) {
	var docs []model.Doc
	for _, doc := range s.docs {
		_ = append(docs, doc)
	}
	return docs, nil
}

func (s *SimpleStorage) DeleteDoc(userId model.Id, docId model.Id) error {
	checkAccess, err := s.CheckAccess(userId, docId)
	if err != nil || checkAccess != "absolute" {
		return model.ErrNoAccess
	}
	_, exists := s.docs[docId]
	if !exists {
		return nil
	}
	delete(s.docs, docId)
	return nil
}

func (s *SimpleStorage) GetFriends(userId model.Id) ([]model.User, error) {
	friends, exists := s.friends[userId]
	if !exists {
		return nil, model.ErrNotFound
	}
	return friends, nil
}

func (s *SimpleStorage) AddFriend(userId model.Id, friendId model.Id) error {
	_, exists := s.friends[userId]
	friend, err := s.GetUserById(friendId)
	if err != nil {
		return model.ErrNotFound
	}
	if !exists {
		s.friends[userId] = []model.User{*friend}
	} else {
		for _,f := range s.friends[userId] {
			if f == *friend {
				return nil
			}
		}
		_ = append(s.friends[userId], *friend)
	}
	return nil
}

func (s *SimpleStorage) RemoveFriend(userId model.Id, friendId model.Id) error {
	_, exists := s.friends[userId]
	friend, err := s.GetUserById(friendId)
	if err != nil {
		return model.ErrNotFound
	}
	if !exists {
		return nil
	}
	var index int
	for i,f := range s.friends[userId] {
		if f == *friend {
			index = i
		}
	}
	s.friends[userId] = append(s.friends[userId][:index], s.friends[userId][index+1:]...)
	return nil
}

func (s *SimpleStorage) CreateGroup(userId model.Id, group model.Group) (*model.Group, error) {
	group = model.Group{
		Id: s.GenerateNewGroupId(),
		Name: group.Name,
		Creator: userId,
	}
	s.groups[group.Id] = group
	return &group, nil
}

func (s *SimpleStorage) DeleteGroup(userId model.Id, groupId model.Id) error {
	_, exists := s.groups[groupId]
	if !exists {
		return nil
	}
	if s.groups[groupId].Creator != userId {
		return model.ErrNoAccess
	}
	delete(s.groups, groupId)
	return nil
}

func (s *SimpleStorage) EditGroup(userId model.Id, groupId model.Id, newGroup model.Group) (*model.Group, error) {
	_, exists := s.groups[groupId]
	if !exists {
		s.groups[groupId] = newGroup
		return &newGroup, nil
	}
	if s.groups[groupId].Creator != userId {
		return nil, model.ErrNoAccess
	}
	s.groups[groupId] = newGroup
	return &newGroup, nil
}

func (s *SimpleStorage) GetGroupById(groupId model.Id) (*model.Group, error) {
	group, exists := s.groups[groupId]
	if !exists {
		return nil, model.ErrNotFound
	}
	return &group, nil
}

func (s *SimpleStorage) AddMember(userId model.Id, groupId model.Id, newMemberId model.Id) error {
	_, exists := s.members[groupId]
	member, err := s.GetUserById(newMemberId)
	if err != nil {
		return model.ErrNotFound
	}
	group, err := s.GetGroupById(groupId)
	if err != nil {
		return model.ErrNotFound
	}
	if group.Creator != userId {
		return model.ErrNoAccess
	}
	if !exists {
		s.members[newMemberId] = []model.User{*member}
	} else {
		for _,m := range s.members[groupId] {
			if m == *member {
				return nil
			}
		}
		_ = append(s.members[newMemberId], *member)
	}
	return nil
}

func (s *SimpleStorage) RemoveMember(userId model.Id, groupId model.Id, memberId model.Id) error {
	group, groupExists := s.GetGroupById(groupId)
	if groupExists != nil {
		return model.ErrNotFound
	}
	if group.Creator != userId {
		return model.ErrNoAccess
	}
	_, exists := s.members[groupId]
	member, err := s.GetUserById(memberId)
	if err != nil {
		return model.ErrNotFound
	}
	if !exists {
		return nil
	}
	var index int
	for i,m := range s.members[groupId] {
		if m == *member {
			index = i
		}
	}
	s.members[groupId] = append(s.members[groupId][:index], s.members[groupId][index+1:]...)
	return nil
}

func (s *SimpleStorage) GetMembers(request model.GroupMembersChunkRequest) ([]model.User, error) {
	if request.Size == 0 {
		return nil, nil
	}
	members, exists := s.members[request.Id]
	if !exists {
		return nil, model.ErrNotFound
	}
	if request.Begin < 0 || request.Size + request.Begin > len(members) {
		return nil, model.ErrNotFound
	}
	return members[request.Begin:(request.Begin + request.Size)], nil
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

func (s *SimpleStorage) GenerateNewGroupId() model.Id {
	s.lastGroupId += 1
	id := strconv.Itoa(s.lastGroupId)

	return model.Id(id)
}
