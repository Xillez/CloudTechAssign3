// Program is made in collaboration with / got help from
//          - Jonas J. Solsvik
// 			- Zohaib Butt
// 			- Eldar Hauge Torkelsen

package main

import (
	"net/http"
	"os"
)

var baseURL = "https://api.github.com/repos/"

// RunningTest tells the program wether we're testing or not
var RunningTest = false

var errorStr = []string{
	"NoError",
	"InvalidServer",
	"StatusCodeOutOfBounds",
	"InvalidVitalData",
	"ReadFileFaliure",
	"UnmarshalFail",
	"InvalidURL",
	"EncodeFail",
	"DecodeFail",
}

// CustError - A custom error type
type CustError struct {
	status int
	msg    string
}

// WebhookReg - Contains all nessecary info to invoke a webhook
type WebhookReg struct {
	URL            string  `json:"webhookURL" bson:"webhookURL"`
	BaseCurrency   string  `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency string  `json:"targetCurrency" bson:"targetCurrency"`
	MinValue       float32 `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxValue       float32 `json:"maxTriggerValue" bson:"maxTriggerValue"`
}

// WebhookInv - Payload for invoking the webhook
type WebhookInv struct {
	BaseCurrency   string  `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency string  `json:"targetCurrency" bson:"targetCurrency"`
	CurrentRate    float32 `json:"currentRate" bson:"currentRate"`
	MinValue       float32 `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxValue       float32 `json:"maxTriggerValue" bson:"maxTriggerValue"`
}

// CurrencyReq - Used variously, requesting exchange, running average, etc.
type CurrencyReq struct {
	BaseCurrency   string `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency string `json:"targetCurrency" bson:"targetCurrency"`
}

var db *MongoDB

// HandlerRoot - Root path handler, handles registration, viewing and
// deleting of webhooks
func handlerRoot(w http.ResponseWriter, r *http.Request) {
	switch RunningTest {
	case false:
		// Get general repo data
		/*if checkPrintErr(fetchDecodedJSON(baseURL+parts[4]+"/"+parts[5], &project), w) {
			return
		}*/
	case true:
		// Load general testrepo data
		/*if checkPrintErr(openMainJSON(&project), w) {
			return
		}*/
	}
}

// HandlerLatest - Latest Currency Exchange handler
func handlerLatest(w http.ResponseWriter, r *http.Request) {
	switch RunningTest {
	case false:
		// Get general repo data
		/*if checkPrintErr(fetchDecodedJSON(baseURL+parts[4]+"/"+parts[5], &project), w) {
			return
		}*/
	case true:
		// Load general testrepo data
		/*if checkPrintErr(openMainJSON(&project), w) {
			return
		}*/
	}
}

// HandlerAverage - 7-day Average Currency Exchange handler
func handlerAverage(w http.ResponseWriter, r *http.Request) {
	switch RunningTest {
	case false:
		// Get general repo data
		/*if checkPrintErr(fetchDecodedJSON(baseURL+parts[4]+"/"+parts[5], &project), w) {
			return
		}*/
	case true:
		// Load general testrepo data
		/*if checkPrintErr(openMainJSON(&project), w) {
			return
		}*/
	}
}

// HandlerEvalTrigger - Evaluation Trigger handler, used when service is evaluated
func handlerEvalTrigger(w http.ResponseWriter, r *http.Request) {
	switch RunningTest {
	case false:
		// Get general repo data
		/*if checkPrintErr(fetchDecodedJSON(baseURL+parts[4]+"/"+parts[5], &project), w) {
			return
		}*/
	case true:
		// Load general testrepo data
		/*if checkPrintErr(openMainJSON(&project), w) {
			return
		}*/
	}
}

func main() {
	db = &MongoDB{"mongodb://localhost", "Books", "book"}
	db.Init()
	http.HandleFunc("/root", handlerRoot)
	http.HandleFunc("/root/latest", handlerLatest)
	http.HandleFunc("/root/average", handlerAverage)
	http.HandleFunc("/root/evaluationtrigger", handlerEvalTrigger)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

/*
book := Book{ /*bson.NewObjectId(), * "The story of Me", "Me", time.Now()}
fmt.Println("Book created!")

/*dbIds := book.ID
fmt.Println("Book's id saved")

db.Add(book)
fmt.Println("Book added!")

fmt.Println("Current nr of books: " + strconv.Itoa(db.Count()))

book, ret := db.Get(string(dbIds))
if ret != true {
	fmt.Println("Retrieved book from database!")
} else {
	fmt.Println("Failed to retrieved book from database!")
}

fmt.Println("Displaying book...")
fmt.Println("Name: " + book.Text)
*/

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
