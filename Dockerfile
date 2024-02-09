FROM golang:1.19

WORKDIR /ovalfi

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ovalfi ./cmd/ovalfi

EXPOSE 8080

CMD [ "./ovalfi" ]