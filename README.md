# Secret Santa

## Introduction
This app helps you to to organize groups of "secret santa", 
it provides a totally anonymous and mail free way to draw lots and keep the lots hidden from each other. 
Every player will see only the one he/she is gifting to. Registration is done by forename and a unique game id. 
You can also add exceptions if couples shouldn't gift each other. A more detailed description will be added later and in the app itself.

## Second use
This app also serves as a simple sample app with a golang backend, Angular frontend, nginx reverse proxy.

## Prerequisites
- git ;-)
- Docker

## Use
- get repo ;-)
- copy .default.env to .env
- change the parameters in this file, they will be used to create the docker containers and as runtime parameters for the apps
  - DB_\*: Cpt Obvious' parameters, used for creating the postgres container and connecting to the db from the backend
  - COOKIE_SECRET: A cookie secret to encrypt session cookies
  - GIN_MODE: "release" should be fine for normal operation mode, you can also set it to debug to trace bugs
  - ALLOWED_HOSTS: the host url of the webapp, used for AllowOrigins (for CORS)
  - REST_SERVICE_URL: the url of the backend, used by the webapp to make rest calls
- run "docker-compose up"
