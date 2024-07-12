package handlers

import (
	"encoding/json"
	"fmt"
	middleware "kangym/api/middleware"
	"kangym/db"
	accounts "kangym/types"
	utils "kangym/utils"
	errors "kangym/utils/errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	listenAddr string
	store      db.Storage
}

func NewAccountHandler(store db.Storage) *AccountHandler {
	return &AccountHandler{
		store: store,
	}
}
func (s *AccountHandler) HandleAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccountById(w, r)
	case "DELETE":
		return s.handleDeleteAccountById(w, r)
	case "POST":
		return s.handleRestoreAccountById(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *AccountHandler) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccounts(w)
	case "POST":
		return s.handleCreateAccount(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *AccountHandler) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}
	return utils.WriteJSON(w, http.StatusOK, account)
}

func (s *AccountHandler) handleGetAccounts(w http.ResponseWriter) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return utils.WriteJSON(w, http.StatusOK, accounts)
}

func (s *AccountHandler) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(accounts.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}
	account, err := s.store.CreateAccount(createAccountReq)
	if err != nil {
		return err
	}
	tokenString, err := middleware.CreateJWT(strconv.Itoa(account.ID))
	if err != nil {
		return err
	}
	fmt.Println(tokenString)
	return utils.WriteJSON(w, http.StatusOK, createAccountReq)
}

func (s *AccountHandler) handleDeleteAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccountById(id); err != nil {
		errors.DBError(w)
		return nil
	}
	return utils.WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *AccountHandler) handleRestoreAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.RestoreAccountById(id); err != nil {
		errors.DBError(w)
		return nil
	}
	return utils.WriteJSON(w, http.StatusOK, map[string]int{"restored": id})
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
