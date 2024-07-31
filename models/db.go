package models

import (
	"log"
	"time"

	"test/helper"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var db = ""
var globalS *mgo.Session

func init() {

	config_vars := helper.GetConfigVars()
	db = config_vars["BBDD_DB"]
	dialInfo := &mgo.DialInfo{
		Addrs:   []string{config_vars["BBDD_HOST"]},
		Source:  config_vars["BBDD_SOURCE"],
		Timeout: time.Second * 10,
	}

	// if config_vars["BBDD_USER"] != "" {
	// 	log.Println("ENTRAMOS")
	// 	dialInfo = &mgo.DialInfo{
	// 		Addrs:    []string{config_vars["BBDD_HOST"]},
	// 		Source:   config_vars["BBDD_SOURCE"],
	// 		Timeout:  time.Second * 10,
	// 		Username: config_vars["BBDD_USER"],
	// 		Password: config_vars["BBDD_PASS"],
	// 	}
	// }
	log.Println("passed")
	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalln("create session error ", err)
	}
	globalS = s
}

func connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	s := globalS.Copy()
	c := s.DB(db).C(collection)
	return s, c
}

func Insert(db, collection string, docs ...interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Insert(docs...)
}

func IsExist(db, collection string, query interface{}) bool {
	ms, c := connect(db, collection)
	defer ms.Close()
	count, _ := c.Find(query).Count()
	return count > 0
}

func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}
func FindOneOrder(db, collection string, query, selector, result interface{}, ordena string) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Sort(ordena).One(result)
}

func FindOneWRespuestas(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	pipeline := []bson.M{
		bson.M{
			"$lookup": bson.M{
				"from":         "respuestas",
				"localField":   "_id",
				"foreignField": "id_pregunta",
				"as":           "respuestas",
			},
		},
		bson.M{"$match": query},
	}
	return c.Pipe(pipeline).One(result)
}
func FindAllWRespuestas(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	pipeline := []bson.M{
		bson.M{
			"$lookup": bson.M{
				"from":         "respuestas",
				"localField":   "_id",
				"foreignField": "id_pregunta",
				"as":           "respuestas",
			},
		},
		bson.M{"$match": query},
	}
	return c.Pipe(pipeline).All(result)
}
func FindAllWPreguntas(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	pipeline := []bson.M{
		bson.M{
			"$lookup": bson.M{
				"from":         "preguntas",
				"localField":   "_id",
				"foreignField": "id_pregunta",
				"as":           "preguntas",
			},
		},
		bson.M{"$match": query},
	}
	return c.Pipe(pipeline).All(result)
}

func FindAll(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).All(result)
}
func FindAllOrder(db, collection string, query, selector, result interface{}, ordena string) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Sort(ordena).All(result)
}

func FindLimit(db, collection string, query, selector, result interface{}, limit int) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Limit(limit).All(result)
}

func Update(db, collection string, query, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Update(query, update)
}

func Remove(db, collection string, query interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Remove(query)
}

func RemoveAll(db, collection string, query interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	_, err := c.RemoveAll(query)
	return err
}
