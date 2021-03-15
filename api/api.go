package api

import (
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
	router := mux.newRouter

	router.HandleFunc("/register", a.register).Methods(http.MethodPost)
	router.HandleFunc("/login", a.login).Methods(http.MethodPost)
	router.HandleFunc("/logout", a.logout).Methods(http.MethodPost)

	router.HandleFunc("/docs/{doc_id}", a.getDoc).Methods(http.MethodGet)
	router.HandleFunc("/docs/{doc_id}", a.createDoc).Methods(http.MethodPost)
	router.HandleFunc("/docs/{doc_id}", a.deleteDoc).Methods(http.MethodDelete)
	router.HandleFunc("/docs/{doc_id}", a.editDoc).Methods(http.MethodPut)

	router.HandleFunc("/docs/get_all_docs", a.getAllDocs).Methods(http.MethodGet)

	router.HandleFunc("/users", a.editUser).Methods(http.MethodPut)
	router.HandleFunc("/users", a.getUser).Methods(http.MethodGet)

	router.HandleFunc("/users/friends", a.getFriends).Methods(http.MethodGet)
	router.HandleFunc("/users/friends", a.putFriend).Methods(http.MethodPut)

	router.HandleFunc("/users/groups", a.createGroup).Methods(http.MethodPost)
	router.HandleFunc("/users/groups", a.deleteGroup).Methods(http.MethodDelete)
	router.HandleFunc("/users/groups", a.editGroup).Methods(http.MethodPut)

	router.HandleFunc("/users/groups/{group_id}/members", a.putMember).Methods(http.MethodPut)
	router.HandleFunc("/users/groups/{group_id}/members", a.getMembers).Methods(http.MethodGet)

	return router
}

func (a *Api) register(w http.ResponseWriter, r *http.Request) {
	var m model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	token := "custom token"

	w.Header().Set("Content-Type", "application/jwt")
	if _, err := w.Write([]byte(token)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) getDoc(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	w.WriteHeader(http.StatusOK)
}

func (a *Api) createDoc(w http.ResponseWriter, r *http.Request) {
	var m model.Doc
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) deleteDoc(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	w.WriteHeader(http.StatusOK)
}

func (a *Api) editDoc(w http.ResponseWriter, r *http.Request) {
	var m model.Doc
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) changeDocAccess(w http.ResponseWriter, r *http.Request) {
	var m model.DocAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) getAllDocs(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) getFriends(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) putFriend(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) createGroup(w http.ResponseWriter, r *http.Request) {
	var m model.Group
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) editGroup(w http.ResponseWriter, r *http.Request) {
	var m model.Group
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) deleteGroup(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	w.WriteHeader(http.StatusOK)
}

func (a *Api) putMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	w.WriteHeader(http.StatusOK)
}

func (a *Api) getMembers(w http.ResponseWriter, r *http.Request) {
	var m model.GroupMembersChunkRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
