# Go Simple Rest

This package provides an out of the box implementation of a rest api router for go.

The database is written for sql but can be easily modified to work with other.

It uses gorilla mux as router and logrus as logger.

The api will create endpoints for each resource provided to the AddHandlers function.

It creates the following routes (models are table names of the resources, converted to kebab case):
- GET /model/{id}
- POST /model
- PUT /model
- DELETE /model/{id}
  
Also, for each belongs to relation, it creates the following routes:
- GET /belongs-to/{id}/model
  
The handlers parameter is a function to wrap the handlers with, for example, authentication and logging
