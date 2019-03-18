Methods available are: http://ip:port/ 
url can be changed in func.go

For example, Methods available are: http://localhost:9000/
http://localhost:port/ProductDeclaration //Product declaration
     
http://localhost:port/ShippingBatchDeclaration   // for Shipping Batch declaration


Transaction Structure
#### First smart contract: Product declaration
Let users key in the product details and raw materials:
The smart contract to check for duplicates of product Batch number. Product batch number has to be unique. Raise error when Product Batch ID has been used. The export will call this API when he have manufactured a new batch of products.
Readable by Importer. 


```
   {
        "SerialNo": 0,
        "ProductCode": "",
        "ProductName": "",
        "ProductBatchNo": "",
        "Quantity": 0,
        "RawMaterial": {
            "RawMaterialBatchNo": "",
            "RawMaterialsID": "",
            "RawMaterialName": "",
            "RawMaterialQuantity": 0,
            "RawMaterialMeasurementUnit": ""
        }
    }


```

#### Second smart contract: Shipping Batch declaration

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
        "Product": {
            "ProductCode": "",
            "ProductName": "",
            "ProductBatch": {
                "ProductBatchNo": "",
                "ProductBatchQuantity": 0
            }
        },
        "Document": {
            "DocumentURL": "",
            "DocumentType": "",
            "DocumentHash": "",
            "DocumentSign": ""
        }
    }



```




### Build Notes 

Refer to (build readme)[docker-build.md] for details
