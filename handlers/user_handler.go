package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"project/foodDeliveringApp/apicontext"
	"project/foodDeliveringApp/config"
	"project/foodDeliveringApp/core"
	"project/foodDeliveringApp/mongodb"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// UserSignUp : handler for user registration
func UserSignUp(w http.ResponseWriter, r *http.Request) {

	ctx := apicontext.UpgradeCtx(r.Context())

	multipartError := r.ParseMultipartForm(32 << 20)
	if multipartError != nil {
		log.Println("parsing mulitpart form failed, err:", multipartError)
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, multipartError)
		return
	}

	var user UserInfo

	formInputs := r.FormValue("inputFields")
	if len(formInputs) == 0 {
		log.Println("no form values found for key-inputfields")
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, errors.New("no inputs for formValue field"))
		return
	}

	unmarshalErr := json.Unmarshal([]byte(formInputs), &user)
	if unmarshalErr != nil {
		log.Println("signUp formValue unmarshalling failed, err:", unmarshalErr)
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, unmarshalErr)
		return
	}

	multiPartFile, header, fileOpenErr := r.FormFile("image")
	if fileOpenErr != nil {
		log.Println("failed to open image file, err:", fileOpenErr)
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, fileOpenErr)
		return
	}

	defer multiPartFile.Close()

	index := strings.LastIndex(header.Filename, ".")

	extension := header.Filename[index+1:]
	if len(extension) == 0 {
		log.Println("file extension not found")
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, errors.New("file extension not found"))
		return
	}

	extension = strings.ToLower(extension)

	if !SupportedImageExt[extension] {
		log.Println("Invalid image (only JPEG or PNG supported)")
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, errors.New("incorrect file extension"))
		return
	}
	if extension == "jpeg" {
		extension = "jpg"
	}

	// newFileName := user.Name + "." + extension
	// originalImage, imgDecodeError := jpeg.Decode(multiPartFile)
	// if imgDecodeError != nil {
	// 	log.Println("image file decoding failed, err:", imgDecodeError)
	// 	core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, imgDecodeError)
	// 	return
	// }

	// upload image to s3 and retrieve the image URL. Update the user struct with the image path

	//to-do: validate signup credentials in front-end like email format, name length, etc

	var collection string
	if ctx.RoleID == apicontext.CustomerRole {
		collection = mongodb.CustomerCollection
	} else if ctx.RoleID == apicontext.AdminRole {
		collection = mongodb.AdminCollection
	}

	err := user.SignUp(collection)
	if err != nil {
		log.Println("user signUp failed, err:", err)
		if strings.Contains(err.Error(), "exists") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDuplicateUser, err)
		} else if strings.Contains(err.Error(), "internal server") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.DefaultServerError, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorQueryingDB, err)
		}
		return
	}

	log.Println("registration successful")
	core.HTTPResponse(ctx, w, http.StatusOK, "registration successfully", nil)
}

// UserLogin : handler for user login
func UserLogin(w http.ResponseWriter, r *http.Request) {

	ctx := apicontext.UpgradeCtx(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println("parse form error, err:", err)
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, err)
		return
	}

	formInput := r.FormValue("loginFields")

	var userLogin UserInfo

	unmarshalErr := json.Unmarshal([]byte(formInput), &userLogin)
	if unmarshalErr != nil {
		log.Println("signin formValue unmarshalling failed, err:", unmarshalErr)
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, unmarshalErr)
		return
	}

	id, err := userLogin.SignIn(ctx.APIContext.RoleID)
	if err != nil {
		log.Println("user signIn failed, err:", err)

		if strings.Contains(err.Error(), "not found") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDocumentNotFound, err)
		} else if strings.Contains(err.Error(), "incorrect") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorIncorrectCredentials, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.DefaultErrorCode, err)
		}

		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)

	tk := &core.Token{
		UserID: id,
		StandardClaims: &jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		log.Println("failed to create auth token, err:", err)
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.DefaultErrorCode, err)
		return
	}

	ctx.APIContext.UserID = id

	http.SetCookie(w, &http.Cookie{
		Name:    core.CookieName,
		Value:   tokenString,
		Path:    "/",
		MaxAge:  60 * 30,
		Expires: expirationTime, //it is ignored if max-age is supoorted
	})

	header := w.Header()
	SetLoginHeader(ctx, &header)

	log.Println("signIn successful")
	core.HTTPResponse(ctx, w, http.StatusOK, "login successfully", nil)
}

// UserUpdateAddress : handler for updating user address
func UserUpdateAddress(w http.ResponseWriter, r *http.Request) {

	ctx := apicontext.UpgradeCtx(r.Context())

	var updateUser struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		log.Println("update user address json decode failed, err:", err)
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, err)
		return
	}

	//to-do: check address is not blank in front end

	usr := UserInfo{
		UserID:  ctx.APIContext.UserID,
		Address: updateUser.Address,
	}

	err := usr.ChangeAddress(ctx.APIContext.RoleID)
	if err != nil {
		log.Println("user update address failed, err:", err)

		if strings.Contains(err.Error(), "not found") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDocumentNotFound, err)
		} else if strings.Contains(err.Error(), "not authentic") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorUnauthorized, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorQueryingDB, err)
		}

		return
	}

	header := w.Header()
	SetLoginHeader(ctx, &header)

	log.Println("user address update successful")
	core.HTTPResponse(ctx, w, http.StatusOK, "user address update successful", nil)
}

