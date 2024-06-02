GEN_DIR=internal/gen
SERVICE_NAME = predictor-service

swagger-gen:
	if ! [ -d $(GEN_DIR) ]; then \
	    mkdir $(GEN_DIR); \
	elif [ -d $(GEN_DIR) ]; then \
		rm -rf $(GEN_DIR); \
		mkdir $(GEN_DIR); \
	fi && \
	swagger generate server -t internal/gen -f ./api/swagger.yml --exclude-main -A $(SERVICE_NAME) && \
	go mod tidy && \
	git add $(GEN_DIR)

go-build+run-dev:
	go build ./cmd/service && \
	sudo ./service -config_path ./configs/dev.yml