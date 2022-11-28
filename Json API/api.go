package main

import (
	"encoding/json"
	"net/http"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHttpHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if e, ok := err.(apiError); ok {
				writeJson(w, e.Status, e)
				return
			}
			writeJson(w, http.StatusInternalServerError, apiError{Err: "internal Server", Status: http.StatusInternalServerError})
		}
	}
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return apiError{Err: "Invalid Method", Status: http.StatusMethodNotAllowed}
	}
	user := User{}
	if !user.Valid {
		return ErrUserInvalid
	}
	return writeJson(w, http.StatusOK, User{})
}
