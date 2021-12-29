# microservice_simple

this project has 3 simple service 

- auth on port 8098 for authentication
- item on port 8099 for item management
- cart on port 8100 for transaction

all of those services will use mongodb on localhost:27017. Bydefault they dont authenticate the database. Edit their config file if you have to use some sort of authentication method

To dockerize app run on the root directory
```
docker build -t <image name> ./webservice/<auth_api/transaction_api/ui_api>/.
```

# Rest API
- auth
  * /auth/registeruser POST
  * /auth POST
- item
  * /item/list POST
  * /item/{id} GET
  * /item/     POST
  * /item/{id} PUT
  * /item/{id} DELETE
- cart
  /cart/list      POST
  /cart/{id}      GET
  /cart/checkout/{id}
  /cart           POST
  /cart/{id}      PUT
  /cart/push/{id} PUT
  /cart/pop/{id}  PUT

# kubernetes
kube directory contains kubernetes deployment yaml. It pull from hub.docker.com/kharism/ repository.
You can use it using ```kubectl apply -f <your yaml file here>```

# testing
in each directory, except model there are already test file. It uses api_test.json on config directory
Just cd into that directory and execute ```go test``` to test the package. Make sure your database is running first

Since We can't do transaction in mongo, we use single routine to handle transaction. Any transaction on cart/checkout will go through single transaction so we can ensure atomicity. In this sample project the transaction order is stored in memory, we should store it in something that is persistent to ensure we don't loose message in case something down. 

In real life production, we should use message queus to make the order persistent and can be consumed by transaction processor. If the order can be handled normally then it goes through. If one of the order cannot be fulfilled then we send notice to the client that the order cannot be fulfilled.

the basic line of process is
- customer checkout their order
- the system notify the client that their order will be processed
- if the order can be processed then continue as usual (wait for customer payment, packing order, send the order to expedition company)
- if the order can't be processed then notify the order can't be fulfilled

