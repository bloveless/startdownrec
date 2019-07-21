FUNCTION_NAME := startdownrec
FUNCTION_RUNTIME := go111
FUNCTION_ENTRY_POINT := Run
FUNCTION_SERVICE_ACCOUNT := startdownrec-function@brennon-loveless.iam.gserviceaccount.com
HOSTNAME := $(shell hostname)
BIN_DIR=$(PWD)/bin

.PHONY: all deploy
all: build

dependencies:
	go mod download

build: dependencies
	go build -o $(BIN_DIR)/startdownrec cmd/*.go

run: build
	$(BIN_DIR)/startdownrec

install-reflex:
	go get github.com/cespare/reflex

debug: dependencies install-reflex
	reflex -c reflex.conf

# TESTING/DEPLOYMENT TOOLS

post:
	curl -X POST -H "Content-Type: application/json" -d '{"hostname":"$(HOSTNAME)","status":"startup"}' http://localhost:8080

post-prod:
	curl -X POST -H "Content-Type: application/json" -d '{"hostname":"$(HOSTNAME)","status":"startup"}' https://us-central1-brennon-loveless.cloudfunctions.net/startdownrec

deploy:
	gcloud functions deploy ${FUNCTION_NAME} \
		--runtime=${FUNCTION_RUNTIME} \
		--entry-point=${FUNCTION_ENTRY_POINT} \
		--service-account=${FUNCTION_SERVICE_ACCOUNT} \
		--trigger-http
