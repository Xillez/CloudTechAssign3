package types

import "gopkg.in/mgo.v2/bson"

// WebhookInfo - Physical format of stored webhook info
type WebhookInfo struct {
	ID       bson.ObjectId `                       bson:"_id,omitempty"`
	URL      string        `json:"webhookURL"      bson:"webhookURL"`
	CurrID   bson.ObjectId `                       bson:"currId"`
	MinValue float64       `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxValue float64       `json:"maxTriggerValue" bson:"maxTriggerValue"`
}

// CurrencyInfo - Physical format of stored currency info
type CurrencyInfo struct {
	ID             bson.ObjectId `                       bson:"_id,omitempty"`
	BaseCurrency   string        `json:"baseCurrency"    bson:"baseCurrency"`
	TargetCurrency string        `json:"targetCurrency"  bson:"targetCurrency"`
	Rate           float64       `json:"rate"            bson:"rate"`
}

// WebhookReg - Contains all nessecary info to invoke a webhook
type WebhookReg struct {
	URL            string  `json:"webhookURL"`
	BaseCurrency   string  `json:"baseCurrency"`
	TargetCurrency string  `json:"targetCurrency"`
	MinValue       float64 `json:"minTriggerValue"`
	MaxValue       float64 `json:"maxTriggerValue"`
}

// WebhookInv - Payload for invoking the webhook
type WebhookInv struct {
	BaseCurrency   string  `json:"baseCurrency"`
	TargetCurrency string  `json:"targetCurrency"`
	CurrentRate    float64 `json:"currentRate"`
	MinValue       float64 `json:"minTriggerValue"`
	MaxValue       float64 `json:"maxTriggerValue"`
}

// CurrencyReq - Used variously, requesting exchange, running average, etc.
type CurrencyReq struct {
	BaseCurrency   string `json:"baseCurrency"`
	TargetCurrency string `json:"targetCurrency"`
}

// LatestCurrResp - Used to store the internal data of response from api.fixer.io/latest
type LatestCurrResp struct {
	BaseCurrency string             `json:"base"`
	Date         string             `json:"date"` // Date format yyyy-mm-dd
	Rates        map[string]float64 `json:"rates"`
}
