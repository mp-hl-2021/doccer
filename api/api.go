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
	router.HandleFunc("/docs/{doc_id}", a.auth(a.createDoc, false)).Methods(http.MethodPost)
	router.HandleFunc("/docs/{doc_id}", a.auth(a.deleteDoc, true)).Methods(http.MethodDelete)
	router.HandleFunc("/docs/{doc_id}", a.auth(a.editDoc, true)).Methods(http.MethodPut)

	router.HandleFunc("/docs", a.auth(a.getAllDocs, true)).Methods(http.MethodGet)

	router.HandleFunc("/users", a.auth(a.editUser, true)).Methods(http.MethodPut)
	router.HandleFunc("/users", a.auth(a.getUser, true)).Methods(http.MethodGet)

	router.HandleFunc("/users/friends", a.auth(a.getFriends, true)).Methods(http.MethodGet)
	router.HandleFunc("/users/friends", a.auth(a.putFriend, true)).Methods(http.MethodPut)

	router.HandleFunc("/users/groups", a.auth(a.createGroup, true)).Methods(http.MethodPost)
	router.HandleFunc("/users/groups", a.auth(a.deleteGroup, true)).Methods(http.MethodDelete)
	router.HandleFunc("/users/groups", a.auth(a.editGroup, true)).Methods(http.MethodPut)

	router.HandleFunc("/users/groups/{group_id}/members", a.auth(a.putMember, true)).Methods(http.MethodPut)
	router.HandleFunc("/users/groups/{group_id}/members", a.auth(a.getMembers, true)).Methods(http.MethodGet)

	return router
}

func (a *Api) register(w http.ResponseWriter, r *http.Request) {
	var m model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := a.useCases.Register(m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err == model.ErrAlreadyExists {
			_, _ = w.Write([]byte("login already exists"))
		}
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
		ctx := context.WithValue(r.Context(), "myUserId", userId)
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
	id := mux.Vars(r)["id"]
	if myId != nil {
		println("Get doc request with id", id, "by user", myId)
	} else {
		println("Get doc request with id", id, "by unauthorized user")
	}

	w.WriteHeader(http.StatusOK)
}

func (a *Api) createDoc(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId != nil {
		println("Create doc request by user", myId)
	} else {
		println("Create doc request by unauthorized user")
	}
	var m model.Doc
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) deleteDoc(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	myId := r.Context().Value("myUserId")
	if myId == nil {
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
	println("Edit doc request with id", id, "by user", myId)

	var m model.Doc
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) changeDocAccess(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	println("Change docs access request by user", myId)

	var m model.DocAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) getAllDocs(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	println("Get all docs request by user", myId)

	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) getUser(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	println("Get user request by user", myId)

	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) editUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) getFriends(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	println("Get friends request by user", myId)

	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) putFriend(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	println("Add friend request by user", myId)

	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) createGroup(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	println("create group request by user", myId)

	var m model.Group
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) editGroup(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	id := mux.Vars(r)["id"]
	println("Edit group request with group id ", id, "by user", myId)

	var m model.Group
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) deleteGroup(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	id := mux.Vars(r)["id"]

	println("Delete group request with group id ", id, "by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) putMember(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	id := mux.Vars(r)["id"]

	println("Add member to group request with group id ", id, "by user", myId)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) getMembers(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("myUserId")
	if myId == nil {
		return
	}
	id := mux.Vars(r)["id"]
	println("Get group members request with group id ", id, "by user", myId)

	var m model.GroupMembersChunkRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
