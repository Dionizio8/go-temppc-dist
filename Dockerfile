FROM golang:1.22.5 as build
WORKDIR /app
COPY . .
ENV WEB_SERVER_PORT=${WEB_SERVER_PORT}
ENV OTEL_SERVICE_NAME=${OTEL_SERVICE_NAME}
ENV OTEL_COLLECTOR_URL=${OTEL_COLLECTOR_URL}
ARG APP_NAME
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/$APP_NAME/

FROM scratch
WORKDIR /app
COPY --from=build /app/app .
COPY --from=build /app/cmd/.env .

EXPOSE 8080
EXPOSE 8090

ENTRYPOINT ["./app"] 