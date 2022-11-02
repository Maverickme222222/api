package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Maverickme222222/api/logic"
)

// UserHandler is a handler for user related requests
type UserHandler struct {
	logic logic.Logic
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(logic logic.Logic) UserHandler {
	return UserHandler{
		logic: logic,
	}
}

func (u UserHandler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	type body struct {
		Name string `json:"name"`
	}

	var reqBody body

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		log.Fatalf("Decode error %v", err)
	}

	res, _ := u.logic.CreateNewUser(ctx, reqBody.Name)

	Respond(w, res, http.StatusOK, true, nil)
}

// Respond converts a Go value to JSON and sends it to the client.
func Respond(w http.ResponseWriter, data interface{}, statusCode int, success bool, errors interface{}) {

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	if data == nil {
		data = map[string]interface{}{}
	}

	if errors == nil {
		errors = []error{}
	}

	resp := struct {
		Success bool        `json:"success"`
		Errors  interface{} `json:"errors"`
		Data    interface{} `json:"data"`
	}{
		Success: success,
		Data:    data,
		Errors:  errors,
	}

	//nolint:golint,errcheck,gosec
	json.NewEncoder(w).Encode(resp)
}
