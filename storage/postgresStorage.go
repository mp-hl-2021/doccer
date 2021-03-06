package storage

import (
	"context"
	"database/sql"
	"doccer/data"
	"doccer/model"
	_ "github.com/lib/pq"
	"strconv"
	"sync"
)

type PostgresStorage struct {
	mu1 sync.Mutex
	mu2 sync.Mutex
	mu3 sync.Mutex
	Dbc *sql.DB
}

func (p *PostgresStorage) ClearAllTables() {
	_, _ = p.Dbc.Exec("TRUNCATE Users, DocGroupRestriction, DocMemberRestriction, Docs, GroupMember, GeneralInfo, Groups1, Password CASCADE ;")
	_, _ = p.Dbc.Exec("insert into GeneralInfo values (0, 0, 0, 0)")
}

func (p *PostgresStorage) GenerateNewUserId() data.Id {
	p.mu1.Lock()
	ctx := context.Background()
	tx, _ := p.Dbc.BeginTx(ctx, nil)

	row := tx.QueryRow("select g.last_user_id from GeneralInfo g where g.base_id = 0")

	lastId := 0
	_ = row.Scan(&lastId)
	id := strconv.Itoa(lastId)
	lastId += 1

	_, _ = tx.ExecContext(ctx, "update GeneralInfo set last_user_id = $1 where base_id = 0", lastId)

	_ = tx.Commit()

	p.mu1.Unlock()
	return data.Id(id)
}

func (p *PostgresStorage) GenerateNewDocId() data.Id {
	p.mu2.Lock()
	ctx := context.Background()
	tx, _ := p.Dbc.BeginTx(ctx, nil)

	row := tx.QueryRow("select g.last_doc_id from GeneralInfo g where g.base_id = 0")

	lastId := 0
	_ = row.Scan(&lastId)
	id := strconv.Itoa(lastId)
	lastId += 1

	_, _ = tx.ExecContext(ctx, "update GeneralInfo set last_doc_id = $1 where base_id = 0", lastId)

	_ = tx.Commit()

	p.mu2.Unlock()
	return data.Id(id)
}

func (p *PostgresStorage) GenerateNewGroupId() data.Id {
	p.mu3.Lock()
	ctx := context.Background()
	tx, _ := p.Dbc.BeginTx(ctx, nil)

	row := tx.QueryRow("select g.last_group_id from GeneralInfo g where g.base_id = 0")

	lastId := 0
	_ = row.Scan(&lastId)
	id := strconv.Itoa(lastId)
	lastId += 1

	_, _ = tx.ExecContext(ctx, "update GeneralInfo set last_group_id = $1 where base_id = 0", lastId)

	_ = tx.Commit()

	p.mu3.Unlock()
	return data.Id(id)
}

func (p *PostgresStorage) AddUser(newUser data.User, password model.Password) error {
	ctx := context.Background()
	tx, err := p.Dbc.BeginTx(ctx, nil);
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(string(newUser.Id))
	_, err = tx.ExecContext(ctx, "insert into Users values ($1, $2)", id , newUser.Login)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, "insert into Password values ($1, $2)", id, password)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p * PostgresStorage) GetUserByLogin(login string) (*data.User, error) {
	res := p.Dbc.QueryRow("select u.id from Users u where u.login = $1", login)
	id := 0
	err := res.Scan(&id)
	if err != nil {
		return nil, model.ErrNotFound
	}

	user := data.User{
		Id:    data.Id(strconv.Itoa(id)),
		Login: login,
	}
	return &user, nil
}

func (p * PostgresStorage) GetUser(userId data.Id) (*data.User, error) {
	id, _ := strconv.Atoi(string(userId))
	res := p.Dbc.QueryRow("select u.login from Users u where u.id = $1", id)
	login := ""
	err := res.Scan(&login)
	if err != nil {
		return nil, model.ErrNotFound
	}

	user := data.User{
		Id: userId,
		Login: login,
	}
	return &user, nil
}

func (p * PostgresStorage) CheckLoginExists(login string) bool {
	_, err := p.GetUserByLogin(login)
	return err == nil
}

func (p * PostgresStorage) GetHashedPassword(userId data.Id) (*model.Password, error) {
	res := p.Dbc.QueryRow("select p.password from Password p where p.id = $1", userId)
	passwordStr := []byte("")
	err := res.Scan(&passwordStr)
	if err != nil {
		return nil, model.ErrNotFound
	}
	password := model.Password(passwordStr)
	return &password, nil
}


func (p * PostgresStorage) EditUser(newUser data.User) (*data.User, error) {
	_, err := p.Dbc.Exec("update Users u set u.login = $1 where u.id = $2", newUser.Id, newUser.Login)
	if err != nil {
		return nil, model.ErrNotFound
	}
	return &newUser, nil
}

