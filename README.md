# MTA Server **WIP**

THIS IS A WORK IN PROGRESS

Go webserver that translates data from the [MTA realtime subway feeds](http://datamine.mta.info/) into a RESTful JSON API.

## Endpoints

GET /stations?query="substring"

Returns all stations with names matching substring.

GET /stations/:id/arrivals

Returns expected arrivals at station specified.

## First time setup

Set up an account on the MTA website to get an API key. Then run:
```
echo 'ApiKey = "YOUR API KEY"' > api.ini
```
