FROM mirror.gcr.io/gocv/opencv:4.8.0 as builder
ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o ./run ./cmd/service

ENTRYPOINT ["./run", "-config_path",  "/app/configs/config.yml"]