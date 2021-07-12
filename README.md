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
  * /cart/list      POST
  * /cart/{id}      GET
  * /cart           POST
  * /cart/{id}      PUT
  * /cart/push/{id} PUT
  * /cart/pop/{id}  PUT

# kubernetes
kube directory contains kubernetes deployment yaml. It pull from hub.docker.com/kharism/ repository.
You can use it using ```kubectl apply -f <your yaml file here>```

# testing
in each directory, except model there are already test file. It uses api_test.json on config directory
Just cd into that directory and execute ```go test``` to test the package. Make sure your database is running first
