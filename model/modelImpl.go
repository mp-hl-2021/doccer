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
	docId, err := s.storage.AddDoc(userId, doc)
	if err != nil {
		return nil, err
	}
	doc.Id = *docId
	return &doc, err
}

func (s *ModelImpl) GetDoc(userId Id, docId Id) (*Doc, error) {
	return s.storage.GetDoc(userId, docId)
}

func (s *ModelImpl) EditDoc(userId Id, newDoc Doc) (*Doc, error) {
	return s.storage.EditDoc(userId, newDoc)
}

func (s *ModelImpl) DeleteDoc(userId Id, docId Id) error {
	return s.storage.DeleteDoc(userId, docId)
}

func (s *ModelImpl) GetAllDocs(userId Id) ([]Doc, error) {
	return s.storage.GetAllDocs(userId)
}

func (s *ModelImpl) GetUserById(userId Id) (*User, error) {
	return s.storage.GetUserById(userId)
}

func (s *ModelImpl) EditUser(userId Id, newUser User) (*User, error) {
	return s.storage.EditUser(userId, newUser)
}

func (s *ModelImpl) CreateGroup(userId Id, group Group) (*Group, error) {
	return s.storage.CreateGroup(userId, group)
}

func (s *ModelImpl) DeleteGroup(userId Id, groupId Id) error {
	return s.storage.DeleteGroup(userId, groupId)
}

func (s *ModelImpl) EditGroup(userId Id, groupId Id, newGroup Group) (*Group, error) {
	return s.storage.EditGroup(userId, groupId, newGroup)
}

func (s *ModelImpl) AddMember(userId Id, groupId Id, newMemberId Id) error {
	return s.storage.AddMember(userId, groupId, newMemberId)
}

func (s *ModelImpl) RemoveMember(userId Id, groupId Id, memberId Id) error {
	return s.storage.RemoveMember(userId, groupId, memberId)
}

func (s *ModelImpl) GetMembers(userId Id, request GroupMembersChunkRequest) ([]User, error) {
	return s.storage.GetMembers(request)
}

func (s *ModelImpl) ChangeDocAccess(userId Id, request DocAccessRequest) (*Doc, error) {
	res, err := s.storage.EditDocAccess(userId, request.DocId, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}



