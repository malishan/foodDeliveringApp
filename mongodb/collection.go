package mongodb

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	DBName               = "foodApp"
	RestaurantCollection = "restaurantDetails"
	CustomerCollection   = "customerInfo"
	AdminCollection      = "adminInfo"
)

// Exists check if data is present in mongo database
func Exists(collectionName string, selector, filter bson.M) error {
	mongo, err := GetSession()
	if err != nil {
		return err
	}
	defer mongo.Close()

	var data interface{}
	err = mongo.DB("").C(collectionName).Find(selector).Select(filter).One(data)
	if err != nil {
		return err
	}

	return nil
}

// Insert - inserts data into mongo database
func Insert(collectionName string, d ...interface{}) error {
	mongo, err := GetSession()
	if err != nil {
		return err
	}
	defer mongo.Close()
	err = mongo.DB("").C(collectionName).Insert(d...)
	if err != nil {
		return err
	}
	return nil
}

// ReadOne - reads single document from mongo database
func ReadOne(collectionName string, selector, filter bson.M, data interface{}) error {
	mongo, err := GetSession()
	if err != nil {
		return err
	}
	defer mongo.Close()
	err = mongo.DB("").C(collectionName).Find(selector).Select(filter).One(data)
	if err != nil {
		return err
	}
	return nil
}

// ReadAll - reads multiple documents from mongo database
func ReadAll(collectionName string, selector, filter bson.M, data interface{}) error {
	mongo, err := GetSession()
	if err != nil {
		return err
	}
	defer mongo.Close()
	err = mongo.DB("").C(collectionName).Find(selector).Select(filter).All(data)
	if err != nil {
		return err
	}
	return nil
}

// Update - updates data into mongo database
func Update(collectionName string, selector, d interface{}) error {
	mongo, err := GetSession()
	if err != nil {
		return err
	}
	defer mongo.Close()
	err = mongo.DB("").C(collectionName).Update(selector, d)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAndReturn - updates data into mongo database and returns the updated document
func UpdateAndReturn(collectionName string, selector, update interface{}, returnObject interface{}) error {
	mongo, err := GetSession()
	if err != nil {
		return err
	}
	defer mongo.Close()

	change := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}
	_, err = mongo.DB("").C(collectionName).Find(selector).Apply(change, returnObject)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAll - updates multiple documents into mongo database
func UpdateAll(collectionName string, selector, d interface{}) (*mgo.ChangeInfo, error) {
	mongo, err := GetSession()
	if err != nil {
		return nil, err
	}
	defer mongo.Close()
	changeInfo, err := mongo.DB("").C(collectionName).UpdateAll(selector, d)
	if err != nil {
		return nil, err
	}
	return changeInfo, nil
}

// Bulk update - updates or creates data into mongo database
func BulkUpdate(collectionName string, pairs ...interface{}) (*mgo.BulkResult, error) {
	mongo, err := GetSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}
	defer mongo.Close()
	bulk := mongo.DB("").C(collectionName).Bulk()
	bulk.Update(pairs...)
	result, err := bulk.Run()
	return result, errors.Wrap(err, "error while updating")
}

// Upsert - updates or creates data into mongo database
func Upsert(collectionName string, selector, d interface{}) (*mgo.ChangeInfo, error) {
	mongo, err := GetSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}
	defer mongo.Close()
	changeInfo, err := mongo.DB("").C(collectionName).Upsert(selector, d)
	if err != nil {
		return nil, errors.Wrap(err, "error while upserting")
	}
	return changeInfo, nil
}

// Bulk Upsert - updates or creates data into mongo database
func BulkUpsert(collectionName string, pairs ...interface{}) (*mgo.BulkResult, error) {
	mongo, err := GetSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}
	defer mongo.Close()
	bulk := mongo.DB("").C(collectionName).Bulk()
	bulk.Upsert(pairs...)
	result, err := bulk.Run()
	return result, errors.Wrap(err, "error while upserting")
}

