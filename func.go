package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

///// GLOBAL FLAGS & VARIABLES
var bcurl = "http://172.21.177.60:6500"

//var bcurl = "http://localhost:6500"

type RawMaterialTransaction struct {
	SerialNo       int    `json:"SerialNo"`
	ProductCode    string `json:"ProductCode"`
	ProductName    string `json:"ProductName"`
	ProductBatchNo string `json:"ProductBatchNo"`
	Quantity       int    `json:"Quantity"`
	//TxnTimestamp string `json:"TxnTimestamp"`
	RawMaterial []RawMaterial
}

type RawMaterial struct { //the raw material document
	RawMaterialBatchNo         string  `json:"RawMaterialBatchNo"`
	RawMaterialsID             string  `json:"RawMaterialsID"`
	RawMaterialName            string  `json:"RawMaterialName"`
	RawMaterialQuantity        float32 `json:"RawMaterialQuantity"`
	RawMaterialMeasurementUnit string  `json:"RawMaterialMeasurementUnit"`
}

var ProductData []Product // `Product` array, to be saved as gob file

/////

type DeliveryTransaction struct {
	SerialNo     int    `json:"SerialNo"`
	RecGenerator string `json:"RecGenerator"`
	ShipmentID   string `json:"ShipmentID"`
	Timestamp    string `json:"TxnTimestamp"`
	//Longitude         string `json:"Longitude"`
	//Latitude          string `json:"Latitude"`
	ShippedFromCompID string `json:"ShippedFromCompID"`
	ShippedToCompID   string `json:"ShippedToCompID"`
	//LocationID        string `json:"LocationID"`
	//DeliveryStatus    string `json:"DeliveryStatus"`
	DeliveryType string `json:"DeliveryType"`
	Product      []Product
	Document     Document
	UserSign     Signature
}

type Product struct {
	ProductCode  string `json:"ProductCode"`
	ProductName  string `json:"ProductName"`
	ProductBatch ProductBatch
}
type ProductBatch struct {
	ProductBatchNo       string `json:"ProductBatchNo"`
	ProductBatchQuantity int    `json:"ProductBatchQuantity"`
}
type Document struct {
	DocumentURL  string
	DocumentType string
	DocumentHash string
	DocumentSign string
}

type Signature struct {
	User     string
	Verify   bool
	USerSign string
}

func CheckSerialNo(SerialNo int) bool { // check the SerialNo with blockchain nodes
	result := false
	receivedSerialNo := strconv.Itoa(SerialNo)
	url := bcurl + "/query/raw/SerialNo/" + receivedSerialNo
	log.Println(url)

	resp, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}
	res, _ := http.DefaultClient.Do(resp)
	defer res.Body.Close()
	//log.Println(resp)

	var body interface{}
	err1 := json.NewDecoder(res.Body).Decode(&body)
	if err1 != nil {
		log.Println(err1)
	}
	//log.Println(body)

	if body == nil { //if the response is null then return ture. It means the SerialNo can be used.
		result = true
	}
	log.Println(result)

	return result
}

func CheckProductBatchNo(ProductBatchNo string) (bool, []RawMaterialTransaction) { // check the ProductBatchNo with blockchain nodes
	//log.Println("called")
	result := false
	receivedProductBatchNo := ProductBatchNo
	url := bcurl + "/query/raw/ProductBatchNo/" + receivedProductBatchNo
	log.Println(url)

	resp, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}
	res, _ := http.DefaultClient.Do(resp)
	defer res.Body.Close()
	//log.Println(resp)

	var body []RawMaterialTransaction
	err1 := json.NewDecoder(res.Body).Decode(&body)
	if err1 != nil {
		log.Println(err1)
	}

	if body == nil { //if the response is null then return ture. It means the ProductCode can be used.
		result = true
	} else {
		result = false
	}
	//log.Println(result)

	return result, body
}

func CheckProductBatchNonCode(ProductCode string, ProductBatchNo string) (bool, []RawMaterialTransaction) { // check the ProductBatchNo with blockchain nodes
	result := false //check if information is correct

	receivedProductBatchNo := ProductBatchNo
	var querybyProductCodeAndBatchNo []RawMaterialTransaction

	//var querybyProductCode []DeliveryTransaction
	status, querybyProductBatchNo := CheckProductBatchNo(receivedProductBatchNo)
	if status == true {
		return result, querybyProductCodeAndBatchNo
	}
	//log.Println(querybyProductCode)

	querybyquerybyProductBatchNoLen := len(querybyProductBatchNo)

	for i := 0; i < querybyquerybyProductBatchNoLen; i++ {

		if querybyProductBatchNo[i].ProductBatchNo == receivedProductBatchNo {
			result = true
			querybyProductCodeAndBatchNo = append(querybyProductCodeAndBatchNo, querybyProductBatchNo[i])
		}

	}

	return result, querybyProductCodeAndBatchNo
}

func CheckProductQuantity(receivedTx []RawMaterialTransaction, ProductBatchQuantity int) bool { // check the ProductCode with blockchain nodes
	var querybyProductCodeAndBatchNo []RawMaterialTransaction
	result := false
	receivedProductBatchQuantity := ProductBatchQuantity
	querybyProductCodeAndBatchNo = receivedTx

	Txlength := len(receivedTx)

	if receivedProductBatchQuantity <= querybyProductCodeAndBatchNo[Txlength-1].Quantity {
		result = true
	}

	return result

}

func GenerateNewRawMaterialTx(newTx RawMaterialTransaction) { //Post this new transaction
	posturl := bcurl + "/raw"
	//log.Println(newTx)
	var receivedNewTx RawMaterialTransaction
	receivedNewTx = newTx

	formdata, err := json.Marshal(receivedNewTx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(string(formdata))

	body := bytes.NewBuffer(formdata)
	rsp, err := http.Post(posturl, "application/json", body)
	// rsp, err := http.NewRequest("POST", posturl, body)
	// rsp.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	body_byte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body_byte))

	log.Println("New TRANSACTION posted!")
}

func GenerateNewDeliveryTx(newDeliveryInfo DeliveryTransaction) { //Post this new transaction
	posturl := bcurl + "/ship"
	//posturl := "http://localhost:8880/raw"
	//log.Println(newDeliveryTx)
	var NewDeliveryTx DeliveryTransaction
	NewDeliveryTx = newDeliveryInfo

	log.Println("Prepare for posting")
	formdata, err := json.Marshal(NewDeliveryTx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(string(formdata))

	body := bytes.NewBuffer(formdata)
	rsp, err := http.Post(posturl, "application/json", body)
	// rsp, err := http.NewRequest("POST", posturl, body)
	// rsp.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	body_byte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body_byte))

	log.Println("New product posted!")
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
