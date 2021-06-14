package storage

import (
	"context"
	"database/sql"
	"doccer/model"
	_ "github.com/lib/pq"
	"strconv"
)

type PostgresStorage struct {
	Dbc *sql.DB
}

func (p *PostgresStorage) GenerateNewUserId() model.Id {
	ctx := context.Background()
	tx, _ := p.Dbc.BeginTx(ctx, nil)

	row := tx.QueryRow("select g.last_user_id from GeneralInfo g where g.base_id = 0")

	lastId := 0
	_ = row.Scan(&lastId)
	id := strconv.Itoa(lastId)
	lastId += 1

	_, _ = tx.ExecContext(ctx, "update GeneralInfo g set g.last_user_id = ? where g.base_id = 0", lastId)

	_ = tx.Commit()

	return model.Id(id)
}

func (p *PostgresStorage) GenerateNewDocId() model.Id {
	ctx := context.Background()
	tx, _ := p.Dbc.BeginTx(ctx, nil)

	row := tx.QueryRow("select g.last_doc_id from GeneralInfo g where g.base_id = 0")

	lastId := 0
	_ = row.Scan(&lastId)
	id := strconv.Itoa(lastId)
	lastId += 1

	_, _ = tx.ExecContext(ctx, "update GeneralInfo g set g.last_doc_id = ? where g.base_id = 0", lastId)

	_ = tx.Commit()

	return model.Id(id)
}

func (p *PostgresStorage) GenerateNewGroupId() model.Id {
	ctx := context.Background()
	tx, _ := p.Dbc.BeginTx(ctx, nil)

	row := tx.QueryRow("select g.last_group_id from GeneralInfo g where g.base_id = 0")

	lastId := 0
	_ = row.Scan(&lastId)
	id := strconv.Itoa(lastId)
	lastId += 1

	_, _ = tx.ExecContext(ctx, "update GeneralInfo g set g.last_group_id = ? where g.base_id = 0", lastId)

	_ = tx.Commit()

	return model.Id(id)
}

func (p *PostgresStorage) AddUser(newUser model.User, password model.Password) error {
	newUser = model.User{
		Id:    p.GenerateNewUserId(),
		Login: newUser.Login,
	}
	ctx := context.Background()
	tx, err := p.Dbc.BeginTx(ctx, nil);
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "insert into Users value (?, ?)", newUser.Id, newUser.Login)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "insert into Password value(?, ?)", newUser.Id, password)
	if err != nil {
		return err
	}
	return nil
}

func (p * PostgresStorage) GetUserByLogin(login string) (*model.User, error) {
	res := p.Dbc.QueryRow("select u.id from Users u where u.login = ?", login)
	id := 0
	err := res.Scan(&id)
	if err != nil {
		return nil, model.ErrNotFound
	}

	user := model.User {
		Id: model.Id(strconv.Itoa(id)),
		Login: login,
	}
	return &user, nil
}

func (p * PostgresStorage) GetUser(userId model.Id) (*model.User, error) {
	res := p.Dbc.QueryRow("select u.login from Users u where u.id = ?", userId)
	login := ""
	err := res.Scan(&login)
	if err != nil {
		return nil, model.ErrNotFound
	}

	user := model.User {
		Id: userId,
		Login: login,
	}
	return &user, nil
}

func (p * PostgresStorage) CheckLoginExists(login string) bool {
	_, err := p.GetUserByLogin(login)
	return err != nil
}

func (p * PostgresStorage) GetHashedPassword(userId model.Id) (*model.Password, error) {
	res := p.Dbc.QueryRow("select u.password from Password p where p.user_id = ?", userId)
	passwordStr := []byte("")
	err := res.Scan(&passwordStr)
	if err != nil {
		return nil, model.ErrNotFound
	}
	password := model.Password(passwordStr)
	return &password, nil
}


func (p * PostgresStorage) EditUser(userId model.Id, newUser model.User) (*model.User, error) {
	if userId != newUser.Id {
		return nil, model.ErrNoAccess
	}
	_, err := p.Dbc.Exec("update Users u set u.login = ? where u.id = ?", userId, newUser.Login)
	if err != nil {
		return nil, model.ErrNotFound
	}
	return &newUser, nil
}

func (p * PostgresStorage) GetUserById(userId model.Id) (*model.User, error) {
	return p.GetUser(userId)
}