func BulkInsert(collectionName string, pairs ...interface{}) (*mgo.BulkResult, error) {
	mongo, err := GetSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}
	defer mongo.Close()
	bulk := mongo.DB("").C(collectionName).Bulk()
	bulk.Insert(pairs...)
	result, err := bulk.Run()
	return result, errors.Wrap(err, "error while Inserting")
}

// Delete - updates data into mongo database
func Delete(collectionName string, selector bson.M) error {
	mongo, err := GetSession()
	if err != nil {
		return err
	}
	defer mongo.Close()
	err = mongo.DB("").C(collectionName).Update(selector, bson.M{"$set": bson.M{"isActive": false}})
	if err != nil {
		return err
	}
	return nil
}

// DeleteAll - updates multiple documents setting their isActive=false in mongo database
func DeleteAll(collectionName string, selector bson.M) error {
	mongo, err := GetSession()
	if err != nil {
		return err
	}
	defer mongo.Close()
	_, err = mongo.DB("").C(collectionName).UpdateAll(selector, bson.M{"$set": bson.M{"isActive": false}})
	if err != nil {
		return err
	}
	return nil
}

// Hard delete - updates data into mongo database
func HardDelete(collectionName string, selector bson.M) error {
	mongo, err := GetSession()
	if err != nil {
		return err
	}
	defer mongo.Close()
	err = mongo.DB("").C(collectionName).Remove(selector)
	if err != nil {
		return err
	}
	return nil
}

// Hard delete all - updates data into mongo database
func HardDeleteAll(collectionName string, selector bson.M) (*mgo.ChangeInfo, error) {
	mongo, err := GetSession()
	if err != nil {
		return nil, err
	}
	defer mongo.Close()
	changeInfo, err := mongo.DB("").C(collectionName).RemoveAll(selector)
	if err != nil {
		return nil, err
	}
	return changeInfo, err
}

// AggregateOne executes aggregation query on db and fetches single document
func AggregateOne(collectionsName string, query []bson.M, data interface{}) error {

	session, err := GetSession()
	if err != nil {
		errSessionFailed := fmt.Errorf("error encountered while creating session : %v", err)
		return errSessionFailed
	}
	defer session.Close()

	mgo.SetDebug(true)

	err = session.DB("").C(collectionsName).Pipe(query).One(data)
	if err != nil {
		errAggregationFailed := fmt.Errorf("error encountered while executing aggregation query : %v", err)
		return errAggregationFailed
	}
	return err
}

// AggregateAll executes aggregation query on db and fetches multiple document
func AggregateMany(collectionsName string, query []bson.M, data interface{}) error {

	session, err := GetSession()
	if err != nil {
		errSessionFailed := fmt.Errorf("error encountered while creating session : %v", err)
		return errSessionFailed
	}
	defer session.Close()

	mgo.SetDebug(true)

	err = session.DB("").C(collectionsName).Pipe(query).All(data)
	if err != nil {
		errAggregationFailed := fmt.Errorf("error encountered while executing aggregation query : %v", err)
		return errAggregationFailed
	}
	return err
}

// ReadDistinct - reads distinct values of the provided field from mongo database
func ReadDistinct(collectionName, field string, selector, data interface{}) error {
	mongo, err := GetSession()
	if err != nil {
		return err
	}
	defer mongo.Close()
	err = mongo.DB("").C(collectionName).Find(selector).Distinct(field, data)
	if err != nil {
		return err
	}
	return nil
}

// Count returns number of documents those satisfied findQuery
func Count(collName string, findQ bson.M) (int, error) {
	s, err := GetSession()
	if err != nil {
		err = fmt.Errorf("failed to get mongo session error: %v", err)
		return 0, err
	}
	defer s.Close()
	n, err := s.DB("").C(collName).Find(findQ).Count()
	if err != nil {
		err = fmt.Errorf("mongo count failed for: %s, findQ: %+v, error: %v", collName, findQ, err)
		return 0, err
	}
	return n, nil
}