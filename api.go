package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

type JwtAuthClaim struct {
	Account *Account `json:"account"`
	jwt.StandardClaims
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

const jwtSecrete = "hunteaar999"

func createJWT(account *Account) (string, error) {

	claims := JwtAuthClaim{
		Account: account,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
	}
	secrete := []byte(jwtSecrete)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secrete)

}

func valuidateJwt(tokenString string) (*jwt.Token, error) {
	// fmt.Printf(tokenString)

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecrete), nil
	})
}

func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")

		bearerToken := strings.Split(tokenString, " ")
		if len(bearerToken) < 2 {
			WriteJson(w, http.StatusUnauthorized, ApiError{Error: "invalid token"})
			return
		}
		_, err := valuidateJwt(bearerToken[1])
		if err != nil {

			WriteJson(w, http.StatusUnauthorized, ApiError{Error: err.Error()})
			return
		}
		handlerFunc(w, r)
	}
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr,
		store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccounts))
	router.HandleFunc("/account/transfer", makeHttpHandleFunc(s.handleTransfer))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHttpHandleFunc(s.handleAccountByID)))

	log.Println("Json Api server running on port", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccounts(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccounts(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	case "PUT":
		return s.handleGetAccounts(w, r)
	default:
		return fmt.Errorf("method now allowd %s", r.Method)
	}
}

func (s *APIServer) handleAccountByID(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccountByID(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)

	default:
		return fmt.Errorf("method now allowd %s", r.Method)
	}

}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, accounts)
}
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	accounts, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	account := new(CreateAccount)

	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		return err
	}

	newAccount := NewAccount(account.FirstName, account.LastName)
	if err := s.store.CreateAccount(newAccount); err != nil {
		return err
	}
	jwtStr, err := createJWT(newAccount)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return WriteJson(w, http.StatusOK, map[string]interface{}{"accountData": newAccount, "token": jwtStr})
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]int{"delete": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := TransferRequest{}
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJson(w, http.StatusOK, transferReq)
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("Invalid id given%s", idStr)
	}
	return id, nil
}
