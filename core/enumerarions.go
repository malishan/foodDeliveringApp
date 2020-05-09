package core

import "net/http"

const (
	RestaurantService = "restaurantservice"
	LoginService      = "loginservice"
)

const (
	// DefaultErrorCode for unknown errors
	DefaultErrorCode = "001"
	// ErrorDecodingPayload error code for decoding payload
	ErrorDecodingPayload = "002"
	// ErrorDecodingData error code for decoding data
	ErrorDecodingData = "003"
	// ErrorWithQueryParams error while sending email
	ErrorWithQueryParams = "004"
	// ErrorInvalidID error code for invalid ID
	ErrorInvalidID = "005"
	// ErrorInvalidRedisToken error code for invalid ID
	ErrorInvalidRedisToken = "006"
	// ErrorEncodingResponse error code for encoding response
	ErrorEncodingResponse = "007"
	// ErrorConnectingWebsocket error code for connecting websocket
	//ErrorConnectingWebsocket = "008"
	// ErrorSendingToMessagebroker error code for sending to message broker
	//ErrorSendingToMessagebroker = "009"
	// ErrorQueryingDB error code for querying database
	ErrorQueryingDB = "010"
	// ErrorInsertingMongoDoc error code for executing database command
	ErrorInsertingMongoDoc = "011"
	// ErrorDeletingMongoDoc error code for executing database command
	ErrorDeletingMongoDoc = "012"
	// ErrorUpdatingMongoDoc error code for executing database command
	ErrorUpdatingMongoDoc = "013"
	// ErrorGettingRedisCache error code for getting redis cache data
	ErrorGettingRedisCache = "014"
	// ErrorSettingRedisCache error code for setting redis cache data
	ErrorSettingRedisCache = "015"
	// ErrorUnauthorized error code when auth token is invalid or missing
	ErrorUnauthorized = "018"
	// ErrorUnsuccessfulAttempt error is when user login crosses limit or by providing invalid credentials
	ErrorUnsuccessfulAttempt = "019"
	// ErrorSendingEmail error while sending email
	ErrorSendingEmail = "020"
	// ErrorReadingVault = "023"
	//ErrorReadingVault = "021"
	//ErrorDocumentNotFound document not found
	ErrorDocumentNotFound = "022"
	// ErrorUploadingFile - during S3 media upload
	ErrorUploadingFile = "023"
	// ErrorDownloadingFile - during S3 media download
	ErrorDownloadingFile = "024"
	// ErrorCreatingBucket - during S3 bucket creation
	ErrorCreatingBucket = "025"
	// ErrorIncorrectCredentials - error when incorrect credentials are passed
	ErrorIncorrectCredentials = "026"
	// ErrorIncorrectEmail - error sent when only email is validated
	ErrorIncorrectEmail = "027"
	// DefaultServerError - default error which can be used to indicate internal server error for any failed internal operations
	DefaultServerError = "028"
	//Already doc exists
	ErrorDocumentExists      = "029"
	ErrorReqValidationFailed = "030"
	ErrorInsertingMySqlDoc   = "031"
	// ErrorDeletingMongoDoc error code for executing database command
	ErrorDeletingMySqlDoc = "032"
	// ErrorUpdatingMongoDoc error code for executing database command
	ErrorUpdatingMySqlDoc = "033"
	//ErrorNoContent No content available
	ErrorNoContent     = "034"
	ErrorMissingHeader = "035"
	ErrorDuplicateUser = "036"
)

