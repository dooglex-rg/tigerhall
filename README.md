# Project Directory Map

## Python Folder
Python folder contains django app which is used for automatic database creation / migrations.
Django makes database migrations easier.

## Golang Folder
Golang folder contains the main webapp which will be handling the incoming API requests.

## golang/images
this is where user uploaded images are stored and resized later (by task queue)

# how to run this project (makefile commands)

## prequsites
configure the .env file for your local settings

## make setup
this command will setup/migrate the database schemas based on python/django models

## make test
runs go testing of our webapp

## make app
initiates the webapp

## make consumer
initiates the task queue for background processing of image resizing

## important note:
Both "make app" and "make consumer" are persistant process. both need to be running parallell in order to run the task queue along with the webapp.