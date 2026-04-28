FROM golang:1.25 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o /krigerforaktforliv

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

VOLUME ["/data"]

COPY --from=build-stage /krigerforaktforliv /krigerforaktforliv

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/krigerforaktforliv"]
