FROM golang:1.21

WORKDIR /auth-app

COPY go.mod go.sum ./
RUN go mod download
# COPY data/* ./data/
COPY . /auth-app/

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth

EXPOSE 8080
CMD [ "/auth" ]