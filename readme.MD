Methods available are: http://localhost:port/

http://localhost:port/newProduct      //generate new ProdectCode

http://localhost:port/newGenerateBatchNo  //generate new ProductBatchNo and ProductBatchQuantity with existing ProdectCode

http://localhost:port/PurchaseOrderRegistry //Post DocumentURL and DocumentHash

http://localhost:port/checkpoacknowledgement/{DocumentURL}/{DocumentHash}  
http://localhost:port/poacknowledgement  //post signature


Transaction Structure
#### For newProduct and newGenerateBatchNo methods, input format should be
```

  {
        "ProductCode": ""
        "ProductName": ""  
        "ProductBatch": [{
            "ProductBatchNo": "",
            "RawMaterialQuantity": 0
        }]
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
        "Product": [{
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

