NAME	=	nobbydobby
LINKERFLAGS = -X main.Name=$(NAME) -X main.Version=`git describe --tags --always --dirty` -X main.BuildTimestamp=`date -u '+%Y-%m-%d_%I:%M:%S_UTC'`
PROJECTROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
DBPATH	=	$(PROJECTROOT)db/



all: clean build

.PHONY: clean release
clean:
	@echo Running clean job...
	rm -f coverage.txt
	rm -rf bin/ release/
	rm -f main

generate:
	@echo Running generate job...


build: #dep generate
	@echo Running build job...
	mkdir -p bin/linux/arm bin/linux/x64 bin/windows bin/mac/x64 bin/mac/arm
	GOOS=linux GOARCH=arm64 go build  -ldflags "$(LINKERFLAGS)" -o bin/linux/arm ./...
	GOOS=linux GOARCH=amd64 go build  -ldflags "$(LINKERFLAGS)" -o bin/linux/x64 ./...
	GOOS=windows GOARCH=amd64 go build  -ldflags "$(LINKERFLAGS)" -o bin/windows ./...
	GOOS=darwin GOARCH=amd64 go build  -ldflags "$(LINKERFLAGS)" -o bin/mac/x64 ./...
	GOOS=darwin GOARCH=arm64 go build  -ldflags "$(LINKERFLAGS)" -o bin/mac/arm ./...

run: #generate
	go run -ldflags "$(LINKERFLAGS)" ./...

test: recreatetables
#	go run -ldflags "$(LINKERFLAGS)" main.go -cfg cselo-local.ini -import data/test.log
	@echo Running test job...
	go test ./... -cover -coverprofile=coverage.txt

coverage: test
	@echo Running coverage job...
	go tool cover -html=coverage.txt

datacreate:
	./scripts/createdata.sh

dataread:
	curl http://localhost:3000/


initdb:
	rm -rf $(DBPATH)
	mkdir -p $(DBPATH)
	initdb -D $(DBPATH)
	make startdb
	sleep 2
	psql postgres -f scripts/create-db.sql
	make recreatetables

startdb:
	postgres -D $(DBPATH)
	#osascript -e 'tell app "Terminal" to do script "postgres -D $(DBPATH)"'
	#execInNewITerm "postgres -D $(DBPATH)"

stopdb:
	pg_ctl stop -D $(DBPATH) -m fast

recreatetables:
	psql $(NAME) -U $(NAME)app -f scripts/create-tables.sql
