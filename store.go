package main

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"

	"github.com/tstranex/u2f"
)

type miniReg struct {
	Raw       []byte
	KeyHandle []byte
	PubKey    struct {
		X big.Int
		Y big.Int
	}
}

// SaveRegistration saves registration to file
func SaveRegistration(remoteAddr string, r u2f.Registration) error {
	if err := write("registration", remoteAddr, &r); err != nil {
		return err
	}
	return nil
}

// LoadRegistration loads a registration based on an ip
// It is not possible to just unmarshal the elliptic.Curve
// object, so we jump through some hoops for that.
func LoadRegistration(remoteAddr string) (*u2f.Registration, error) {
	mreg := miniReg{}
	if err := load("registration", remoteAddr, &mreg); err != nil {
		return nil, err
	}
	var pubk = ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &mreg.PubKey.X,
		Y:     &mreg.PubKey.Y,
	}
	var reg = u2f.Registration{
		Raw:       mreg.Raw,
		KeyHandle: mreg.KeyHandle,
		PubKey:    pubk,
	}
	return &reg, nil
}

// SaveChallenge saves a challenge to file
func SaveChallenge(remoteAddr string, c *u2f.Challenge) error {
	log.Printf("saving challenge to %v", remoteAddr)
	if err := write("challenge", remoteAddr, c); err != nil {
		return err
	}
	return nil
}

// LoadChallenge loads a registration based on an ip
func LoadChallenge(remoteAddr string) (*u2f.Challenge, error) {
	var chl = u2f.Challenge{}
	if err := load("challenge", remoteAddr, &chl); err != nil {
		log.Printf("Failed to load challenge for ip %v", remoteAddr)
		return nil, err
	}
	log.Printf("challenge; %v", chl)
	return &chl, nil
}

func write(t string, remoteAddr string, val interface{}) error {
	if err := ensureDir(t); err != nil {
		return err
	}
	sIp, err := safeIp(remoteAddr)
	if err != nil {
		return errors.New("Invalid ip")
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.Encode(&val)
	filePath := fmt.Sprintf("%v/%v", t, sIp)
	log.Printf("will write %v to %v", buf, filePath)

	fd, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file %v", filePath)
		return err
	}
	defer fd.Close()

	w := bufio.NewWriter(fd)
	writtenBytes, err := buf.WriteTo(w)
	w.Flush()
	if err != nil {
		log.Printf("failed to write file %v", filePath)
		return err
	}
	log.Printf("Success writing to file %v", writtenBytes)
	return nil
}

func load(t string, remoteAddr string, val interface{}) error {
	if err := ensureDir(t); err != nil {
		return err
	}
	sIp, err := safeIp(remoteAddr)
	if err != nil {
		return err
	}
	filePath := fmt.Sprintf("%v/%v", t, sIp)
	fd, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file: %v", filePath)
		return err
	}
	defer fd.Close()

	ioR := bufio.NewReader(fd)

	decoder := json.NewDecoder(ioR)
	err = decoder.Decode(&val)
	if err != nil {
		log.Printf("Failed to decode json from %v, %v", filePath, err)
		return err
	}
	return nil
}

func safeIp(remoteAddr string) (string, error) {
	splitted, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to split port from ip %v", remoteAddr))
	}
	ipObj := net.ParseIP(splitted)
	if ipObj == nil {
		log.Printf("Failed to parse ip: %v", remoteAddr)
		return "", errors.New("Failed to parse ip")
	}
	return strings.Replace(ipObj.String(), ":", "x", -1), nil

}

func ensureDir(t string) error {
	path := t
	if _, err := os.Stat(path); os.IsNotExist(err) {
		osErr := os.Mkdir(path, 0755)
		if osErr != nil {
			log.Printf("failed to create dir %v due to %v", path, osErr)
			return osErr
		}
	}
	return nil
}
