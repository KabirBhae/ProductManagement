# ProductManagement

Simple prouctmanagement project. Used golang and echo framework for the backend and mongodb as database. Routes are protected using JWT.

Kindly create a .env file in the root folder and add two variables named "MONGOURI" and "SECRET_KEY". The first variable should be a valid mongodb connection string that you will use to connect your mongo database. (All collections will automatically be created when you try to visit the routes). The second variable is a string you want to encrypt your JWT token with.

### example .env file:

MONGOURI = mongodb+srv://username:>password>@cluster0.d2cphw.mongodb.net/?retryWrites=true&w=majority

SECRET_KEY = "water"

Start reading from the main.go file. To get a summary of what the project does, see the controller files under "controllers" folder.
