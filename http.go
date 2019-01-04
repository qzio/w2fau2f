package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"u2f"
)

var challenge *u2f.Challenge
var registrations []u2f.Registration
var counter uint32

type httpMessage struct {
	Message   string
	ErrorCode int
}

func errorResponse(w http.ResponseWriter, msg string, errorCode int) {
	log.Printf("failed: %v", msg)
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(httpMessage{Message: msg, ErrorCode: errorCode})
}

func Register(w http.ResponseWriter, r *http.Request) {
	log.Println("doing register")
	if r.Method == "GET" {
		u2fRegReq, err := NewU2FRegReq(r.RemoteAddr)
		if err != nil {
			errorResponse(w, "failed to create new registration", http.StatusInternalServerError)
			return
		}
		log.Printf("register request: %v", u2fRegReq)
		json.NewEncoder(w).Encode(u2fRegReq)
		return

	} else if r.Method == "POST" {
		log.Println("Doing register.post")
		var regResp u2f.RegisterResponse
		if err := json.NewDecoder(r.Body).Decode(&regResp); err != nil {
			log.Printf("failed to decode req body %v", err)
			http.Error(w, "error "+err.Error(), http.StatusBadRequest)
		}

		if err := CompleteRegReq(r.RemoteAddr, regResp); err != nil {
			var errorCode int
			if strings.Contains(err.Error(), "challenge not found") {
				errorCode = http.StatusBadRequest
			} else {
				errorCode = http.StatusInternalServerError
			}
			errorResponse(w, err.Error(), errorCode)
			return
		}
		log.Printf("Registration success: %+v", regResp)
		msg := httpMessage{Message: "registration success"}
		json.NewEncoder(w).Encode(msg)
		return
	} else {
		errorResponse(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("doing login")
	if r.Method == "GET" {
		signReq, err := NewSignReq(r.RemoteAddr)
		if err != nil {
			log.Printf("failed to parse json %v", err)
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("sign request: %+v", signReq)
		json.NewEncoder(w).Encode(signReq)
		return
	} else if r.Method == "POST" {
		var signResp u2f.SignResponse
		if err := json.NewDecoder(r.Body).Decode(&signResp); err != nil {
			log.Printf("failed to parse json %v", err)
			errorResponse(w, "Failed to parse json body", http.StatusBadRequest)
			return
		}

		log.Printf("signResponse: %+v", signResp)
		err := CompleteSignReq(r.RemoteAddr, signResp)
		if err != nil {
			errorResponse(w, "failure!", http.StatusBadRequest)
			return
		} else {
			msg := httpMessage{Message: "Signin success"}
			json.NewEncoder(w).Encode(msg)
			return
		}
	} else {
		errorResponse(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}
