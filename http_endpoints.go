package heroku_store

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HttpEndpoints interface {
	MakeCreateUserEndpoint() func(w http.ResponseWriter, r *http.Request)
	MakeListUserEndpoint() func(w http.ResponseWriter, r *http.Request)
}

type httpEndpoints struct {
	store UserStore
}

func NewHttpEndpoints(s UserStore) HttpEndpoints {
	return &httpEndpoints{store: s}
}

func (h *httpEndpoints) MakeListUserEndpoint() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.store.List()
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err)
			return
		}
		respondJSON(w, http.StatusOK, users)
	}
}

func (h *httpEndpoints) MakeCreateUserEndpoint() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &User{}
		dataBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err)
			return
		}
		err = json.Unmarshal(dataBytes, &user)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err)
			return
		}
		newUser, err := h.store.Create(user)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err)
			return
		}
		respondJSON(w, http.StatusCreated, newUser)
	}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
