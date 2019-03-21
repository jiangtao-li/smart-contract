# How to use

API can be called with a HTTP request:
```
 request line:  <method type> ip:port/<API name>
 body: <body>
```

We always use the  `POST` method to post a transaction; and `GET` method to query something.

For now we only have one query API with name: `Inventoryquery` .

Here are examples:

 - `POST  192.168.50.184:9000/ProductDeclaration` with well formatted body will generate a type 1 transaction.
 - `GET  192.168.50.184:9000/Inventoryquery` will return all the inventory items.

# Transaction Structures

#### First smart contract: Product declaration
> Let users key in the product details and raw materials: 

> The smart contract to check for duplicates of product Batch number. Product batch number has to be unique. Raise error when Product Batch ID has been used. The export will call this API when he have manufactured a new batch of products.
Readable by Importer. 

**API name:**  `ProductDeclaration` 

**body format:**


```
  {
        "SerialNo": 0,
        "ProductCode": "",
        "ProductName": "",
        "ProductBatchNo": "",
        "Quantity": 0,
        "RawMaterial": 
        [
            {
                "RawMaterialBatchNo": "",
                "RawMaterialsID": "",
                "RawMaterialName": "",
                "RawMaterialQuantity": 0,
                "RawMaterialMeasurementUnit": ""
            }
        ]
    }


```

#### Second smart contract: Shipping Batch declaration

> Let users create shipment based on PO.Group different batches of product under a single shipment ID. The contract is called by the export to the importer.

> The smart contract checks for the quantity inputted by the exporter against current stock of the batches. At the same time, the export uploads the documents details of the PO. The importer will check through and once verified, will sign off with a signature to complete the transaction.

**API name:**  `ShippingBatchDeclaration` 

**body format:**

```
   {
        "SerialNo": 0,
        "RecGenerator": "",
        "ShipmentID": "",
        "Timestamp": "",
        "Longitude": "",
        "Latitude": "",
        "ShippedFromCompID": "",
        "ShippedToCompID": "",
        "LocationID": "",
        "DeliveryStatus": "",
        "DeliveryType": "",
        "Product": 
        [
        	{
                "ProductCode": "",
                "ProductName": "",
                "ProductBatch": {
                    "ProductBatchNo": "",
                    "ProductBatchQuantity": 0
                }
            }
        ],
        "Document": {
            "DocumentURL": "",
            "DocumentType": "",
            "DocumentHash": "",
            "DocumentSign": ""
        },
        "UserSign":{
			"User":     "",
			"Verify":   true,
			"UserSign": ""
		}
    }



```

#### Fourth Smart: MSDS and DG upload

> Let exporter upload required DG and MSDS. 

**API name:**  `DocumentsUpload` 

**body format:**

```
 {
	"SerialNo":        0,
	"RecGenerator":      "",
	"ShipmentID":       "",
	"Timestamp":     "",
	"ShippedFromCompID":  "",
	"ShippedToCompID": "",
	"DeliveryType":      "",
     "Document": {
            "DocumentURL": "",
            "DocumentType": "",
            "DocumentHash": "",
            "DocumentSign": ""
        },
     "UserSign":{
			"User":     "",
			"Verify":   true,
			"UserSign": ""
		}

}
```

#### Fifth Smart: Ownership change

> Let Transfer goods to the Freight forwarder. Require Freight forwarder to acknowledge and upload bill of laden. The smart contract has to check for the Shipping documents i.e. the transaction for the shipping document must be valid.  

**API name:**  `OwnershipChange` 

**body format:**


```
{
	"SerialNo":        0,
	"RecGenerator":      "",
	"ShipmentID":        "",
	"Timestamp":         "",
	"ShippedFromCompID": "",
	"ShippedToCompID":   "",
	"DeliveryType":      "",
	"UserSign":{
        "User":     "",
        "Verify":   true,
        "UserSign": "",
        "Document": {
            "DocumentURL":  "",
            "DocumentType": "",
            "DocumentHash": "",
            "DocumentSign": ""
        }
	}
}
```


### Build Notes 

Refer to (build readme)[docker-build.md] for details
