package mongodb

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

var (
	session *mgo.Session
	err     error
)

func init() {
	session, err = connect()
	if err != nil {
		log.Fatalln(err)
	}
}

func connect() (*mgo.Session, error) {

	address := "127.0.0.1:27017"
	//username := ""
	//password := ""

	dial := &mgo.DialInfo{
		Addrs: []string{address},
		//Username: username,
		//Password: password,
		Timeout:  15 * time.Second,
		Database: DBName,
		//Source:   "admin",
	}

	session, err := mgo.DialWithInfo(dial)
	if err != nil {
		return session, errors.Wrap(err, "error dialing mongo")
	}

	return session, nil
}

//GetSession returns a copy of the current running mongo session
func GetSession() (*mgo.Session, error) {

	if session == nil {
		session, err = connect()
		if err != nil {
			return session, err
		}
	}

	return session.Copy(), nil
}