func (p * PostgresStorage) CheckAccess(userId model.Id, docId model.Id) (string, error) {
	res := p.Dbc.QueryRow("select access from DocMemberAccess where doc_id = ? and user_id = ?", docId, userId)
	userAccess := 0
	err := res.Scan(&userAccess)
	if err == nil {
		return accessIntToStr(userAccess), nil
	}

	res2, err := p.Dbc.Query("select group_id, access from DocGroupAccess where doc_id = ?", docId, userId)
	if err == nil {
		for res2.Next() {
			naccess := 0
			groupId := ""
			err = res2.Scan(&groupId, &naccess)
			if err == nil {
				members, err := p.GetMembers(model.GroupMembersChunkRequest{
					Id:    model.Id(groupId),
					Begin: 0,
					Size:  0,
				})

				if err == nil {
					contains := false
					for _, member := range members {
						if member.Id == userId {
							contains = true
						}
					}

					if contains && naccess > userAccess {
						userAccess = naccess
					}
				}
			}
		}
		return accessIntToStr(userAccess), nil
	}


	doc, err := p.getDoc(userId, docId, false)
	if err != nil {
		return "none", err
	}
	return doc.Access, nil
}

func (p * PostgresStorage) GetDoc(userId model.Id, docId model.Id) (*model.Doc, error) {
	return p.getDoc(userId, docId, true)
}

func (p * PostgresStorage) getDoc(userId model.Id, docId model.Id, shouldCheck bool) (*model.Doc, error) {
	var realAccess string = ""
	if shouldCheck {
		checkAccess, err := p.CheckAccess(userId, docId)
		if err != nil || checkAccess == "none" {
			return nil, model.ErrNoAccess
		}
		realAccess = checkAccess
	}


	res := p.Dbc.QueryRow("select u.text, u.creator_id, u.public_access_type from Users u where u.id = ?", userId)
	text := ""
	creatorId := ""
	pub_access := ""
	err := res.Scan(&text, &creatorId, &pub_access)
	if err != nil {
		return nil, model.ErrNotFound
	}

	if realAccess == "" {
		realAccess = pub_access
	}

	doc := model.Doc{
		Id:       docId,
		AuthorId: model.Id(creatorId),
		Text:     text,
		Access:   realAccess,
	}
	return &doc, nil
}

func (p * PostgresStorage) AddDoc(userId model.Id, doc model.Doc) error {
	doc = model.Doc{
		Id:       p.GenerateNewDocId(),
		AuthorId: userId,
		Text:     doc.Text,
		Access:   doc.Access,
	}
	_, err := p.Dbc.Exec("insert into Docs values (?, ?, ?, ?)", doc.Id, doc.AuthorId, doc.Text, doc.Access)
	if err != nil {
		return model.ErrAlreadyExists
	}
	return nil
}

func (p * PostgresStorage) EditDoc(userId model.Id, newDoc model.Doc) (*model.Doc, error) {
	checkAccess, err := p.CheckAccess(userId, newDoc.Id)
	if err != nil || checkAccess == "none" || checkAccess == "read" {
		return nil, model.ErrNoAccess
	}
	oldDoc, err := p.GetDoc(userId, newDoc.Id)
	if err != nil {
		return nil, model.ErrNotFound
	}

	if oldDoc.Access != newDoc.Access && checkAccess != "absolute" {
		return nil, model.ErrNoAccess
	}
	publicAccType := accessStrToInt(newDoc.Access)
	_, _ = p.Dbc.Exec("update Docs set text = ?, public_access_type = ?", publicAccType)
	return &newDoc, nil
}

func (p * PostgresStorage) EditDocAccess(userId model.Id, docId model.Id, editRequest model.DocAccessRequest) (*model.Doc, error) {
	acc, err := p.CheckAccess(userId, docId)
	if err != nil || acc != "absolute" {
		return nil, model.ErrNoAccess
	}
	if editRequest.Type == 0 {
		_, err := p.Dbc.Exec("insert into DocMemberRestriction values (?, ?, ?) on conflict(doc_id, member_id) update", docId, editRequest.ItemId, accessStrToInt(editRequest.Access))
		if err != nil {
			return nil, err
		}
	} else {
		_, err := p.Dbc.Exec("insert into DocGroupRestriction values (?, ?, ?) on conflict(doc_id, member_id) update ", docId, editRequest.ItemId, accessStrToInt(editRequest.Access))
		if err != nil {
			return nil, err
		}
	}
	return p.GetDoc(userId, docId)
}

func (p * PostgresStorage) GetAllDocs(userId model.Id) ([]model.Doc, error) {
	res, err := p.Dbc.Query("select u.id, u.text, u.access from Docs d where d.creator_id = ?", userId)
	if err != nil {
		return nil, err
	}
	var docs []model.Doc

	for res.Next() {
		text := ""
		id := ""
		access := 0
		_ = res.Scan(&id, &text, &access)

		accessStr := accessIntToStr(access)
		_ = append(docs, model.Doc {
			Id: model.Id(id),
			AuthorId: userId,
			Text: text,
			Access: accessStr,
		})
	}
	return docs, nil
}

