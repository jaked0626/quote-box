#!/bin/bash 

if [ -f app.env ]
then
  export $(cat app.env | xargs)
fi

wgo run ./cmd/web -addr=$HTTP_TARGET_ADDRESS -dsn=$DB_SOURCE

