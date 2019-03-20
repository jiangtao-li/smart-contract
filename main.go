package main

import (
	"io"
	//"os"
	"encoding/gob"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

///// GLOBAL FLAGS & VARIABLES

var StartTime time.Time

var listenPort *int // listen port & total locations in the supply chain

///// LIST OF FUNCTIONS

func init() {

	gob.Register(Product{})
	gob.Register(map[string]interface{}{})

	log.SetFlags(log.Lshortfile)

	log.Printf("========================================")
	listenPort = flag.Int("port", 9000, "mux server listen port")
	//logDir = flag.String("pdts", "pdts", "pathname of log data directory")

	flag.Parse()

	StartTime = time.Now()
	StartTime = StartTime.AddDate(0, -6, 10) // random negative offset

}

func main() {
	log.Fatal(launchMUXServer())
}

func launchMUXServer() error { // launch MUX server
	mux := makeMUXRouter()
	log.Println("HTTP Server Listening on port:", *listenPort) // listenPort is a global flag
	s := &http.Server{
		Addr:           ":" + strconv.Itoa(*listenPort),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMUXRouter() http.Handler { // create handlers
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleHome).Methods("GET")

	muxRouter.HandleFunc("/ProductDeclaration", handleProductDeclaration).Methods("POST") //post a new product
	muxRouter.HandleFunc("/ShippingBatchDeclaration", handleShippingBatchDeclaration).Methods("POST")
	muxRouter.HandleFunc("/DocumentsUpload", handleDocumentsUpload).Methods("POST")
	muxRouter.HandleFunc("/OwnershipChange", handleOwnershipChange).Methods("POST")

	return muxRouter
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	log.Println("handleHome() API called")
	io.WriteString(w, " ")
}

func handleProductDeclaration(w http.ResponseWriter, r *http.Request) { //handle new product request
	log.Println("called")
	w.Header().Set("Content-Type", "application/json")
	log.Println("handleNewProduct() API called")
	var receivedNewTx RawMaterialTransaction

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&receivedNewTx); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	if receivedNewTx.ProductCode == "" {
		respmsg := "Emptry Product Code!"
		bytes, _ := json.MarshalIndent(respmsg, "", "  ")
		io.WriteString(w, string(bytes))
		return
	}

	if receivedNewTx.ProductBatchNo == "" {
		respmsg := "Emptry Product Batch Number!"
		bytes, _ := json.MarshalIndent(respmsg, "", "  ")
		io.WriteString(w, string(bytes))
		return
	}

	//newProduct.TxnTimestamp = time.Now().Format("02-01-2006 15:04:05 Mon")
	log.Println("New transaction received:", receivedNewTx)

	result, _ := CheckProductCodeBatch(receivedNewTx.ProductCode, receivedNewTx.ProductBatchNo)

	//result, _ := CheckProductBatchNo(receivedNewTx.ProductBatchNo) // check if the SerialNo exists
	var respmsg string

	if result == true {
		GenerateNewRawMaterialTx(receivedNewTx)
		//log.Println(result)
		respmsg = "Transaction has been posted, please wait for comfirmation"
		bytes, _ := json.MarshalIndent(respmsg, "", "  ")
		io.WriteString(w, string(bytes))
	} else {
		//log.Println(result)
		respmsg = "ProductBatch existed!"
		bytes, _ := json.MarshalIndent(respmsg, "", "  ")
		io.WriteString(w, string(bytes))
	}

}

func handleShippingBatchDeclaration(w http.ResponseWriter, r *http.Request) { //handle new product request
	w.Header().Set("Content-Type", "application/json")
	log.Println("handleGenerateBatchNo() API called")
	var receivedNewTx DeliveryTransaction
	var respmsg string

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&receivedNewTx); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	productLen := len(receivedNewTx.Product)

	for i := 0; i < productLen; i++ {
		if receivedNewTx.Product[i].ProductBatch.ProductBatchNo == "" || receivedNewTx.Product[i].ProductCode == "" || receivedNewTx.Product[i].ProductBatch.ProductBatchQuantity == 0 {
			respmsg := "Error, Emptry item!"
			bytes, _ := json.MarshalIndent(respmsg, "", "  ")
			io.WriteString(w, string(bytes))
			return
		}
	}

	log.Println("New product information received:", receivedNewTx)

	for j := 0; j < productLen; j++ {
		result, queryresult := CheckProductBatchNonCode(receivedNewTx.Product[j].ProductCode, receivedNewTx.Product[j].ProductBatch.ProductBatchNo) // check if the SerialNo exists
		if result == false {
			respmsg = "The product batch number or product code you provided is incorrect!"
			bytes, _ := json.MarshalIndent(respmsg, "", "  ")
			io.WriteString(w, string(bytes))
		} else {
			result2 := CheckProductQuantity(queryresult, receivedNewTx.Product[j].ProductBatch.ProductBatchQuantity)
			if result2 == false {
				respmsg = "The product batch quantity you provided is incorrect!"
				bytes, _ := json.MarshalIndent(respmsg, "", "  ")
				io.WriteString(w, string(bytes))
			} else {
				if receivedNewTx.UserSign.Verify == true {
					//newTx.ProductBatch[0].ProductBatchQuantity = receivedNewTx.ProductBatch[0].ProductBatchQuantity
					GenerateNewDeliveryTx(receivedNewTx)
					//log.Println(result)
					respmsg = "Transaction has been posted, please wait for comfirmation"
					bytes, _ := json.MarshalIndent(respmsg, "", "  ")
					io.WriteString(w, string(bytes))
				}
			}
		}
	}

}

