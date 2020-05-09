package handlers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"project/foodDeliveringApp/apicontext"
	"project/foodDeliveringApp/mongodb"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// SignUp checks and stores customer info
func (usr UserInfo) SignUp(collection string) error {

	err := mongodb.Exists(collection, bson.M{"email": usr.Email}, nil)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	} else if err == nil {
		return errors.New("email already exists")
	}

	err = mongodb.Exists(collection, bson.M{"name": usr.Name}, nil)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	} else if err == nil {
		return errors.New("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 5)
	if err != nil {
		return fmt.Errorf("internal server error, %v", err)
	}

	usr.Password = string(hash)

	err = mongodb.Insert(collection, usr)
	return err
}

// SignIn validates user credentials
func (usr UserInfo) SignIn(roleID string) (id string, err error) {

	var rsltUser struct {
		UserID   []byte `bson:"_id"`
		Password string `bson:"pswd"`
	}

	var query bson.M

	if len(usr.Email) > 0 {
		query = bson.M{"email": usr.Email}
	} else {
		query = bson.M{"name": usr.Name}
	}

	if roleID == apicontext.CustomerRole {
		err = mongodb.ReadOne(mongodb.CustomerCollection, query, bson.M{"_id": 1, "pswd": 1}, &rsltUser)
	} else if roleID == apicontext.AdminRole {
		err = mongodb.ReadOne(mongodb.AdminCollection, query, bson.M{"_id": 1, "pswd": 1}, &rsltUser)
	} else {
		err = errors.New("user role not authentic")
	}

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(rsltUser.Password), []byte(usr.Password))
	if err != nil {
		return "", errors.New("incorrect password")
	}

	id = hex.EncodeToString(rsltUser.UserID)
	return id, nil
}

// ChangeAddress checks user role and updates user address to database
func (usr UserInfo) ChangeAddress(roleID string) (err error) {

	if roleID == apicontext.CustomerRole {
		err = mongodb.Update(mongodb.CustomerCollection, bson.M{"_id": bson.ObjectIdHex(usr.UserID)}, bson.M{"$set": bson.M{"address": usr.Address}})
	} else if roleID == apicontext.AdminRole {
		err = mongodb.Update(mongodb.AdminCollection, bson.M{"_id": bson.ObjectIdHex(usr.UserID)}, bson.M{"$set": bson.M{"address": usr.Address}})
	} else {
		err = errors.New("user role not authentic")
	}

	return
}

// ChangePhoneNo checks user role and updates user phone number  to database
func (usr UserInfo) ChangePhoneNo(roleID string) (err error) {

	if roleID == apicontext.CustomerRole {
		err = mongodb.Update(mongodb.CustomerCollection, bson.M{"_id": bson.ObjectIdHex(usr.UserID)}, bson.M{"$set": bson.M{"phoneNo": usr.Phone}})
	} else if roleID == apicontext.AdminRole {
		err = mongodb.Update(mongodb.AdminCollection, bson.M{"_id": bson.ObjectIdHex(usr.UserID)}, bson.M{"$set": bson.M{"phoneNo": usr.Phone}})
	} else {
		err = errors.New("user role not authentic")
	}

	return
}

// ChangeProfilePic checks user role and updates user profile pic to database
func (usr UserInfo) ChangeProfilePic(roleID string) (err error) {

	if roleID == apicontext.CustomerRole {
		err = mongodb.Update(mongodb.CustomerCollection, bson.M{"_id": bson.ObjectIdHex(usr.UserID)}, bson.M{"$set": bson.M{"profilePic": usr.ProfilePic}})
	} else if roleID == apicontext.AdminRole {
		err = mongodb.Update(mongodb.AdminCollection, bson.M{"_id": bson.ObjectIdHex(usr.UserID)}, bson.M{"$set": bson.M{"profilePic": usr.ProfilePic}})
	} else {
		err = errors.New("user role not authentic")
	}

	return
}

// ChangePassword checks user role and updates user password to database
func (usr UserInfo) ChangePassword(roleID string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 5)
	if err != nil {
		return fmt.Errorf("internal server error, %v", err)
	}

	usr.Password = string(hash)

	if roleID == apicontext.CustomerRole {
		err = mongodb.Update(mongodb.CustomerCollection, bson.M{"_id": bson.ObjectIdHex(usr.UserID)}, bson.M{"$set": bson.M{"pswd": usr.Password}})
	} else if roleID == apicontext.AdminRole {
		err = mongodb.Update(mongodb.AdminCollection, bson.M{"_id": bson.ObjectIdHex(usr.UserID)}, bson.M{"$set": bson.M{"pswd": usr.Password}})
	} else {
		err = errors.New("user role not authentic")
	}

	return err
}

// Signout logouts user
func (usr UserInfo) Signout(roleID, userID string) (err error) {

	if roleID == apicontext.CustomerRole {
		err = mongodb.Exists(mongodb.CustomerCollection, bson.M{"_id": bson.ObjectIdHex(usr.UserID)}, nil)
	} else if roleID == apicontext.AdminRole {
		err = mongodb.Exists(mongodb.AdminCollection, bson.M{"_id": bson.ObjectIdHex(usr.UserID)}, nil)
	} else {
		err = errors.New("user role not authentic")
	}

	return
}
