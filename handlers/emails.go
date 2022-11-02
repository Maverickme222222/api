package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Maverickme222222/api/logic"
)

// EmailHandler is a handler for user related requests
type EmailHandler struct {
	logic logic.Logic
}

// NewUserHandler creates a new UserHandler
func NewEmailHandler(logic logic.Logic) EmailHandler {
	return EmailHandler{
		logic: logic,
	}
}

func (e EmailHandler) CreateNewEmail(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	type body struct {
		Name string `json:"name"`
	}

	var reqBody body

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		log.Fatalf("Decode error %v", err)
	}

	res, _ := e.logic.CreateNewEmail(ctx, reqBody.Name)

	Respond(w, res, http.StatusOK, true, nil)
}
