package storage

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"log"
	"gopkg.in/mgo.v2/bson"
)

const (
	COLLECTION = "gifs"
)


type Element struct {
	Id    string `bson:"id" json:"id"`
	Url   string `bson:"url" json:"url"`
	Title string `bson:"title" json:"title"`
}

type Results struct {
	Data []Element `json:"data"`
}

//TODO add session for DBase interface
//type Session *mgo.Session

// ----------------------------------------

func NewStorage(db *mgo.Database) Storage {
	return Storage{db}
}

type Storage struct {
	DB *mgo.Database
}

func (h *Storage) InsertInitialData() {
	gifs := loadData()
	for _, element := range gifs.Data {
		err := h.DB.C(COLLECTION).Insert(element)
		check(err)
	}
}

func (h *Storage) SaveData() {
	var elements []Element
	err := h.DB.C(COLLECTION).Find(bson.M{}).All(&elements)
	check(err)
	res := Results{elements}

	f, err := os.Create("data_gifs")
	check(err)
	fmt.Println("Store data.")
	b, err := json.Marshal(res)
	f.Write(b)

	defer f.Close()
}

func loadData() Results {
	var s = new(Results)
	body, err := ioutil.ReadFile("data_gifs")
	if err != nil {
		os.Exit(1)
	}
	json.Unmarshal(body, &s)
	return *s
}

// ----------------------------------------------------------------

type DBase interface {
	Connect() (*mgo.Session, error)
	Insert(Element) error
	Remove(Element) error
	Get(string) (Element, error)
	GetAll() ([]Element, error)
	RemoveAll() (*mgo.ChangeInfo, error)
}

type DataBase struct {
	ConnectionPoint string
	Session *mgo.Session
}

func (db *DataBase) Connect() (*mgo.Session, error) {
	dbSession, err := mgo.Dial(db.ConnectionPoint)
	dbSession.SetMode(mgo.Monotonic, true)
	// Error check on every access
	dbSession.SetSafe(&mgo.Safe{})

	if err == nil {
		db.Session = dbSession
	}

	return dbSession, err
}

func (db *DataBase) GetAll() ([]Element, error){
	var elements []Element
	err := db.Session.DB("testDB").C(COLLECTION).Find(bson.M{}).All(&elements)
	return elements, err
}

func (db *DataBase) Get(id string) (Element, error){
	var element Element
	err := db.Session.DB("testDB").C(COLLECTION).Find(bson.M{"id": id}).One(&element)
	return element, err
}

func (db *DataBase) Insert(element Element) error{
	err := db.Session.DB("testDB").C(COLLECTION).Insert(element)
	return err
}

func (db *DataBase) Remove(element Element) error{
	err := db.Session.DB("testDB").C(COLLECTION).Remove(element)
	return err
}

func (db *DataBase) RemoveAll() (*mgo.ChangeInfo, error) {
	info, err := db.Session.DB("testDB").C(COLLECTION).RemoveAll(bson.M{})
	return info, err
}

func check(e error) {
	if e != nil {
		//TODO: Add better logging
		log.Fatal(e)
		os.Exit(1)
	}
}

