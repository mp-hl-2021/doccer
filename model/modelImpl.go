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

type UserClaims struct {
	UserId    Id
	jwt.StandardClaims
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

func (s *ModelImpl) Register(request LoginRequest) error {
	user := User{
		Id:    s.nextId(true),
		Login: request.Login,
	}

	if s.storage.CheckLoginExists(user.Login) {
		return ErrAlreadyExists
	}

	encryptedPassword, err := auth.EncodeStr(request.Password)
	if err != nil {
		return err
	}
	err = s.storage.AddUser(user, encryptedPassword)
	if err != nil {
		return err
	}
	return nil
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

	claims := UserClaims{
		user.Id,
		jwt.StandardClaims{
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
	claims, err := s.jwtHandler.ParseClaims(tokenStr, UserClaims{})
	if err != nil {
		return nil, err
	}
	userClaims := (*claims).(*UserClaims)
	return &userClaims.Id, nil
}

func (s *ModelImpl) Logout(token Token) error {
	return nil
}

func (s *ModelImpl) CreateDoc(userId Id, doc Doc) (*Doc, error) {
	return nil, ErrNotImplemented
}

func (s *ModelImpl) GetDoc(userId Id, docId Id) (*Doc, error) {
	return nil, ErrNotImplemented
}

func (s *ModelImpl) EditDoc(userId Id, newDoc Doc) (*Doc, error) {
	return nil, ErrNotImplemented
}

func (s *ModelImpl) DeleteDoc(userId Id, docId Id) error {
	return ErrNotImplemented
}

func (s *ModelImpl) GetAllDocs(userId Id) ([]Doc, error) {
	return nil, ErrNotImplemented
}

func (s *ModelImpl) GetUserById(userId Id) (*User, error) {
	return nil, ErrNotImplemented
}

func (s *ModelImpl) EditUser(userId Id, newUser User) (*User, error) {
	return nil, ErrNotImplemented
}

func (s *ModelImpl) GetFriends(userId Id) ([]User, error) {
	return nil, ErrNotImplemented
}

func (s *ModelImpl) AddFriend(userId Id, friendId Id) error {
	return ErrNotImplemented
}


func (s *ModelImpl) RemoveFriend(userId Id, friendId Id) error {
	return ErrNotImplemented
}

func (s *ModelImpl) CreateGroup(userId Id, group Group) (*Group, error) {
	return nil, ErrNotImplemented
}

func (s *ModelImpl) DeleteGroup(userId Id, groupId Id) error {
	return ErrNotImplemented
}

func (s *ModelImpl) EditGroup(userId Id, groupId Id, newGroup Group) (*Group, error) {
	return nil, ErrNotImplemented
}

func (s *ModelImpl) AddMember(userId Id, groupId Id, newMemberId Id) error {
	return ErrNotImplemented
}

func (s *ModelImpl) RemoveMember(userId Id, groupId Id, memberId Id) error {
	return ErrNotImplemented
}

func (s *ModelImpl) GetMembers(userId Id, request GroupMembersChunkRequest) ([]User, error) {
	return nil, ErrNotImplemented
}


