FROM golang:1.24.2 AS build

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY api ./api/
COPY database ./database/
COPY errors ./errors/
COPY session ./session/

RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=build /server ./
COPY ./static ./static
COPY ./templates ./templates
COPY ./certs ./certs
COPY ./README.md ./README.md

USER nonroot:nonroot

EXPOSE 80
EXPOSE 443
CMD ["./server"]