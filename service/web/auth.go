package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	OTP      string `json:"otp"`
	StoreID  string `json:"store"`
	Series   string `json:"series"`
}

// LoginUser logins a user into the store
func (srv Web) LoginUser(w http.ResponseWriter, r *http.Request) {
	req := srv.decodeLoginRequest(w, r)
	if req == nil {
		return
	}

	if err := loginValidate(req); err != nil {
		formatStandardResponse("login", err.Error(), w)
		return
	}

	if err := srv.Store.LoginUser(req.Email, req.Password, req.OTP, req.StoreID, req.Series); err != nil {
		formatStandardResponse("login", err.Error(), w)
		return
	}

	formatStandardResponse("", "", w)
}

// Macaroon checks that the store macaroon is available
func (srv Web) Macaroon(w http.ResponseWriter, r *http.Request) {
	headers, err := srv.Store.Macaroon()
	if err != nil {
		log.Printf("Error retrieving macaroon: %v", err)
		formatStandardResponse("auth", "No store macaroon found", w)
		return
	}

	formatRecordResponse(headers, w)
}

func (srv Web) decodeLoginRequest(w http.ResponseWriter, r *http.Request) *loginRequest {
	// Decode the JSON body
	req := loginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	switch {
	// Check we have some data
	case err == io.EOF:
		formatStandardResponse("data", "No request data supplied.", w)
		return nil
		// Check for parsing errors
	case err != nil:
		formatStandardResponse("decode-json", err.Error(), w)
		return nil
	}
	return &req
}

func loginValidate(req *loginRequest) error {
	if req.Email == "" || req.Password == "" || req.StoreID == "" || req.Series == "" {
		return fmt.Errorf("email, password, store and series must be entered")
	}
	if req.Series != "16" && req.Series != "18" {
		return fmt.Errorf("the series must be 16 or 18")
	}
	return nil
}
