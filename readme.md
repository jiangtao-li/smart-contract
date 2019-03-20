Methods available are: http://ip:port/ 
url can be changed in func.go

For example, Methods available are: http://localhost:9000/
http://localhost:port/ProductDeclaration //Product declaration
     
http://localhost:port/ShippingBatchDeclaration   // for Shipping Batch declaration


Transaction Structure
#### First smart contract: Product declaration
> Let users key in the product details and raw materials: 

> The smart contract to check for duplicates of product Batch number. Product batch number has to be unique. Raise error when Product Batch ID has been used. The export will call this API when he have manufactured a new batch of products.
Readable by Importer. 

`ProductDeclaration` with body


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

`ShippingBatchDeclaration` with body

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

`DocumentsUpload` with body

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

`OwnershipChange` with body


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
