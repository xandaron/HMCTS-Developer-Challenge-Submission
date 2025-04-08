FROM golang:1.24.2 AS build

WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

FROM gcr.io/distroless/static:nonroot

COPY --from=build /backend/server .

EXPOSE 8080
CMD ["./server"]