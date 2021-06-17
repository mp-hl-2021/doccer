package model

import (
	"doccer/auth"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type ModelImpl struct {
	storage Storage
	jwtHandler auth.JwtHandler
}

func NewModelImpl(storage Storage, secret []byte) ModelImpl {
	return ModelImpl{
		storage: storage,
		jwtHandler: auth.NewJwtHandler(secret, 24 * time.Hour),
	}
}

func (s *ModelImpl) nextId(isUserId bool) Id {
	if isUserId {
		return s.storage.GenerateNewUserId()
	}
	return s.storage.GenerateNewDocId()
}

func (s *ModelImpl) Register(request LoginRequest) (*User, error) {
	user := User{
		Id:    s.nextId(true),
		Login: request.Login,
	}

	if s.storage.CheckLoginExists(user.Login) {
		return nil, ErrAlreadyExists
	}

	encryptedPassword, err := auth.EncodeStr(request.Password)
	if err != nil {
		return nil, err
	}
	err = s.storage.AddUser(user, encryptedPassword)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *ModelImpl) Login(request LoginRequest) (*LoginResponse, error) {
	user, err := s.storage.GetUserByLogin(request.Login)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := s.storage.GetHashedPassword(user.Id)
	if err != nil {
		return nil, err
	}
	err = auth.Compare([]byte(request.Password), *hashedPassword)
	if err != nil {
		return nil, ErrWrongPassword
	}

	claims := auth.UserClaims{
		UserId: string(user.Id),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.jwtHandler.ExpirationTime).Unix(),
		},
	}
	token, err := s.jwtHandler.GetNewToken(claims)
	if err != nil {
		return nil, err
	}
	resp := LoginResponse{User: *user, Token: Token(token)}
	return &resp, nil
}

func (s *ModelImpl) Auth(tokenStr string) (*string, error) {
	claims, err := s.jwtHandler.ParseClaims(tokenStr, auth.UserClaims{})
	if err != nil {
		return nil, err
	}
	userClaims := (*claims).(*auth.UserClaims)
	return &userClaims.UserId, nil
}

func (s *ModelImpl) Logout(token Token) error {
	return nil
}

func (s *ModelImpl) CreateDoc(userId Id, doc Doc) (*Doc, error) {
	doc = Doc {
		Id:       s.storage.GenerateNewDocId(),
		AuthorId: userId,
		Text:     doc.Text,
		Access:   doc.Access,
	}
	docId, err := s.storage.AddDoc(doc)
	if err != nil {
		return nil, err
	}
	doc.Id = *docId
	return &doc, err
}

func (s *ModelImpl) GetDoc(userId Id, docId Id) (*Doc, error) {
	return s.getDoc(userId, docId, true)
}

func (s *ModelImpl) getDoc(userId Id, docId Id, shouldCheck bool) (*Doc, error) {
	var realAccess string = ""
	if shouldCheck {
		checkAccess, err := s.storage.CheckAccess(userId, docId)
		if err != nil || checkAccess == "none" {
			return nil, ErrNoAccess
		}
		realAccess = checkAccess
	}


	res, err := s.storage.GetDoc(docId)
	if err != nil {
		return nil, err
	}

	doc := Doc{
		Id:       docId,
		AuthorId: res.AuthorId,
		Text:     res.Text,
		Access:   realAccess,
	}
	return &doc, nil
}

func (s *ModelImpl) EditDoc(userId Id, newDoc Doc) (*Doc, error) {
	checkAccess, err := s.storage.CheckAccess(userId, newDoc.Id)
	if err != nil || checkAccess == "none" || checkAccess == "read" {
		return nil, ErrNoAccess
	}
	oldDoc, err := s.getDoc(userId, newDoc.Id, false)
	if err != nil {
		return nil, ErrNotFound
	}

	if oldDoc.Access != newDoc.Access && checkAccess != "absolute" {
		return nil, ErrNoAccess
	}
	return s.storage.EditDoc(newDoc)
}

func (s *ModelImpl) DeleteDoc(userId Id, docId Id) error {
	checkAccess, err := s.storage.CheckAccess(userId, docId)
	if err != nil || checkAccess != "absolute" {
		return ErrNoAccess
	}
	return s.storage.DeleteDoc(docId)
}

func (s *ModelImpl) ChangeDocAccess(userId Id, request DocAccessRequest) (*Doc, error) {
	acc, err := s.storage.CheckAccess(userId, request.DocId)
	if err != nil || acc != "absolute" {
		return nil, ErrNoAccess
	}
	err = s.storage.EditDocAccess(request.DocId, request)
	if err != nil {
		return nil, err
	}
	return s.getDoc(userId, request.DocId, false)
}

func (s *ModelImpl) GetAllDocs(userId Id) ([]Doc, error) {
	return s.storage.GetAllDocs(userId)
}

func (s *ModelImpl) GetUserById(userId Id) (*User, error) {
	return s.storage.GetUser(userId)
}

func (s *ModelImpl) EditUser(userId Id, newUser User) (*User, error) {
	if userId != newUser.Id {
		return nil, ErrNoAccess
	}
	return s.storage.EditUser(newUser)
}

func (s *ModelImpl) CreateGroup(userId Id, group Group) (*Group, error) {
	group = Group{
		Id:      s.storage.GenerateNewGroupId(),
		Name:    group.Name,
		Creator: userId,
	}

	return s.storage.CreateGroup(group)
}

func (s *ModelImpl) DeleteGroup(userId Id, groupId Id) error {
	g, err := s.storage.GetGroupById(groupId)
	if err != nil {
		return err
	}
	if g.Creator != userId {
		return ErrNoAccess
	}
	return s.storage.DeleteGroup(groupId)
}

func (s *ModelImpl) EditGroup(userId Id, newGroup Group) (*Group, error) {
	oldGroup, err := s.storage.GetGroupById(newGroup.Id)
	if err != nil || userId != oldGroup.Creator  {
		return nil, ErrNoAccess
	}
	return s.storage.EditGroup(newGroup)
}

func (s *ModelImpl) AddMember(userId Id, groupId Id, newMemberId Id) error {
	group, err := s.storage.GetGroupById(groupId)
	if err != nil {
		return err
	}
	if group.Creator != userId {
		return ErrNoAccess
	}

	_, err = s.storage.GetUser(newMemberId)
	if err != nil {
		return ErrNotFound
	}
	return s.storage.AddMember(groupId, newMemberId)
}

func (s *ModelImpl) RemoveMember(userId Id, groupId Id, memberId Id) error {
	group, err := s.storage.GetGroupById(groupId)
	if err != nil {
		return err
	}
	if group.Creator != userId {
		return ErrNoAccess
	}

	_, err = s.storage.GetUser(memberId)
	if err != nil {
		return ErrNotFound
	}
	return s.storage.RemoveMember(groupId, memberId)
}

func (s *ModelImpl) GetMembers(userId Id, request GroupMembersChunkRequest) ([]User, error) {
	g, err := s.storage.GetGroupById(request.Id)
	if err != nil {
		return nil, err
	}
	if g.Creator != userId {
		return nil, ErrNoAccess
	}
	return s.storage.GetMembers(request)
}



