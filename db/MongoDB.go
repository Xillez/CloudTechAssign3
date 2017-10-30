package db

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	types "github.com/Xillez/CloudTechAssign2/types"
)

// Init - initializes the mongodb database
func (db *MongoDB) Init() error {
	// Dial database
	session, err := mgo.Dial(db.DatabaseURL)

	// Check for session error
	if err != nil {
		return err
	}

	// Set up currency collation indexing
	index := mgo.Index{
		Key:        []string{"targetCurrency"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	// Ensure Currency collation follows the index
	err = session.DB(db.DatabaseName).C(db.CurrCollName).EnsureIndex(index)
	if err != nil {
		return err
	}

	// Set up currency collection for future webhooks
	currUpdate()

	// Postpone closing connection until we return
	defer session.Close()

	// Nothing bad happened!
	return nil
}

// GetWebhook - Gets the webhook with the given id
// if no entry is found with the given id, it returns an empty webhook struct
func (db *MongoDB) GetWebhook(id string, webhook *types.WebhookInfo) CustError {
	// Dial database
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Try to get webhook details from database
	errFind := session.DB(db.DatabaseName).C(db.WebCollName).Find(bson.M{"_id": id}).One(&webhook)
	if errFind != nil {
		// Something went wrong! Inform the user!
		return CustError{http.StatusInternalServerError, "Failed to find webhook with id: \"" + id + "\" in the database | " + errFind.Error()}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

// GetCurrByTarget - Gets the currency with the given target currency
// if no entry is found with the given id, it returns an error and ignores updating the given interface
func (db *MongoDB) GetCurrByTarget(target string, curr *types.CurrencyInfo) CustError {
	// Dial database
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Try to get currency details from database
	errFind := session.DB(db.DatabaseName).C(db.CurrCollName).Find(bson.M{"targetCurrency": target}).One(&curr)
	if errFind != nil {
		// Something went wrong! Inform the user!
		return CustError{http.StatusInternalServerError, "Failed to find currency with target: \"" + target + "\" in the database | " + errFind.Error()}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

// GetCurrByID - Gets the currency with the given id
// if no entry is found with the given id, it returns an error and ignores updating the given interface
func (db *MongoDB) GetCurrByID(id string, curr *types.CurrencyInfo) CustError {
	// Dial database
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Try to get currency details from database
	errFind := session.DB(db.DatabaseName).C(db.CurrCollName).Find(bson.M{"_id": id}).One(&curr)
	if errFind != nil {
		// Something went wrong! Inform the user!
		return CustError{http.StatusInternalServerError, "Failed to find currency with id: \"" + id + "\" in the database | " + errFind.Error()}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

// AddWebhook - Adds "webhook" to Webhook collection
func (db *MongoDB) AddWebhook(webhook types.WebhookInfo) CustError {
	// Dial database
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Trying to insert a webhook,
	errInsert := session.DB(db.DatabaseName).C(db.WebCollName).Insert(webhook)
	if errInsert != nil {
		// Something went wrong! Inform the user!
		return CustError{http.StatusInternalServerError, "Failed Inserting to database | " + errInsert.Error()}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

// AddCurr - Adds "curr" to Currency collection
func (db *MongoDB) AddCurr(curr types.CurrencyInfo) CustError {
	// Dial database
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Trying to insert a webhook,
	errInsert := session.DB(db.DatabaseName).C(db.CurrCollName).Insert(curr)
	if errInsert != nil {
		// Something went wrong! Inform the user!
		return CustError{http.StatusInternalServerError, "Failed Inserting to database | " + errInsert.Error()}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

// UpdateCurr - Updates "target" currency with "curr"
func (db *MongoDB) UpdateCurr(target string, curr interface{}) CustError {
	// Dial database
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Trying to insert a webhook,
	errInsert := session.DB(db.DatabaseName).C(db.CurrCollName).Update(bson.M{"targetCurrency": target}, curr)
	if errInsert != nil {
		// Something went wrong! Inform the user!
		return CustError{http.StatusBadRequest, "Failed Inserting to database | " + errInsert.Error()}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

// DelWebhook - Gets the curreny with the given target currency
// if no entry is found with the given id, it returns an empty currency struct
func (db *MongoDB) DelWebhook(id string) CustError {
	// Dial database
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Try to get currency details from database
	errDel := session.DB(db.DatabaseName).C(db.CurrCollName).Remove(bson.M{"_id": id})
	if errDel != nil {
		// Something went wrong! Inform the user!
		return CustError{http.StatusBadRequest, "Failed to find webhook with id: \"" + id + "\" in the database | " + errDel.Error()}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

/*
Add adds new students to the storage.
*
func (db *MongoDB) Add(book Book) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(book)

	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
		return err
	}

	return nil
}

/*
Count returns the current count of the students in in-memory storage.
*
func (db *MongoDB) Count() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// handle to "db"
	count, err := session.DB(db.DatabaseName).C(db.CollectionName).Count()
	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}

	return count
}

/*
Get returns a student with a given ID or empty student struct.
*
func (db *MongoDB) Get(keyID string) (Book, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	book := Book{}
	allWasGood := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{"_id": keyID}).One(&book)
	if err != nil {
		allWasGood = false
	}

	return book, allWasGood
}*/
