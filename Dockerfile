FROM golang:1.21.6 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/tb

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /bin/tb /bin/tb
COPY --from=build-stage /app/conf.yml /etc/tb.yml

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/bin/tb"]
