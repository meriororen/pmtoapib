all:
	go build
	./pmtoapib -force-apib -force-responses -c Bill*.json -environment-path Staging_GCP.postman_environment.json

copy:
	cp pmtoapib ../../../nodejs/api-blueprints/
