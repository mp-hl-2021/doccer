package api

import (
	"context"
	"doccer/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Api struct {
	useCases model.UseCasesInterface
}
func NewApi(x model.UseCasesInterface) *Api {
	return &Api{
		useCases: x,
	}
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/register", a.register).Methods(http.MethodPost)
	router.HandleFunc("/login", a.login).Methods(http.MethodPost)
	router.HandleFunc("/logout", a.auth(a.logout, true)).Methods(http.MethodPost)

	router.HandleFunc("/docs/{doc_id}", a.auth(a.getDoc, false)).Methods(http.MethodGet)
	router.HandleFunc("/docs", a.auth(a.createDoc, false)).Methods(http.MethodPost)
	router.HandleFunc("/docs/{doc_id}", a.auth(a.deleteDoc, true)).Methods(http.MethodDelete)
	router.HandleFunc("/docs/{doc_id}", a.auth(a.editDoc, true)).Methods(http.MethodPut)

	router.HandleFunc("/docs/{doc_id}/access", a.auth(a.changeDocAccess, true)).Methods(http.MethodPost)

	router.HandleFunc("/docs", a.auth(a.getAllDocs, true)).Methods(http.MethodGet)

	router.HandleFunc("/users", a.auth(a.editUser, true)).Methods(http.MethodPut)
	router.HandleFunc("/users", a.auth(a.getUser, true)).Methods(http.MethodGet)

	router.HandleFunc("/users/groups", a.auth(a.createGroup, true)).Methods(http.MethodPost)
	router.HandleFunc("/users/groups", a.auth(a.deleteGroup, true)).Methods(http.MethodDelete)
	router.HandleFunc("/users/groups", a.auth(a.editGroup, true)).Methods(http.MethodPut)

	router.HandleFunc("/users/groups/{group_id}/members", a.auth(a.getMembers, true)).Methods(http.MethodGet)
	router.HandleFunc("/users/groups/{group_id}/members", a.auth(a.removeMember, true)).Methods(http.MethodDelete)
	router.HandleFunc("/users/groups/{group_id}/members", a.auth(a.addMember, true)).Methods(http.MethodPut)

	return router
}

func (a *Api) register(w http.ResponseWriter, r *http.Request) {
	var m model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := a.useCases.Register(m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err == model.ErrAlreadyExists {
			_, _ = w.Write([]byte("login already exists"))
		}
		return
	}
	respJson, err := json.Marshal(user)
	_, err = w.Write(respJson)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) login(w http.ResponseWriter, r *http.Request) {
	var m model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loginResponse, err := a.useCases.Login(m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/jwt")
	respJson, err := json.Marshal(loginResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) auth(f func (w http.ResponseWriter, r *http.Request), isRequired bool) func (w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("AuthToken")
		if token == "" {
			if isRequired {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Missing authorization token"))
			}
			return
		}

		userId, err := a.useCases.Auth(token)
		if err != nil {
			if isRequired {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Wrong authorization token"))
			}
			return
		}
		ctx := context.WithValue(r.Context(), "myUserId", *userId)
		f(w, r.WithContext(ctx))
	}
}

func (a *Api) logout(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	println("Logout request by user ", myId)
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) getDoc(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	id := model.Id(mux.Vars(r)["doc_id"])
	var newDoc *model.Doc
	if myId != nil {
		doc, err := a.useCases.GetDoc(model.Id(myId.(string)), id)
		newDoc = doc
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		println("Get doc request with id", id, "by user", myId)
	} else {
		doc, err := a.useCases.GetDoc("-1", id)
		newDoc = doc
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		println("Get doc request with id", id, "by unauthorized user")
	}
	respJson, err := json.Marshal(newDoc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) createDoc(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	var m model.Doc
	var newDoc *model.Doc
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if myId != nil {
		doc, err := a.useCases.CreateDoc(model.Id(myId.(string)), m)
		newDoc = doc
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		println("Create doc request by user", myId)
	} else {
		doc, err := a.useCases.CreateDoc("-1", m)
		newDoc = doc
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		println("Create doc request by unauthorized user")
	}
	respJson, err := json.Marshal(newDoc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) deleteDoc(w http.ResponseWriter, r *http.Request) {
	id := model.Id(mux.Vars(r)["id"])
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	err := a.useCases.DeleteDoc(model.Id(myId.(string)), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Delete doc request with id", id, "by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) editDoc(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	var m model.Doc
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m = model.Doc{
		Id: model.Id(id),
		AuthorId: m.AuthorId,
		Access: m.Access,
		Text: m.Text,
	}
	doc, err := a.useCases.EditDoc(model.Id(myId.(string)), m)
	respJson, err := json.Marshal(doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Edit doc request with id", id, "by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) changeDocAccess(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}

	var m model.DocAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	doc, err := a.useCases.ChangeDocAccess(model.Id(myId.(string)), m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respJson, err := json.Marshal(doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Change docs access request by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) getAllDocs(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}

	docs, err := a.useCases.GetAllDocs(model.Id(myId.(string)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respJson, err := json.Marshal(docs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Get all docs request by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) getUser(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	user, err := a.useCases.GetUserById(model.Id(myId.(string)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respJson, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Get user request by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) editUser(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	var m model.User
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m = model.User{
		Id: model.Id(myId.(string)),
		Login: m.Login,
	}
	user, err := a.useCases.EditUser(model.Id(myId.(string)), m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respJson, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Edit user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) createGroup(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	var m model.Group
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	group, err := a.useCases.CreateGroup(model.Id(myId.(string)), m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respJson, err := json.Marshal(group)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("create group request by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) editGroup(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	id := mux.Vars(r)["id"]

	var m model.Group
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m = model.Group{
		Id: model.Id(id),
		Name: m.Name,
		Creator: m.Creator,
	}
	group, err := a.useCases.EditGroup(model.Id(myId.(string)), model.Id(id), m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respJson, err := json.Marshal(group)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Edit group request with group id ", id, "by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) deleteGroup(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	id := mux.Vars(r)["id"]
	err := a.useCases.DeleteGroup(model.Id(myId.(string)), model.Id(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Delete group request with group id ", id, "by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) removeMember(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	var m model.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := a.useCases.RemoveMember(model.Id(myId.(string)), m.GroupId, m.MemberId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Add member to group request with group id ", m.GroupId, "by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) addMember(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	var m model.MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := a.useCases.AddMember(model.Id(myId.(string)), m.GroupId, m.MemberId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Add member to group request with group id ", m.GroupId, "by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) getMembers(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	var m model.GroupMembersChunkRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	members, err := a.useCases.GetMembers(model.Id(myId.(string)), m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respJson, err := json.Marshal(members)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	println("Get group members request with group id ", m.Id, "by user", myId)
	w.WriteHeader(http.StatusOK)
}
