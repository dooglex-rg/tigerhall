# Project Directory Map

## Python Folder
Python folder contains django app which is used for automatic database creation / migrations.
Django makes database migrations easier.

## Golang Folder
Golang folder contains the main webapp which will be handling the incoming API requests.

## golang/images
this is where user uploaded images are stored and resized later (by task queue)

## API documentation
The API is documented based on swagger standards. The documenation can be found in urlpath `/api-docs/`
The deployed webapp/docs is available in [this url](https://tigerhall.dooglex.com/api-docs/)

# how to run this project (makefile commands)

## prequsites/installation
configure the .env file for your local settings

```sh
make setup
```
this command will setup/migrate the database schemas based on python/django models

```sh
make test
```
runs go testing of our webapp

```sh
make app
```
initiates the webapp

```sh
make consumer
```
initiates the task queue for background processing of image resizing

## important note:
Both `make app` and `make consumer` are persistant process. both need to be running parallell in order to run the task queue along with the webapp.