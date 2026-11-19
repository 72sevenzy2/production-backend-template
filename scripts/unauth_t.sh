#!/bin/bash

cd .. || exit 1

go test -run '^TestGetUser|TestEditUser|TestDeleteUser|TestCreateUser' -v
