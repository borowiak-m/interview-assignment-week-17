FROM golang:latest 
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o interview-assignment-week-17 .
CMD ["./interview-assignment-week-17"]