func (p * PostgresStorage) CheckAccess(userId data.Id, docId data.Id) (string, error) {
	doc, err := p.GetDoc(docId)
	if err != nil {
		return "none", err
	}

	if userId == "-1" {
		return doc.Access, nil
	}

	if doc.AuthorId == userId {
		return "absolute", nil
	}

	res := p.Dbc.QueryRow("select type from DocMemberRestriction where doc_id = $1 and member_id = $2", docId, userId)
	userAccess := 0
	err = res.Scan(&userAccess)
	if err == nil {
		return accessIntToStr(userAccess), nil
	}

	res2, err := p.Dbc.Query("select group_id, type from DocGroupRestriction where doc_id = $1", docId)
	if err == nil {
		userAccess = accessStrToInt(doc.Access)
		for res2.Next() {
			naccess := 0
			groupId := ""
			err = res2.Scan(&groupId, &naccess)
			if err == nil {
				members, err := p.GetMembers(model.GroupMembersChunkRequest{
					Id:    data.Id(groupId),
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
	return doc.Access, nil
}

func (p *PostgresStorage) GetDoc(docId data.Id) (*data.Doc, error) {
	res := p.Dbc.QueryRow("select d.text, d.creator_id, d.public_access_type, d.lang, d.lstatus from Docs d where d.id = $1", docId)
	text := ""
	creatorId := ""
	pubAccess := 0
	lang := ""
	lstatus := ""
	err := res.Scan(&text, &creatorId, &pubAccess, &lang, &lstatus)
	if err != nil {
		return nil, model.ErrNotFound
	}
	return &data.Doc{
		Id:           docId,
		AuthorId:     data.Id(creatorId),
		Text:         text,
		Access:       accessIntToStr(pubAccess),
		Lang:         lang,
		LinterStatus: lstatus,
	}, nil
}

func (p * PostgresStorage) AddDoc(doc data.Doc) (*data.Id, error) {
	_, err := p.Dbc.Exec("insert into Docs values ($1, $2, $3, $4, $5, $6)",
		doc.Id, doc.AuthorId, doc.Text, accessStrToInt(doc.Access), doc.Lang, doc.LinterStatus)
	if err != nil {
		return nil, model.ErrAlreadyExists
	}
	return &doc.Id, nil
}

func (p * PostgresStorage) EditDoc(newDoc data.Doc) (*data.Doc, error) {
	publicAccType := accessStrToInt(newDoc.Access)
	_, _ = p.Dbc.Exec("update Docs set text = $1, public_access_type = $2, lang = $3, lstatus = $4 where id = $5",
		newDoc.Text, publicAccType, newDoc.Lang, newDoc.LinterStatus, newDoc.Id)
	return &newDoc, nil
}

func (p * PostgresStorage) EditDocAccess(docId data.Id, editRequest model.DocAccessRequest) error {
	if editRequest.Type == 0 {
		_, err := p.Dbc.Exec("insert into DocMemberRestriction values ($1, $2, $3) on conflict(doc_id, member_id) do update set type = excluded.type;", docId, editRequest.ItemId, accessStrToInt(editRequest.Access))
		if err != nil {
			return err
		}
	} else {
		_, err := p.Dbc.Exec("insert into DocGroupRestriction values ($1, $2, $3) on conflict(doc_id, group_id) do update set type = excluded.type", docId, editRequest.ItemId, accessStrToInt(editRequest.Access))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p * PostgresStorage) GetAllDocs(userId data.Id) ([]data.Doc, error) {
	res, err := p.Dbc.Query("select d.id, d.text, d.access, d.lang, d.lstatus from Docs d where d.creator_id = $1", userId)
	if err != nil {
		return nil, err
	}
	var docs []data.Doc

	for res.Next() {
		text := ""
		id := ""
		access := 0
		lang := ""
		lstatus := ""
		_ = res.Scan(&id, &text, &access, &lang, &lstatus)

		accessStr := accessIntToStr(access)
		docs = append(docs, data.Doc{
			Id:           data.Id(id),
			AuthorId:     userId,
			Text:         text,
			Access:       accessStr,
			Lang:         lang,
			LinterStatus: lstatus,
		})
	}
	return docs, nil
}

func (p * PostgresStorage) DeleteDoc(docId data.Id) error {
	_, err := p.Dbc.Exec("delete from Docs d where d.id = $1", docId)
	if err != nil {
		return model.ErrNotFound
	}
	return nil
}

func (p * PostgresStorage) CreateGroup(group data.Group) (*data.Group, error) {
	_, err := p.Dbc.Exec("insert into Groups1 values ($1, $2, $3)", group.Id, group.Creator, group.Name)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (p * PostgresStorage) DeleteGroup(groupId data.Id) error {
	_, err := p.Dbc.Exec("delete from Groups1 g where g.id = $1", groupId)
	if err != nil {
		return model.ErrNotFound
	}
	return nil
}

func (p * PostgresStorage) EditGroup(newGroup data.Group) (*data.Group, error) {
	_, err := p.Dbc.Exec("update Groups1 set name = $1 where id = $2", newGroup.Name, newGroup.Id)
	if err != nil {
		return nil, err
	}
	return &newGroup, nil
}

func (p * PostgresStorage) GetGroupById(groupId data.Id) (*data.Group, error) {
	res := p.Dbc.QueryRow("select g.name, g.creator_id from Groups1 g where g.id = $1", groupId)
	name := ""
	creatorId := ""
	err := res.Scan(&name, &creatorId)

	if err != nil {
		return nil, model.ErrNotFound
	}

	group := data.Group{
		Id:      groupId,
		Name:    name,
		Creator: data.Id(creatorId),
	}
	return &group, nil
}

func (p * PostgresStorage) AddMember(groupId data.Id, newMemberId data.Id) error {
	_, err := p.Dbc.Exec("insert into GroupMember values ($1, $2)", groupId, newMemberId)
	if err != nil {
		return model.ErrAlreadyExists
	}

	return nil
}

func (p * PostgresStorage) RemoveMember(groupId data.Id, memberId data.Id) error {
	_, err := p.Dbc.Exec("delete from GroupMember g where g.group_id = $1 and g.member_id = $2", groupId, memberId)
	if err != nil {
		return model.ErrNotFound
	}

	return nil
}

func (p * PostgresStorage) GetMembers(request model.GroupMembersChunkRequest) ([]data.User, error) {
	res, err := p.Dbc.Query("select g.member_id from GroupMember g where g.group_id = $1", request.Id)
	if err != nil {
		return nil, err
	}
	var users []data.User

	for res.Next() {
		user_id := ""
		_ = res.Scan(&user_id)
		user, _ := p.GetUser(data.Id(user_id))
		users = append(users, *user)
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