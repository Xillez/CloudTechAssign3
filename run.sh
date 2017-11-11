#!/bin/sh
go install -v ./cmd/clock
go install -v ./cmd/web
heroku local