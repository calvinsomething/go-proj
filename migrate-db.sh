#! /bin/bash

mode=$1
option=$2
# [ verbose ]
docker compose up -d db && docker compose run server ./wait-for db:3306 -t 30 -- go run . migrate ${mode} ${option}