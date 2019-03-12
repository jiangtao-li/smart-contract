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

type RawMaterialTransaction struct {
	SerialNo       int    `json:"SerialNo"`
	ProductCode    string `json:"ProductCode"`
	ProductName    string `json:"ProductName"`
	ProductBatchNo string `json:"ProductBatchNo"`
	Quantity       int    `json:"Quantity"`
	//TxnTimestamp string `json:"TxnTimestamp"`
	RawMaterialInfo RawMaterial
}

type RawMaterial struct { //the raw material document
	RawMaterialBatchNo         string  `json:"RawMaterialBatchNo"`
	RawMaterialsID             string  `json:"RawMaterialsID"`
	RawMaterialName            string  `json:"RawMaterialName"`
	RawMaterialQuantity        float32 `json:"RawMaterialQuantity"`
	RawMaterialMeasurementUnit string  `json:"RawMaterialMeasurementUnit"`
}

//var ProductData []Product // `Product` array, to be saved as gob file

/////

type DeliveryTransaction struct {
	SerialNo          int    `json:"SerialNo"`
	RecGenerator      string `json:"RecGenerator"`
	ShipmentID        string `json:"ShipmentID"`
	Timestamp         string `json:"TxnTimestamp"`
	Longitude         string `json:"Longitude"`
	Latitude          string `json:"Latitude"`
	ShippedFromCompID string `json:"ShippedFromCompID"`
	ShippedToCompID   string `json:"ShippedToCompID"`
	LocationID        string `json:"LocationID"`
	DeliveryStatus    string `json:"DeliveryStatus"`
	DeliveryType      string `json:"DeliveryType"`
	Product           Product
	Document          Document
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

type NewProductTransaction struct {
	SerialNo        int
	ProductID       string
	ProductName     string
	RawMaterialsID  string
	RawMaterialName string
}

//WarehousingTransaction type 4 Txn: recored the quantity of each 'Product ID+ batch ID + raw material ID'
type WarehousingTransaction struct {
	SerialNo        int
	ProductID       string
	ProductName     string
	ProductBatch    string
	RawMaterialsID  string
	RawMaterialName string
	MaterialBatch   string
	Quantity        int
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

func CheckProductCode(ProductID string) (bool, []NewProductTransaction) { // check the ProductCode with blockchain nodes
	//log.Println("called")
	result := false
	receivedProductID := ProductID
	url := bcurl + "/query/new/ProductID/" + receivedProductID
	log.Println(url)

	resp, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}
	res, _ := http.DefaultClient.Do(resp)
	defer res.Body.Close()
	//log.Println(resp)

	var body []NewProductTransaction
	err1 := json.NewDecoder(res.Body).Decode(&body)
	if err1 != nil {
		log.Println(err1)
	}

	if body == nil { //if the response is null then return ture. It means the ProductCode can be used.
		result = true
	}
	//log.Println(result)

	return result, body
}

func CheckProductIDforBatch(ProductID string) (bool, []WarehousingTransaction) { // check the ProductCode with blockchain nodes
	//log.Println("called")
	result := false
	receivedProductID := ProductID
	url := bcurl + "/query/war/ProductID/" + receivedProductID
	log.Println(url)

	resp, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}
	res, _ := http.DefaultClient.Do(resp)
	defer res.Body.Close()
	//log.Println(resp)

	var body []WarehousingTransaction
	err1 := json.NewDecoder(res.Body).Decode(&body)
	if err1 != nil {
		log.Println(err1)
	}

	if body == nil { //if the response is null then return ture. It means the ProductCode can be used.
		result = true
	}
	//log.Println(result)

	return result, body
}

func CheckProductBatchNo(ProductID string, ProductBatch string) (bool, []WarehousingTransaction) { // check the ProductBatchNo with blockchain nodes
	result := true
	receivedProductID := ProductID
	receivedProductBatch := ProductBatch
	var querybyProductCodeAndBatch []WarehousingTransaction

	//var querybyProductCode []DeliveryTransaction
	status, _ := CheckProductCode(receivedProductID)
	if status == true {
		return result, querybyProductCodeAndBatch
	}

	status_2, querybyProductCode := CheckProductIDforBatch(receivedProductID)
	if status_2 == true {
		return result, querybyProductCode
	}
	//log.Println(querybyProductCode)

	querybyProductCodeLen := len(querybyProductCode)

	for i := 0; i < querybyProductCodeLen; i++ {

		if querybyProductCode[i].ProductBatch == receivedProductBatch {
			result = false
			querybyProductCodeAndBatch = append(querybyProductCodeAndBatch, querybyProductCode[i])

		}
	}

	return result, querybyProductCodeAndBatch
}

func CheckQuantity(ProductCode string, ProductBatchNo string, Quantity int) (bool, WarehousingTransaction) { // check the ProductCode with blockchain nodes
	var WarehousingTransaction WarehousingTransaction
	result := false
	receivedProductCode := ProductCode
	receivedProductBatchNo := ProductBatchNo
	receivedProductBatchQuantity := Quantity

	status, querybyProductCodeAndBatchNo := CheckProductBatchNo(receivedProductCode, receivedProductBatchNo)
	if status == true {
		return result, WarehousingTransaction
	}

	Txlength := len(querybyProductCodeAndBatchNo)

	if receivedProductBatchQuantity == querybyProductCodeAndBatchNo[Txlength-1].Quantity {
		result = true
	}

	return result, querybyProductCodeAndBatchNo[Txlength-1]

}

func GenerateNewProductTx(newProductInfo NewProductTransaction) { //Post this new transaction
	posturl := bcurl + "/new"
	//posturl := "http://localhost:8880/raw"
	log.Println(newProductInfo)
	var NewDeliveryTx NewProductTransaction
	NewDeliveryTx = newProductInfo

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

func GenerateWarehousingTransaction(newProductInfo WarehousingTransaction) { //Post this new transaction
	posturl := bcurl + "/war"
	//posturl := "http://localhost:8880/raw"
	log.Println(newProductInfo)
	var NewDeliveryTx WarehousingTransaction
	NewDeliveryTx = newProductInfo

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

func GenerateNewTxbyDelivery(newDeliveryInfo DeliveryTransaction) { //Post this new transaction
	posturl := bcurl + "/del"
	//posturl := "http://localhost:8880/raw"
	log.Println(newDeliveryInfo)
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
