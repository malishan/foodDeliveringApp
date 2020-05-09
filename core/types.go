package core

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// Routes - list of route
type Routes []route

type route struct {
	Name                   string
	MethodType             string
	Pattern                string
	ResourcesPermissionMap map[string]uint8
	HandlerFunc            http.HandlerFunc
}

//Response - complete structure of  HTTP response Meta + Data
type Response struct {
	Meta MetaData    `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

// MetaData of HTTP API response
type MetaData struct {
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
	RequestID     string `json:"requestId,omitempty"`
	CorrelationID string `json:"correlationId,omitempty"`
}

// Token forms the structure to validate user
type Token struct {
	UserID string
	*jwt.StandardClaims
}
