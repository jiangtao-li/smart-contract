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

type RawMaterialTransaction struct {
	SerialNo       int    `json:"SerialNo"`
	ProductCode    string `json:"ProductCode"`
	ProductName    string `json:"ProductName"`
	ProductBatchNo string `json:"ProductBatchNo"`
	Quantity       int    `json:"Quantity"`
	//TxnTimestamp string `json:"TxnTimestamp"`
	RawMaterialInfo []RawMaterial
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
	Product           []Product
	Document          []Document
}

type Product struct {
	ProductCode  string `json:"ProductCode"`
	ProductName  string `json:"ProductName"`
	ProductBatch []ProductBatch
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

type Block struct { // An element of Blockchain
	Index      int
	Timestamp  string
	TxnType    int
	TxnPayload string
	Comment    string
	Proposer   string
	PrevHash   string
	ThisHash   string
}

func CheckSerialNo(SerialNo int) bool { // check the SerialNo with blockchain nodes
	result := false
	receivedSerialNo := strconv.Itoa(SerialNo)
	url := "http://172.21.177.60:5000/query/raw/SerialNo/" + receivedSerialNo
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

func CheckProductCode(ProductCode string) (bool, []DeliveryTransaction) { // check the ProductCode with blockchain nodes
	//log.Println("called")
	result := false
	receivedProductCode := ProductCode
	url := "http://172.21.177.60:5000/query/del/Product.ProductCode/" + receivedProductCode
	log.Println(url)

	resp, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}
	res, _ := http.DefaultClient.Do(resp)
	defer res.Body.Close()
	//log.Println(resp)

	var body []DeliveryTransaction
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

func CheckProductBatchNo(ProductCode string, ProductBatchNo string) (bool, []DeliveryTransaction) { // check the ProductBatchNo with blockchain nodes
	result := true
	receivedProductCode := ProductCode
	receivedProductBatchNo := ProductBatchNo

	//var querybyProductCode []DeliveryTransaction
	status, querybyProductCode := CheckProductCode(receivedProductCode)
	if status == true {
		return result, querybyProductCode
	}
	//log.Println(querybyProductCode)

	querybyProductCodeLen := len(querybyProductCode)
	var querybyProductCodeAndBatchNo []DeliveryTransaction

	for i := 0; i < querybyProductCodeLen; i++ {
		for j := 0; j < len(querybyProductCode[i].Product); j++ {
			for k := 0; k < len(querybyProductCode[i].Product[j].ProductBatch); k++ {
				if querybyProductCode[i].Product[j].ProductBatch[k].ProductBatchNo == receivedProductBatchNo {
					result = false
					querybyProductCodeAndBatchNo = append(querybyProductCodeAndBatchNo, querybyProductCode[i])
				}
			}
		}
	}

	return result, querybyProductCodeAndBatchNo
}

func CheckQuantity(ProductCode string, ProductBatchNo string, ProductBatchQuantity int) (bool, DeliveryTransaction) { // check the ProductCode with blockchain nodes
	var DeliveryNewTransaction DeliveryTransaction
	result := false
	receivedProductCode := ProductCode
	receivedProductBatchNo := ProductBatchNo
	receivedProductBatchQuantity := ProductBatchQuantity

	status, querybyProductCodeAndBatchNo := CheckProductBatchNo(receivedProductCode, receivedProductBatchNo)
	if status == true {
		return result, DeliveryNewTransaction
	}

	Txlength := len(querybyProductCodeAndBatchNo)
	Productlen := len(querybyProductCodeAndBatchNo[Txlength-1].Product)
	ProductBatchlen := len(querybyProductCodeAndBatchNo[Txlength-1].Product[Productlen-1].ProductBatch)

	if receivedProductBatchQuantity == querybyProductCodeAndBatchNo[Txlength-1].Product[Productlen-1].ProductBatch[ProductBatchlen-1].ProductBatchQuantity {
		result = true
	}

	return result, querybyProductCodeAndBatchNo[Txlength-1]

}

func GenerateNewProductTx(newProductInfo Product) { //Post this new transaction
	posturl := "http://172.21.177.60:5000/del"
	//posturl := "http://localhost:8880/raw"
	log.Println(newProductInfo)
	var NewDeliveryTx DeliveryTransaction
	NewDeliveryTx.Product = append(NewDeliveryTx.Product, newProductInfo)

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
	posturl := "http://172.21.177.60:5000/del"
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
