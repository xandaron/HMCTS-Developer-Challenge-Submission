FROM golang:1.24.2 AS build

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

FROM gcr.io/distroless/static:nonroot

COPY --from=build /server .
COPY ./static ./static
COPY ./templates ./templates

EXPOSE 80
CMD ["./server"]