#! /bin/bash

mode=$1
docker compose up -d db && docker compose run server ./wait-for db:3306 -t 30 -- go run . migrate ${mode}