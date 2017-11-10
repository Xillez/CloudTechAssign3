package mongodb

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Xillez/CloudTechAssign2/types"

	"github.com/Xillez/CloudTechAssign2/utils"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB stores the details of the DB connection.
type MongoDB struct {
	DatabaseURL  string
	DatabaseName string
	WebCollName  string
	CurrCollName string
}

var Warn = "[WARNING]: "
var Error = "[ERROR]: "
var Info = "[INFO]: "

// Init - initializes the mongoDB database
func (db *MongoDB) Init() error {
	// Dial database
	log.Println(Info + "Dialing database")
	session, err := mgo.Dial(db.DatabaseURL)

	// Check for session error
	if err != nil {
		log.Println(Error + "Failed dialing database")
		return err
	}

	// Set up currency collation indexing
	index := mgo.Index{
		Key:        []string{"date"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err != nil {
		log.Println(Error + "Failed to drop currency databse!")
		return err
	}

	// Ensure Currency collation follows the index
	log.Println(Info + "Make collection: " + db.CurrCollName + "Ensure \"index\"")
	err = session.DB(db.DatabaseName).C(db.CurrCollName).EnsureIndex(index)
	if err != nil {
		log.Println(Error + "Failed to ensure \"index\" on collection: " + db.CurrCollName + "!")
		return err
	}

	// Postpone closing connection until we return
	defer session.Close()

	log.Println(Info + "Initializing successfull!")
	// Nothing bad happened!
	return nil
}

// GetWebhook - Gets the webhook with the given id
// if no entry is found with the given id, it returns an empty webhook struct
func (db *MongoDB) GetWebhook(id string, webhook *types.WebhookInfo) utils.CustError {
	// Dial database
	log.Println(Info + "Dialing database")
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		log.Println(Error + "Failed dialing the database!")
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Try to get webhook details from database
	log.Println(Info + "Trying to find the webhook with the given id")
	errFind := session.DB(db.DatabaseName).C(db.WebCollName).Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&webhook)
	if errFind != nil {
		// Something went wrong! Inform the user!
		log.Println(Error + "Failed to find a webhook! Informing user!")
		return utils.CustError{http.StatusInternalServerError, "Failed to find webhook with id: \"" + id + "\" in the database!"}
	}

	log.Println(Info + "GetWebhook finnished successfull!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

// GetAllWebhooks - Gets all the webhooks
// if no entries are found, an error is returned
func (db *MongoDB) GetAllWebhooks(webhook *[]types.WebhookInfo) utils.CustError {
	// Dial database
	log.Println(Info + "Dialing databse")
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Try to get webhook details from database
	log.Println(Info + "Trying to get all wbhooks in the database")
	errFind := session.DB(db.DatabaseName).C(db.WebCollName).Find(nil).All(webhook)
	if errFind != nil {
		// Something went wrong! Inform the user!
		log.Println(Error + "Couldn't get all webhooks from database! Inform the user!")
		return utils.CustError{http.StatusInternalServerError, ""}
	}

	log.Println(Info + "GetAllWebhooks finished successfull!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

// GetAllCurr - Gets all currencies
// if no entries are found, an error is returned
func (db *MongoDB) GetAllCurr(curr *[]types.CurrencyInfo) utils.CustError {
	// Dial database
	log.Println(Info + "Dialing database!")
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		log.Println(Error + "Failed to dail database!")
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Try to get currency details from database
	log.Println(Info + "Trying to get all currencies from database")
	errFind := session.DB(db.DatabaseName).C(db.CurrCollName).Find(nil).All(curr)
	if errFind != nil {
		// Something went wrong! Inform the user!
		log.Println(Error + "Failed getting all currencies from database! Inform the user!")
		return utils.CustError{http.StatusInternalServerError, "Failed to get all currencies from the database!"}
	}

	log.Println(Info + "GetAllCurr finished successfull!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

// GetCurrByDate - Gets the currency with the given date
// if no entry is found, it returns an error and ignores updating the given interface
func (db *MongoDB) GetCurrByDate(date string, curr *types.CurrencyInfo) utils.CustError {
	// Dial database
	log.Println(Info + "Dialing database!")
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		log.Println(Error + "Failed to dail database!")
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Try to get currency details from database
	log.Println(Info + "Trying to get currency with given date")
	errFind := session.DB(db.DatabaseName).C(db.CurrCollName).Find(bson.M{"date": date}).One(&curr)
	if errFind != nil {
		// Something went wrong! Inform the user!
		log.Println(Error + "Failed getting currency with given date!")
		return utils.CustError{http.StatusBadRequest, "Failed to find currency with date: \"" + date + "\" in the database!"}
	}

	log.Println(Info + "GetCurrByDate finished successfull!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

// AddWebhook - Adds "webhook" to Webhook collection
func (db *MongoDB) AddWebhook(webhook types.WebhookInfo) utils.CustError {
	// Dial database
	log.Println(Info + "Dialing database!")
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		log.Println(Error + "Failed to dail database!")
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Trying to insert a webhook
	log.Println(Info + "Trying to insert a given webhook")
	errInsert := session.DB(db.DatabaseName).C(db.WebCollName).Insert(webhook)
	if errInsert != nil {
		// Something went wrong! Inform the user!
		log.Println(Error + "Failed inserting given wehook! Inform the user!")
		return utils.CustError{http.StatusInternalServerError, "Failed inserting webhook to database!"}
	}

	log.Println(Info + "AddWebhook finished successfull!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

// AddCurr - Adds "curr" to Currency collection
func (db *MongoDB) AddCurr(curr types.CurrencyInfo) utils.CustError {
	// Dial database
	log.Println(Info + "Dialing database!")
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		log.Println(Error + "Failed to dail database!")
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Trying to insert a currency
	log.Println(Info + "Trying to insert a given currency")
	errInsert := session.DB(db.DatabaseName).C(db.CurrCollName).Insert(curr)
	if errInsert != nil {
		// Something went wrong! Inform the user!
		log.Println(Error + "Failed inserting given currency! Inform the user!")
		return utils.CustError{http.StatusInternalServerError, "Failed inserting currency to database!"}
	}

	log.Println(Info + "AddCurr finished successfull!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

// UpdateCurr - Updates "date" currency with "curr"
func (db *MongoDB) UpdateCurr() utils.CustError {
	// Dial database
	log.Println(Info + "Dialing database!")
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		log.Println(Error + "Failed to dail database!")
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	currResp := types.CurrencyInfo{}

	// Fetch latest currencies and decode them
	log.Println(Info + "Fetching and Decoding latest currency request from fixer")
	errCurr := utils.FetchDecodedJSON(utils.FixerURL+"/latest?base=EUR", &currResp)
	if errCurr.Status != 0 {
		log.Println(Error + "Failed fetching latest currencies! | Status: " + strconv.Itoa(errCurr.Status) + "Msg: " + errCurr.Msg) //errCurr.Error())
		return utils.CustError{http.StatusInternalServerError, ""}
	}

	currResp.Date = time.Now().Format("2006-01-02")

	// Trying to update a currency with given date
	log.Println(Info + "Trying to insert new currency")
	errInsert := session.DB(db.DatabaseName).C(db.CurrCollName).Insert(&currResp)
	if errInsert != nil {
		// User doesn't need to know that the program couldn't update currencies
		// correctly, he only needs to know that something went wrong on the server
		// not a user induced funciton
		log.Println(Error + "Failed to insert the new currency!")
		// Something went wrong! Inform the user!
		return utils.CustError{http.StatusInternalServerError, ""}
	}

	log.Println(Info + "UpdateCurr finished successfull!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

// DelWebhook - Deletes the webhook with the given id
// if no entry is found with the given id, it returns an empty currency struct
func (db *MongoDB) DelWebhook(id string) utils.CustError {
	// Dial database
	log.Println(Info + "Dialing database!")
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		log.Println(Error + "Failed to dail database!")
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Try to get currency details from database
	log.Println(Info + "Trying to delete a webhook with given id")
	errDel := session.DB(db.DatabaseName).C(db.WebCollName).Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if errDel != nil {
		// Something went wrong! Inform the user!
		log.Println(Warn + "Failed to delete webhook from the database! Inform the user!")
		return utils.CustError{http.StatusBadRequest, "Failed to delete webhook from the database!"}
	}

	log.Println(Info + "DelWebhook finished successfull!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

func (db *MongoDB) Count(collName string) (int, utils.CustError) {
	// Dial database
	log.Println(Info + "Dialing database!")
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		log.Println(Error + "Failed to dail database!")
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Try to get nr of elements from database
	log.Println(Info + "Trying to count nr elements in collection: " + collName)
	count, err := session.DB(db.DatabaseName).C(collName).Count()
	if err != nil {
		// Something went wrong! Inform the user!
		log.Println(Warn + "Failed to count nr elements in collection: " + collName + "! Inform the user!")
		return 0, utils.CustError{http.StatusInternalServerError, ""}
	}

	log.Println(Info + "Count finished successfull!")
	// Nothing bad happened
	return count, utils.CustError{0, utils.ErrorStr[0]}
}

// InvokeWebhooks - Invokes webhooks,
// if true, it'll check against min/max values
// if false these check are ignored
func (db *MongoDB) InvokeWebhooks(checkMinMax bool) utils.CustError {
	webhooks := []types.WebhookInfo{}
	curr := types.CurrencyInfo{}
	client := &http.Client{}

	// Get all currency info
	log.Println(Info + "Trying to get currency by date for invoking of webhooks")
	err := db.GetCurrByDate(time.Now().Format("2006-01-02"), &curr)
	if err.Status != 0 {
		log.Println(Error + "Failed to get currency by date for invoking of webhooks!")
	}

	// Get all webhooks
	log.Println(Info + "Trying to get all webhooks from database for invoking")
	err = db.GetAllWebhooks(&webhooks)
	if err.Status != 0 {
		log.Println(Error + "Failed to get all webhooks from database for invoking!")
	}

	// Find how may webhooks we got
	count := len(webhooks)

	// Loop through them
	log.Println(Info + "Trying to invoke webhooks!")
	for i := 0; i < count; i++ {
		log.Println("--------------- Log for " + strconv.Itoa(i) + "---------------")
		// Store target currency rate for convinience
		targetCurrRate := curr.Rates[webhooks[i].TargetCurrency]
		// Check if we should invoke all hooks, or by min/max condition only
		if !checkMinMax {
			log.Println(Warn + "Going to ignore minValue/maxValue checking, being evaluated!")
		}
		if !checkMinMax || // min && max => false, we can still invoke
			targetCurrRate < webhooks[i].MinValue ||
			targetCurrRate > webhooks[i].MaxValue {
			// Stringify WebhookInv struct
			reqBody := "{\"baseCurrency\":\"" + webhooks[i].BaseCurrency +
				"\", \"targetCurrency\":\"" + webhooks[i].TargetCurrency +
				"\", \"currentRate\":\"" + strconv.FormatFloat(targetCurrRate, 'f', -1, 64) +
				"\", \"minTriggerValue\":" + strconv.FormatFloat(webhooks[i].MinValue, 'f', -1, 64) +
				", \"maxTriggerValue\":" + strconv.FormatFloat(webhooks[i].MaxValue, 'f', -1, 64) + "}"

			// Build request
			log.Println(Info + "Trying to build a new POST request with webhooks given URL")
			req, errNewReq := http.NewRequest("POST", webhooks[0].URL, bytes.NewReader([]byte(reqBody)))
			if errNewReq != nil {
				log.Println(Error + "Failed to build a new POST request to webhook!")
			}
			req.Header.Set("Content-Type", "application/json")

			// Send request
			log.Println(Info + "Trying to send POST request with requested data")
			_, errClientDo := client.Do(req)
			if errClientDo != nil {
				// LOG: RUNNING THE REQUEST FAILED!
				log.Println(Error + "Failed to send newly created POST request to webhook!")
			}
		}
	}

	log.Println(Info + "InvokeWebhooks finished successfull!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}
