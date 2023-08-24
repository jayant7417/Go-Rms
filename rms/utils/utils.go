package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"math/big"
	"net/http"
	"strings"
)

type clientError struct {
	ID            string `json:"id"`
	MessageToUser string `json:"messageToUser"`
	DeveloperInfo string `json:"developerInfo"`
	Err           string `json:"error"`
	StatusCode    int    `json:"statusCode"`
	IsClientError bool   `json:"isClientError"`
}

var generator *shortid.Shortid

const generatorSeed = 1000

func init() {
	n, err := rand.Int(rand.Reader, big.NewInt(generatorSeed))
	if err != nil {
		logrus.Panicf("failed to initialize utilities with random seed, %+v", err)
		return
	}

	g, err := shortid.New(1, shortid.DefaultABC, n.Uint64())

	if err != nil {
		logrus.Panicf("Failed to initialize utils package with error: %+v", err)
	}

	generator = g
}

func ParseBody(body io.Reader, out interface{}) error {
	err := json.NewDecoder(body).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

func EncodeJSONBody(w http.ResponseWriter, body interface{}) interface{} {
	return json.NewEncoder(w).Encode(body)
}

func RespondJSON(w http.ResponseWriter, StatusCode int, body interface{}) {
	w.WriteHeader(StatusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			logrus.Errorf("Failed to respond JSON with error: #{err}")
		}
	}
}

func RespondError(w http.ResponseWriter, statusCode int, err error, messageToUser string, additionalInfoForDevs ...string) {
	logrus.Errorf("status: #{statusCode}, message: #{messageToUser},err: #{err}")
	clientError := NewClientError(err, statusCode, messageToUser, additionalInfoForDevs...)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientError); err != nil {
		logrus.Errorf("Failed to send to caller with error: #{err}")
	}
}

func NewClientError(err error, statusCode int, messageToUser string, additionalInfoForDevs ...string) *clientError {
	additionalInfoJoined := strings.Join(additionalInfoForDevs, "\n")
	if additionalInfoJoined == "" {
		additionalInfoJoined = messageToUser
	}

	errorID, _ := generator.Generate()
	var errString string
	if err != nil {
		errString = err.Error()
	}
	return &clientError{
		ID:            errorID,
		MessageToUser: messageToUser,
		DeveloperInfo: additionalInfoJoined,
		Err:           errString,
		StatusCode:    statusCode,
		IsClientError: true,
	}
}

// HashPassword return hashed password and error
func HashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(pass), nil
}

// CheckPassword returns nil on success, or an error on failure.
func CheckPassword(pass, hashedPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
}