func (p * PostgresStorage) DeleteDoc(userId model.Id, docId model.Id) error {
	checkAccess, err := p.CheckAccess(userId, docId)
	if err != nil || checkAccess != "absolute" {
		return model.ErrNoAccess
	}
	_, err = p.Dbc.Exec("delete from Docs d where d.id = ?", docId)
	if err != nil {
		return model.ErrNotFound
	}
	return nil
}

func (p * PostgresStorage) CreateGroup(userId model.Id, group model.Group) (*model.Group, error) {
	group = model.Group{
		Id:      p.GenerateNewGroupId(),
		Name:    group.Name,
		Creator: userId,
	}

	_, _ = p.Dbc.Exec("insert into Groups1 value (?, ?, ?)", group.Id, group.Creator, group.Name)
	return &group, nil
}

func (p * PostgresStorage) DeleteGroup(userId model.Id, groupId model.Id) error {
	g, _ := p.GetGroupById(groupId)
	if g.Creator != userId {
		return model.ErrNoAccess
	}
	_, err := p.Dbc.Exec("delete from Groups1 g where g.id = ?", groupId)
	if err != nil {
		return model.ErrNotFound
	}
	return nil
}

func (p * PostgresStorage) EditGroup(userId model.Id, groupId model.Id, newGroup model.Group) (*model.Group, error) {
	oldGroup, err := p.GetGroupById(groupId)
	if err != nil || userId != oldGroup.Id  {
		return nil, model.ErrNoAccess
	}

	_, err = p.Dbc.Exec("update Groups1 set name = ? where id = ?", newGroup.Name, groupId)
	return &newGroup, nil
}

func (p * PostgresStorage) GetGroupById(groupId model.Id) (*model.Group, error) {
	res := p.Dbc.QueryRow("select g.name, g.creator_id from Groups1 g where g.id = ?", groupId)
	name := ""
	creatorId := ""
	err := res.Scan(&name, creatorId)

	if err != nil {
		return nil, model.ErrNotFound
	}

	group := model.Group{
		Id:      groupId,
		Name:    name,
		Creator: model.Id(creatorId),
	}
	return &group, nil
}

func (p * PostgresStorage) AddMember(userId model.Id, groupId model.Id, newMemberId model.Id) error {
	group, err := p.GetGroupById(groupId)
	if err != nil {
		return err
	}
	if group.Creator != userId {
		return model.ErrNoAccess
	}

	_, err = p.GetUserById(newMemberId)
	if err != nil {
		return model.ErrNotFound
	}

	_, err = p.Dbc.Exec("insert into GroupMemebers value (?, ?)", groupId, userId)
	if err != nil {
		return model.ErrAlreadyExists
	}

	return nil
}

func (p * PostgresStorage) RemoveMember(userId model.Id, groupId model.Id, memberId model.Id) error {
	group, err := p.GetGroupById(groupId)
	if err != nil {
		return err
	}
	if group.Creator != userId {
		return model.ErrNoAccess
	}

	_, err = p.GetUserById(memberId)
	if err != nil {
		return model.ErrNotFound
	}

	_, err = p.Dbc.Exec("delete from GroupMemebers g where g.id = ? and g.member_id = ?", groupId, userId)
	if err != nil {
		return model.ErrNotFound
	}

	return nil
}

func (p * PostgresStorage) GetMembers(request model.GroupMembersChunkRequest) ([]model.User, error) {
	res, err := p.Dbc.Query("select g.member_id from GroupMembers g where g.group_id = ?", request.Id)
	if err != nil {
		return nil, err
	}
	var users []model.User

	for res.Next() {
		user_id := ""
		_ = res.Scan(&user_id)
		user, _ := p.GetUserById(model.Id(user_id))
		_ = append(users, *user)
	}
	return users, nil
}

func accessStrToInt(accessStr string) int {
	accessCode := 1 //read
	if accessStr == "absolute" {
		accessCode = 3
	}
	if accessStr == "edit" {
		accessCode = 2
	}
	if accessStr == "none" {
		accessCode = 0
	}
	return accessCode
}

func accessIntToStr(accessCode int) string {
	accessStr := "read"
	if accessCode == 3 {
		accessStr = "absolute"
	}
	if accessCode == 1 {
		accessStr = "edit"
	}
	if accessCode == 0 {
		accessStr = "none"
	}
	return accessStr
}