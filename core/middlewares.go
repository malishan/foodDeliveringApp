package core

import (
	"errors"
	"log"
	"net/http"
	"project/foodDeliveringApp/apicontext"
	"project/foodDeliveringApp/config"

	"github.com/dgrijalva/jwt-go"
)

//logRequest logs each HTTP incoming Requests
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, ": API called")
		next.ServeHTTP(w, r)
	})
}

// validateContext incoming request context
func validateContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		apiCtx, ok := ctx.Value(apicontext.APICtx).(apicontext.APIContext)

		tempContext := apicontext.AppContext{}
		tempContext.APIContext = apiCtx

		if !ok {
			HTTPErrorResponse(tempContext, w, config.BasicModule, ErrorUnauthorized, errors.New("unable to cast ctx to apiContext"))
			log.Println("context not found")
			return
		}

		if len(apiCtx.RestaurantID) == 0 {
			HTTPErrorResponse(tempContext, w, config.BasicModule, ErrorUnauthorized, errors.New(" missing restaurant id"))
			log.Println("context restaurantID not found")
			return
		}

		if len(apiCtx.BranchID) == 0 {
			HTTPErrorResponse(tempContext, w, config.BasicModule, ErrorUnauthorized, errors.New(" missing branch id"))
			log.Println("context branchID not found")
			return
		}

		if len(apiCtx.UserID) == 0 {
			HTTPErrorResponse(tempContext, w, config.BasicModule, ErrorUnauthorized, errors.New(" missing user id"))
			log.Println("context userID not found")
			return
		}

		if len(apiCtx.RoleID) == 0 {
			HTTPErrorResponse(tempContext, w, config.BasicModule, ErrorUnauthorized, errors.New(" missing role id"))
			log.Println("context roleID not found")
			return
		}

		// if len(apiCtx.Token) == 0 {
		// 	HTTPErrorResponse(tempContext, w, "config.BasicModule", ErrorUnauthorized, errors.New(" missing API token "))
		// 	log.Println("context token not found")
		// 	return
		// }

		ctx = apicontext.WithAPICtx(ctx, apiCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// createContext generates API context for incoming request and appends in request Context
func createResContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header
		ctx := r.Context()

		// requestID := header.Get(RequestID)
		// if requestID == "" {
		// 	requestID = uuid.NewV4().String()
		// }

		// correlationID := header.Get(CorrelationID)
		// if correlationID == "" {
		// 	correlationID = uuid.NewV4().String()
		// }

		restaurantID, branchID, userID, roleID := header.Get(RestaurantID), header.Get(BranchID), header.Get(UserID), header.Get(RoleID)
		//token := header.Get(AppAPIToken)

		apiCtx := apicontext.APIContext{
			RestaurantID: restaurantID,
			BranchID:     branchID,
			RoleID:       roleID,
			//Token:         token,
			//RequestID:     requestID,
			//CorrelationID: correlationID,
			UserID: userID,
		}

		ctx = apicontext.WithAPICtx(ctx, apiCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loginContextWrapper incoming request context
func loginContextWrapper(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		tempContext := apicontext.AppContext{}

		roleID := r.Header.Get(RoleID)
		if len(roleID) == 0 {
			log.Println("roleID header missing, err:", errors.New("missing header"))
			CustomHTTPErrorResponse(tempContext, w, config.BasicModule, ErrorMissingHeader, http.StatusBadRequest, "unaouthorized user access", 0, errors.New("missing roleID"))
			return
		}

		apiCtx := apicontext.APIContext{
			RoleID: roleID,
		}

		ctx = apicontext.WithAPICtx(ctx, apiCtx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateLoginToken validates FoodApp-api-token of incoming request
func validateUserLogin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tempContext := apicontext.AppContext{}
		apiCtx, ok := ctx.Value(apicontext.APICtx).(apicontext.APIContext)
		if !ok {
			HTTPErrorResponse(tempContext, w, config.BasicModule, ErrorUnauthorized, errors.New("unable to cast ctx to apiContext"))
			log.Println("context not found")
			return
		}

		tempContext.APIContext = apiCtx

		userID := r.Header.Get(UserID)
		if len(userID) == 0 {
			log.Println("userID header missing, err:", errors.New("missing header"))
			CustomHTTPErrorResponse(tempContext, w, config.BasicModule, ErrorMissingHeader, http.StatusBadRequest, "header for userID missing", 0, errors.New("missing userID"))
			return
		}

		tempContext.APIContext.UserID = userID

		ck, err := r.Cookie(CookieName)
		if err != nil {
			log.Println("user cookie missing, err:", errors.New("missing cookie"))
			CustomHTTPErrorResponse(tempContext, w, config.BasicModule, ErrorUnauthorized, http.StatusUnauthorized, "cookie not set", 0, errors.New("missing roleID"))
			return
		}

		token, err := jwt.ParseWithClaims(ck.Value, &Token{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		})

		if err != nil {
			log.Println("parsing with JWT failed, err:", errors.New("failed token parsing"))
			CustomHTTPErrorResponse(tempContext, w, config.BasicModule, DefaultErrorCode, http.StatusInternalServerError, "parsing with JWT failed", 0, err)
			return
		}

		if !token.Valid {
			log.Println("invalid authtoken, err:", errors.New("incorrect authtoken"))
			CustomHTTPErrorResponse(tempContext, w, config.BasicModule, ErrorUnauthorized, http.StatusUnauthorized, "incorrect authtoken", 0, errors.New("user not logged in"))
			return
		}

		//to-do: validate that the cookie value is similar to the userid value

		ctx = apicontext.WithAPICtx(ctx, tempContext.APIContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func verifySession(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		apiCtx, _ := ctx.Value(apicontext.APICtx).(apicontext.APIContext)

		tempContext := apicontext.AppContext{}
		tempContext.APIContext = apiCtx

		ck, err := r.Cookie(CookieName)
		if err != nil {
			log.Println("user cookie missing, err:", errors.New("missing cookie"))
			CustomHTTPErrorResponse(tempContext, w, config.BasicModule, ErrorUnauthorized, http.StatusUnauthorized, "cookie not set", 0, errors.New("missing roleID"))
			return
		}

		token, err := jwt.ParseWithClaims(ck.Value, &Token{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		})

		if err != nil {
			log.Println("parsing with JWT failed, err:", errors.New("failed token parsing"))
			CustomHTTPErrorResponse(tempContext, w, config.BasicModule, DefaultErrorCode, http.StatusInternalServerError, "parsing with JWT failed", 0, err)
			return
		}

		if !token.Valid {
			log.Println("invalid authtoken, err:", errors.New("incorrect authtoken"))
			CustomHTTPErrorResponse(tempContext, w, config.BasicModule, ErrorUnauthorized, http.StatusUnauthorized, "incorrect authtoken", 0, errors.New("user not logged in"))
			return
		}

		//to-do: validate that the cookie value is similar to the userid value

		ctx = apicontext.WithAPICtx(ctx, apiCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