// UserUpdatePhone : handler for updating user phone number
func UserUpdatePhone(w http.ResponseWriter, r *http.Request) {

	ctx := apicontext.UpgradeCtx(r.Context())

	var updateUser struct {
		Phone []string `json:"phoneNo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		log.Println("update user phone json decode failed, err:", err)
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, err)
		return
	}

	//to-do: check phone number is not blank in front end

	usr := UserInfo{
		UserID: ctx.APIContext.UserID,
		Phone:  updateUser.Phone,
	}

	err := usr.ChangePhoneNo(ctx.APIContext.RoleID)
	if err != nil {
		log.Println("user update phone number failed, err:", err)

		if strings.Contains(err.Error(), "not found") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDocumentNotFound, err)
		} else if strings.Contains(err.Error(), "not authentic") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorUnauthorized, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorQueryingDB, err)
		}

		return
	}

	header := w.Header()
	SetLoginHeader(ctx, &header)

	log.Println("user phoneNo update successful")
	core.HTTPResponse(ctx, w, http.StatusOK, "user phoneNo update successful", nil)
}

// UserUpdateProfile : handler for updating user profile pic
func UserUpdateProfile(w http.ResponseWriter, r *http.Request) {

	ctx := apicontext.UpgradeCtx(r.Context())

	var updateUser struct {
		ProfilePic string `json:"profilePic"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		log.Println("update user profilePic json decode failed, err:", err)
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, err)
		return
	}

	//to-do: check profilePic is not blank in front end

	usr := UserInfo{
		UserID:     ctx.APIContext.UserID,
		ProfilePic: updateUser.ProfilePic,
	}

	err := usr.ChangeProfilePic(ctx.APIContext.RoleID)
	if err != nil {
		log.Println("user update profile pic failed, err:", err)

		if strings.Contains(err.Error(), "not found") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDocumentNotFound, err)
		} else if strings.Contains(err.Error(), "not authentic") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorUnauthorized, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorQueryingDB, err)
		}

		return
	}

	header := w.Header()
	SetLoginHeader(ctx, &header)

	log.Println("user profilePic update successful")
	core.HTTPResponse(ctx, w, http.StatusOK, "user profilePic update successful", nil)
}

// UserUpdatePassword : handler for updating user password
func UserUpdatePassword(w http.ResponseWriter, r *http.Request) {

	ctx := apicontext.UpgradeCtx(r.Context())

	var pswdUpdate struct {
		Password string `json:"pswd"`
	}

	if err := json.NewDecoder(r.Body).Decode(&pswdUpdate); err != nil {
		log.Println("update user password json decode failed, err:", err)
		core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDecodingPayload, err)
		return
	}

	//to-do: check password is not blank in front end

	usr := UserInfo{
		UserID:   ctx.APIContext.UserID,
		Password: pswdUpdate.Password,
	}

	err := usr.ChangePassword(ctx.APIContext.RoleID)
	if err != nil {
		log.Println("user password update failed, err:", err)

		if strings.Contains(err.Error(), "not found") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDocumentNotFound, err)
		} else if strings.Contains(err.Error(), "not authentic") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorUnauthorized, err)
		} else if strings.Contains(err.Error(), "internal server") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.DefaultServerError, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorQueryingDB, err)
		}

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    core.CookieName,
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0), //it is ignored if max-age is supported
	})

	log.Println("user password updated successful")
	core.HTTPResponse(ctx, w, http.StatusOK, "user password updated successful, please login", nil)
}

// UserLogout : handler for user logout
func UserLogout(w http.ResponseWriter, r *http.Request) {

	ctx := apicontext.UpgradeCtx(r.Context())

	usr := UserInfo{
		UserID: ctx.APIContext.UserID,
	}

	err := usr.Signout(ctx.APIContext.RoleID, ctx.APIContext.UserID)
	if err != nil {
		log.Println("user logout failed, err:", err)

		if strings.Contains(err.Error(), "not found") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorDocumentNotFound, err)
		} else if strings.Contains(err.Error(), "not authentic") {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorUnauthorized, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.LoginModule, core.ErrorQueryingDB, err)
		}

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    core.CookieName,
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0), //it is ignored if max-age is supported
	})

	log.Println("user signOut successful")
	core.HTTPResponse(ctx, w, http.StatusOK, "user signout successful", nil)
	//to-do: http.Redirect()
}

// SetLoginHeader adds new header with roleID and userID to the HTTP request
func SetLoginHeader(ctx apicontext.AppContext, h *http.Header) {
	h.Set(core.RoleID, ctx.APIContext.RoleID)
	h.Set(core.UserID, ctx.APIContext.UserID)
}
