Methods available are: http://localhost:port/

http://localhost:port/newProduct      //generate new ProdectCode

http://localhost:port/newGenerateBatchNo  //generate new ProductBatchNo and ProductBatchQuantity with existing ProdectCode

http://localhost:port/PurchaseOrderRegistry //Post DocumentURL and DocumentHash

http://localhost:port/checkpoacknowledgement/{DocumentURL}/{DocumentHash}  
http://localhost:port/poacknowledgement  //post signature


Transaction Structure
#### For newProduct (Transaction type 3), input format should be

```
   {
    "SerialNo":        0,
    "ProductID":       "",
    "ProductName":     "",
    "RawMaterialsID":  "",
    "RawMaterialName": ""
}

```

#### For newGenerateBatchNo (Transaction type 4)

```
{
    "SerialNo":        0,
    "ProductID":       "",
    "ProductName":     "",
    "ProductBatch":    "",
    "RawMaterialsID":  "",
    "RawMaterialName": "",
    "MaterialBatch":   "",
    "Quantity":        1
}


```


#### For PurchaseOrderRegistry and poacknowledgement methods
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
        }],
        "Document": [{
            "DocumentURL": "",
            "DocumentType": "",
            "DocumentHash": "",
            "DocumentSign": ""
        }]
  }
```

### Build Notes 

Refer to (build readme)[docker-build.md] for details
