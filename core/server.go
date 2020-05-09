package core

import (
	"log"
	"net/http"
	"project/foodDeliveringApp/apicontext"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

const (
	CookieName    = "foodApp-access-token"
	RestaurantID  = "restaurantID"
	BranchID      = "branchID"
	RoleID        = "roleID"
	UserID        = "userID"
	TransactionID = "transactionID"
	TimeZone      = "timezone"
	Locale        = "locale"
	RequestID     = "requestId"
	CorrelationID = "correlationId"
	timeFormat    = time.RFC3339
)

var routes = make(Routes, 0)

//StartServer - http servers
func StartServer(port, subroute string) {

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	allowedHeaders := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"X-CSRF-Token",
		"X-Auth-Token",
		"Content-Type",
		"processData",
		"contentType",
		"Origin",
		"Authorization",
		"Accept",
		"Client-Security-Token",
		"Accept-Encoding",
		TimeZone,
		Locale,
		RestaurantID,
		BranchID,
		RoleID,
		RequestID,
		CorrelationID,
		UserID,
	})

	allowedMethods := handlers.AllowedMethods([]string{
		"POST",
		"GET",
		"DELETE",
		"PUT",
		"PATCH",
		"OPTIONS"})

	allowCredential := handlers.AllowCredentials()

	handlers := handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins, allowCredential)(context.ClearHandler(newRouter(subroute)))

	log.Fatal(http.ListenAndServe(":"+port, handlers))
}

