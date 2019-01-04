package main

import (
	"errors"
	"log"
	"u2f"
)

// NewU2FRegReq returns a new registration request struct
func NewU2FRegReq(ip string) (*u2f.WebRegisterRequest, error) {
	chl, err := u2f.NewChallenge(appID, trustedFacets)
	if err != nil {
		log.Printf("Failed to create new challenge: %v", err)
		return nil, err
	}

	// store challenge to be able to save a completed registration request challenge response
	SaveChallenge(ip, chl)
	u2fReq := u2f.NewWebRegisterRequest(chl, registrations)
	log.Printf("registerRequest: %+v", u2fReq)
	return u2fReq, nil
}

// CompleteRegReq completes a registration
func CompleteRegReq(ip string, regResp u2f.RegisterResponse) error {
	log.Printf("decoded json: %v \n", regResp)
	challenge, err := LoadChallenge(ip)
	if err != nil {
		log.Printf("Failed to load challenge for ip %s", ip)
		return errors.New("challenge not found")
	}

	config := &u2f.Config{
		// Chrome 66+ doesn't return the device's attestation
		// certificate by default.
		SkipAttestationVerify: true,
	}

	reg, err := u2f.Register(regResp, *challenge, config)
	if err != nil {
		log.Printf("u2f.Register error: %v", err)
		return errors.New("error verifyfing response")
	}

	if err := SaveRegistration(ip, *reg); err != nil {
		return errors.New("Failed to save registration")
	}

	return nil
}

func NewSignReq(remoteAddr string) (*u2f.WebSignRequest, error) {

	reg, err := LoadRegistration(remoteAddr)
	if err != nil {
		return nil, err
	}
	chl, err := u2f.NewChallenge(appID, trustedFacets)
	if err != nil {
		log.Printf("u2f.NewChallenge error: %v", err)
		return nil, err
	}

	if err := SaveChallenge(remoteAddr, chl); err != nil {
		return nil, err
	}
	var regs = []u2f.Registration{*reg}
	signReq := chl.SignRequest(regs)
	return signReq, nil
}

// CompleteSignReq completes a signature (login) and returns error upon failure
// @TODO counter/newCounter needs to be handled
func CompleteSignReq(remoteAddr string, sig u2f.SignResponse) error {
	challenge, err := LoadChallenge(remoteAddr)
	if err != nil {
		return err
	}
	reg, err := LoadRegistration(remoteAddr)
	if err != nil {
		return err
	}
	var counter uint32
	counter = 0
	newCounter, authErr := reg.Authenticate(sig, *challenge, counter)
	if authErr != nil {
		log.Printf("Failed to authenticate!: %v", authErr)
		return authErr
	}
	log.Printf("new counter is %v", newCounter)
	if err := OpenAccess(remoteAddr); err != nil {
		return err
	}
	return nil
}