func handleDocumentsUpload(w http.ResponseWriter, r *http.Request) { //handle new product request
	log.Println("called")
	w.Header().Set("Content-Type", "application/json")
	log.Println("handleDocumentsUpload() API called")
	var receivedNewTx ShippingDocTransaction

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&receivedNewTx); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	if receivedNewTx.Document.DocumentURL == "" {
		respmsg := "Empty Document!"
		bytes, _ := json.MarshalIndent(respmsg, "", "  ")
		io.WriteString(w, string(bytes))
		return
	}

	if receivedNewTx.Document.DocumentType == "" {
		respmsg := "Empty Document Type!"
		bytes, _ := json.MarshalIndent(respmsg, "", "  ")
		io.WriteString(w, string(bytes))
		return
	}

	//newProduct.TxnTimestamp = time.Now().Format("02-01-2006 15:04:05 Mon")
	log.Println("New transaction received:", receivedNewTx)
	var respmsg string

	GenerateNewShippingDocTransaction(receivedNewTx)
	respmsg = "Transaction has been posted, please wait for comfirmation"
	bytes, _ := json.MarshalIndent(respmsg, "", "  ")
	io.WriteString(w, string(bytes))

}

func handleOwnershipChange(w http.ResponseWriter, r *http.Request) { //handle new product request
	log.Println("called")
	w.Header().Set("Content-Type", "application/json")
	log.Println("handleNewProduct() API called")
	var receivedNewTx OwnerTransaction

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&receivedNewTx); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	if receivedNewTx.UserSign.Verify == false {
		respmsg := "Invalid Signature!"
		bytes, _ := json.MarshalIndent(respmsg, "", "  ")
		io.WriteString(w, string(bytes))
		return
	}

	if receivedNewTx.UserSign.User == "" {
		respmsg := "Empty freight forwarder name!"
		bytes, _ := json.MarshalIndent(respmsg, "", "  ")
		io.WriteString(w, string(bytes))
		return
	}

	//newProduct.TxnTimestamp = time.Now().Format("02-01-2006 15:04:05 Mon")
	log.Println("New transaction received:", receivedNewTx)

	result, _ := CheckShipmentID(receivedNewTx.ShipmentID)

	//result, _ := CheckProductBatchNo(receivedNewTx.ProductBatchNo) // check if the SerialNo exists
	var respmsg string

	if result == true {
		GenerateNewOwenerTransaction(receivedNewTx)
		//log.Println(result)
		respmsg = "Transaction has been posted, please wait for comfirmation"
		bytes, _ := json.MarshalIndent(respmsg, "", "  ")
		io.WriteString(w, string(bytes))
	} else {
		//log.Println(result)
		respmsg = "ProductBatch existed!"
		bytes, _ := json.MarshalIndent(respmsg, "", "  ")
		io.WriteString(w, string(bytes))
	}

}
