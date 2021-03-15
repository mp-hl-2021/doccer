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

	router.HandleFunc("/docs/{doc_id}/get", a.getDoc).Methods(http.MethodGet)
	router.HandleFunc("/docs/{doc_id}/create", a.createDoc).Methods(http.MethodPost)
	router.HandleFunc("/docs/{doc_id}/delete", a.deleteDoc).Methods(http.MethodPost)
	router.HandleFunc("/docs/{doc_id}/edit", a.editDoc).Methods(http.MethodPost)
	router.HandleFunc("/docs/{doc_id}/change_doc_access", a.changeDocAccess).Methods(http.MethodPost)
	router.HandleFunc("/docs/get_all_docs", a.getAllDocs).Methods(http.MethodGet)

	router.HandleFunc("/users/get_friends", a.getFriends).Methods(http.MethodGet)

	router.HandleFunc("/users/groups/create", a.createGroup).Methods(http.MethodPost)
	router.HandleFunc("/users/groups/delete", a.deleteGroup).Methods(http.MethodPost)
	router.HandleFunc("/users/groups/add_member", a.addMember).Methods(http.MethodPost)
	router.HandleFunc("/users/groups/remove_member", a.removeMember).Methods(http.MethodPost)
	router.HandleFunc("/users/groups/get_members", a.getGroupMembers).Methods(http.MethodGet)

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

func (a *Api) createGroup(w http.ResponseWriter, r *http.Request) {
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

func (a *Api) addMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	w.WriteHeader(http.StatusOK)
}

func (a *Api) removeMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	w.WriteHeader(http.StatusOK)
}

func (a *Api) getGroupMembers(w http.ResponseWriter, r *http.Request) {
	var m model.GroupMembersChunkRequest
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