// Errordefs definition of error
var Errordefs = map[string]ErrorMessage{

	// default error
	DefaultErrorCode:   {DefaultErrorCode, http.StatusNotFound, "Server Error", "Sorry! An error occurred, please try again. If the problem persists, please contact support."},
	DefaultServerError: {DefaultServerError, http.StatusInternalServerError, "Server Error", "An error occurred on our system. Please contact the support team"},

	// json related error
	ErrorDecodingPayload:  {ErrorDecodingPayload, http.StatusBadRequest, "", "Invalid request"},
	ErrorDecodingData:     {ErrorDecodingData, http.StatusBadRequest, "", "Unrecognized data"},
	ErrorEncodingResponse: {ErrorEncodingResponse, http.StatusBadRequest, "", "Cannot encode data"},

	ErrorWithQueryParams: {ErrorWithQueryParams, http.StatusBadRequest, "", "Invalid request"},
	ErrorMissingHeader:   {ErrorMissingHeader, http.StatusBadGateway, "", "Certain headers missing"},
	ErrorInvalidID:       {ErrorInvalidID, http.StatusBadRequest, "", "Invalid ID"},

	//database related error
	ErrorInvalidRedisToken: {ErrorInvalidRedisToken, http.StatusInternalServerError, "", "Looks like link was invalid, possibly it has already been used"},
	ErrorGettingRedisCache: {ErrorGettingRedisCache, http.StatusBadRequest, "", "Cache error"},
	ErrorSettingRedisCache: {ErrorSettingRedisCache, http.StatusInternalServerError, "", "Cache error"},

	ErrorQueryingDB:        {ErrorQueryingDB, http.StatusBadRequest, "", "Database error"},
	ErrorInsertingMongoDoc: {ErrorInsertingMongoDoc, http.StatusBadRequest, "", "Database error"},
	ErrorDeletingMongoDoc:  {ErrorDeletingMongoDoc, http.StatusBadRequest, "", "Database error"},
	ErrorUpdatingMongoDoc:  {ErrorUpdatingMongoDoc, http.StatusBadRequest, "", "Database error"},

	ErrorInsertingMySqlDoc: {ErrorInsertingMongoDoc, http.StatusBadRequest, "", "Database error"},
	ErrorDeletingMySqlDoc:  {ErrorDeletingMySqlDoc, http.StatusBadRequest, "", "Database error"},
	ErrorUpdatingMySqlDoc:  {ErrorUpdatingMySqlDoc, http.StatusBadRequest, "", "Database error"},

	// duplicate or not present
	ErrorNoContent:        {ErrorNoContent, http.StatusNoContent, "", "No content available"},
	ErrorDocumentNotFound: {ErrorDocumentNotFound, http.StatusNotFound, "", "Document not found"},
	ErrorDocumentExists:   {ErrorDocumentExists, http.StatusNotAcceptable, "", "Document already exists"},
	ErrorDuplicateUser:    {ErrorDuplicateUser, http.StatusBadRequest, "user already exists", "user exists"},

	// s3 related error
	ErrorCreatingBucket:  {ErrorCreatingBucket, http.StatusInternalServerError, "", "Sorry! Could not create bucket"},
	ErrorUploadingFile:   {ErrorUploadingFile, http.StatusInternalServerError, "", "Sorry! Could not upload file"},
	ErrorDownloadingFile: {ErrorDownloadingFile, http.StatusInternalServerError, "", "Sorry! Could not generate signed URL for download file"},

	// user access related error
	ErrorUnauthorized:         {ErrorUnauthorized, http.StatusForbidden, "", "Unauthorized request"},
	ErrorUnsuccessfulAttempt:  {ErrorUnsuccessfulAttempt, http.StatusBadRequest, "", "Unsuccessful attempt"},
	ErrorReqValidationFailed:  {ErrorReqValidationFailed, http.StatusBadRequest, "", "Request validation failed"},
	ErrorIncorrectCredentials: {ErrorIncorrectCredentials, http.StatusBadRequest, "", "Incorrect username or password"},

	// email related issue
	ErrorIncorrectEmail: {ErrorIncorrectEmail, http.StatusBadRequest, "", "Incorrect email"},
	ErrorSendingEmail:   {ErrorSendingEmail, http.StatusBadRequest, "", "Email error"},
}

// ErrorMessage is an internal structure containing the application or system error code along with the http status code
type ErrorMessage struct {
	ErrCode     string `json:"errCode,omitempty"`
	HTTPErrCode int    `json:"httpErrCode,omitempty"`
	LogMessage  string `json:"logMessage,omitempty"`  //message to be logged in our logs (err.Error string)
	UserMessage string `json:"userMessage,omitempty"` //message to be displayed to the user
}
