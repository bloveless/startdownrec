GOOGLE_APPLICATION_CREDENTIALS := ~/.gcp/local-startdownrec-service-account.json
FUNCTION_NAME := startdownrec
FUNCTION_RUNTIME := go111
FUNCTION_ENTRY_POINT := Run
FUNCTION_SERVICE_ACCOUNT := startdownrec-function@brennon-loveless.iam.gserviceaccount.com
HOSTNAME := $(shell hostname)

run:
	GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS} go run cmd/main.go

post:
	curl -X POST -H "Content-Type: application/json" -d '{"hostname":"$(HOSTNAME)","status":"startup"}' http://localhost:8080

deploy:
	gcloud functions deploy ${FUNCTION_NAME} \
		--runtime=${FUNCTION_RUNTIME} \
		--entry-point=${FUNCTION_ENTRY_POINT} \
		--service-account=${FUNCTION_SERVICE_ACCOUNT} \
		--set-env-vars=MYSQL_HOST=
		--trigger-http
