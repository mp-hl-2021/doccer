package model

import (
	"doccer/auth"
	"doccer/data"
	"doccer/linter"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type ModelImpl struct {
	storage Storage
	jwtHandler auth.JwtHandler
	processChannel chan data.Doc
	resChannel chan data.Doc
}

func NewModelImpl(
	storage Storage,
	secret []byte,
	linter linter.GeneralLinter,
    saveWorkersCnt int,
	linterWorkersCnt int,
	) ModelImpl {

	processChannel := make(chan data.Doc, linterWorkersCnt * 2)
	resChannel := make(chan data.Doc, saveWorkersCnt * 2)

	res := ModelImpl{
		storage: storage,
		jwtHandler: auth.NewJwtHandler(secret, 24 * time.Hour),
		processChannel: processChannel,
		resChannel: resChannel,
	}

	for i := 0; i < saveWorkersCnt; i++ {
		go func() {
			for {
				newDoc, ok := <-resChannel
				if !ok {
					break
				}
				_, _ = res.editDoc(newDoc.AuthorId, newDoc, false)
			}
		}()
	}

	for i := 0; i < linterWorkersCnt; i++ {
		go func() {
			for {
				docToProcess, ok := <- processChannel
				if !ok {
					break
				}
				resChannel <- linter.CheckCode(docToProcess)
			}
		}()
	}

	return res
}

func (s *ModelImpl) nextId(isUserId bool) data.Id {
	if isUserId {
		return s.storage.GenerateNewUserId()
	}
	return s.storage.GenerateNewDocId()
}

func (s *ModelImpl) Register(request LoginRequest) (*data.User, error) {
	user := data.User{
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

func (s *ModelImpl) CreateDoc(userId data.Id, doc data.Doc) (*data.Doc, error) {
	doc = data.Doc{
		Id:       s.storage.GenerateNewDocId(),
		AuthorId: userId,
		Text:     doc.Text,
		Access:   doc.Access,
		Lang: doc.Lang,
		LinterStatus: "No inspection",
	}
	docId, err := s.storage.AddDoc(doc)

	if err != nil {
		return nil, err
	}
	doc.Id = *docId
	s.processChannel <- doc
	return &doc, err
}

func (s *ModelImpl) GetDoc(userId data.Id, docId data.Id) (*data.Doc, error) {
	return s.getDoc(userId, docId, true)
}

func (s *ModelImpl) getDoc(userId data.Id, docId data.Id, shouldCheck bool) (*data.Doc, error) {
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

	res.Access = realAccess
	return res, nil
}

func (s *ModelImpl) EditDoc(userId data.Id, newDoc data.Doc) (*data.Doc, error) {
	return s.editDoc(userId, newDoc, true)
}

func (s *ModelImpl) LaunchLinter(userId data.Id, docId data.Id) error {
	doc, err := s.getDoc(userId, docId, true)
	if err != nil {
		return err
	}
	access := doc.Access
	if access != "edit" && access != "absolute" {
		return ErrNoAccess
	}
	s.processChannel <- *doc
	return nil
}

func (s *ModelImpl) editDoc(userId data.Id, newDoc data.Doc, updateLinter bool) (*data.Doc, error) {
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
	res, err := s.storage.EditDoc(newDoc)
	if err != nil {
		return nil, err
	}

	if updateLinter {
		s.processChannel <- *res
	}

	return res, nil
}

func (s *ModelImpl) DeleteDoc(userId data.Id, docId data.Id) error {
	checkAccess, err := s.storage.CheckAccess(userId, docId)
	if err != nil || checkAccess != "absolute" {
		return ErrNoAccess
	}
	return s.storage.DeleteDoc(docId)
}

func (s *ModelImpl) ChangeDocAccess(userId data.Id, request DocAccessRequest) (*data.Doc, error) {
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

func (s *ModelImpl) GetAllDocs(userId data.Id) ([]data.Doc, error) {
	return s.storage.GetAllDocs(userId)
}

func (s *ModelImpl) GetUserById(userId data.Id) (*data.User, error) {
	return s.storage.GetUser(userId)
}

func (s *ModelImpl) EditUser(userId data.Id, newUser data.User) (*data.User, error) {
	if userId != newUser.Id {
		return nil, ErrNoAccess
	}
	return s.storage.EditUser(newUser)
}

func (s *ModelImpl) CreateGroup(userId data.Id, group data.Group) (*data.Group, error) {
	group = data.Group{
		Id:      s.storage.GenerateNewGroupId(),
		Name:    group.Name,
		Creator: userId,
	}

	return s.storage.CreateGroup(group)
}

func (s *ModelImpl) DeleteGroup(userId data.Id, groupId data.Id) error {
	g, err := s.storage.GetGroupById(groupId)
	if err != nil {
		return err
	}
	if g.Creator != userId {
		return ErrNoAccess
	}
	return s.storage.DeleteGroup(groupId)
}

func (s *ModelImpl) EditGroup(userId data.Id, newGroup data.Group) (*data.Group, error) {
	oldGroup, err := s.storage.GetGroupById(newGroup.Id)
	if err != nil || userId != oldGroup.Creator  {
		return nil, ErrNoAccess
	}
	return s.storage.EditGroup(newGroup)
}

func (s *ModelImpl) AddMember(userId data.Id, groupId data.Id, newMemberId data.Id) error {
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

func (s *ModelImpl) RemoveMember(userId data.Id, groupId data.Id, memberId data.Id) error {
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

func (s *ModelImpl) GetMembers(userId data.Id, request GroupMembersChunkRequest) ([]data.User, error) {
	g, err := s.storage.GetGroupById(request.Id)
	if err != nil {
		return nil, err
	}
	if g.Creator != userId {
		return nil, ErrNoAccess
	}
	return s.storage.GetMembers(request)
}