// NewRouter provides a mux Router.
// Handles all incoming request who matches registered routes against the request.
func newRouter(subroute string) *mux.Router {
	muxRouter := mux.NewRouter().StrictSlash(true)
	subRouter := muxRouter.PathPrefix(subroute).Subrouter()
	for _, route := range routes {
		subRouter.
			Methods(route.MethodType).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return muxRouter
}

// useMiddleware applies chains of middleware (ie: log, contextWrapper, validateAuth) handler into incoming request
// For example, logging middleware might write the incoming request details to a log
// Note - It applies in reverse order
func useMiddleware(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

// Ping API used to check the status of the service
func Ping(response http.ResponseWriter, request *http.Request) {
	HTTPPingResponse(response, http.StatusOK, map[string]string{"ping": "pong"})
}

// HTTPPingResponse writes the HTTPResponse and renders the json: Uses context
func HTTPPingResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	renderer := render.New()
	renderer.JSON(w, statusCode, data)
}

// HTTPErrorResponse writes the HTTPResponse and renders the json
func HTTPErrorResponse(ctx apicontext.AppContext, w http.ResponseWriter, moduleID string, errorCode string, err error) {
	renderer := render.New()
	requestID := ctx.RequestID
	correlationID := ctx.CorrelationID
	res := Response{}
	erratumErr := ConstructErrorObjectWithContext(ctx, moduleID, errorCode, err)
	res.Meta.Code = erratumErr.HTTPErrCode
	res.Meta.Msg = erratumErr.UserMessage
	res.Meta.RequestID = requestID
	res.Meta.CorrelationID = correlationID
	renderer.JSON(w, erratumErr.HTTPErrCode, res)
}

// HTTPResponse writes the HTTPResponse and renders the json: Uses context
func HTTPResponse(ctx apicontext.AppContext, w http.ResponseWriter, statusCode int, msg string, data interface{}) {
	renderer := render.New()
	requestID := ctx.RequestID
	correlationID := ctx.CorrelationID

	res := Response{}
	res.Meta.Code = statusCode
	res.Meta.Msg = msg
	res.Meta.RequestID = requestID
	res.Meta.CorrelationID = correlationID
	res.Data = data
	renderer.JSON(w, statusCode, res)
}

// ConstructErrorObjectWithContext constructs the error object structure based on the moduleID and the erratum code
func ConstructErrorObjectWithContext(ctx apicontext.AppContext, moduleID, errorCode string, innerError error) ErrorMessage {
	errorObject, ok := Errordefs[errorCode]
	if !ok {
		errorObject = Errordefs[DefaultErrorCode]
	}
	if innerError != nil {
		errorObject.LogMessage = innerError.Error()
	}
	errorObject.ErrCode = moduleID + errorCode
	return errorObject
}

// CustomHTTPErrorResponse writes the HTTPResponse and renders the json.
// This can be used to mask the internal log messages from being sent to requester and instead send a custom code and message.
// ctx - request context
// customMessage - string, if specified, would override the erratum's error message.
// httpErrorCode - int, if specified, would override the erratum's HTTP error code
// customErrorCode - int, if specified, would override the Meta.Code. Can be used for internal error mapping.
func CustomHTTPErrorResponse(ctx apicontext.AppContext, w http.ResponseWriter, moduleID string, errorCode string, httpErrorCode int, customMessage string, customErrorCode int, err error) {
	renderer := render.New()
	requestID := ctx.RequestID
	correlationID := ctx.CorrelationID
	res := Response{}
	erratumErr := ConstructCustomErrorMessageWithContext(ctx, moduleID, errorCode, err, customMessage, httpErrorCode)
	if customErrorCode != 0 {
		res.Meta.Code = customErrorCode
	} else {
		res.Meta.Code = erratumErr.HTTPErrCode
	}
	res.Meta.Msg = erratumErr.UserMessage
	res.Meta.RequestID = requestID
	res.Meta.CorrelationID = correlationID
	renderer.JSON(w, erratumErr.HTTPErrCode, res)
}

// ConstructCustomErrorMessageWithContext constructs the error message structure and assigns the actual
// error message in to the LogMessage and user supplied message
func ConstructCustomErrorMessageWithContext(ctx apicontext.AppContext, moduleID, errorCode string, innerError error, userMessage string, httpErrorCode int) ErrorMessage {
	errorObject, ok := Errordefs[errorCode]
	if !ok {
		// log.GenericError(ctx, errors.New("looks like there is an issue with the error code used"), nil)
		errorObject = Errordefs[DefaultErrorCode]
	}
	if innerError != nil {
		errorObject.LogMessage = innerError.Error()
	}
	if userMessage != "" {
		errorObject.UserMessage = userMessage
	}
	if httpErrorCode != 0 {
		errorObject.HTTPErrCode = httpErrorCode
	}
	errorObject.ErrCode = moduleID + errorCode
	//Log the error message
	//log.GenericError(ctx, errors.New("|"+errorObject.ErrCode+"|"+strconv.Itoa(errorObject.HTTPErrCode)+"|"+errorObject.LogMessage+"|"+errorObject.UserMessage), nil)
	return errorObject
}

// AddRoute is to create routes with auth, user and restaurant details
func AddRoute(methodName, methodType, mRoute string, m map[string]uint8, handlerFunc http.HandlerFunc) {
	r := route{
		Name:                   methodName,
		MethodType:             methodType,
		Pattern:                mRoute,
		ResourcesPermissionMap: m,
		HandlerFunc:            useMiddleware(handlerFunc, verifySession ,validateContext, logRequest, createResContext),
	}

	routes = append(routes, r)
}

//AddLoginRoutes is to create routes without authentication check
func AddLoginRoutes(methodName string, methodType string, mRoute string, m map[string]uint8, handlerFunc http.HandlerFunc) {
	r := route{
		Name:                   methodName,
		MethodType:             methodType,
		Pattern:                mRoute,
		ResourcesPermissionMap: m,
		HandlerFunc:            useMiddleware(handlerFunc, loginContextWrapper, logRequest)}
	routes = append(routes, r)

}

// AddUserRoute is to create routes with user authentication only
func AddUserRoute(methodName string, methodType string, mRoute string, m map[string]uint8, handlerFunc http.HandlerFunc) {
	r := route{
		Name:                   methodName,
		MethodType:             methodType,
		Pattern:                mRoute,
		ResourcesPermissionMap: m,
		HandlerFunc:            useMiddleware(handlerFunc, validateUserLogin, loginContextWrapper, logRequest)}
	routes = append(routes, r)
}

//AddNoAuthRoutes - Route without any Auth
func AddNoAuthRoutes(methodName string, methodType string, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        methodName,
		MethodType:  methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, logRequest)}
	routes = append(routes, r)

}
