#!/bin/bash

make postgres
sleep 2

make createdb
sleep 1

make migrateup
sleep 1

make createdbrole
sleep 1

make grantpermission
sleep 1

make serve