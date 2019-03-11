# Pull base image
FROM golang:latest 

RUN apt-get update && apt-get upgrade -y && apt-get autoremove && apt-get autoclean

# Set environment variables
ARG PROJECT=smart-contract
ARG PROJECT_DIR=/go/src/${PROJECT}

# Set work directory
RUN mkdir -p $PROJECT_DIR
WORKDIR $PROJECT_DIR

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...


# Server
EXPOSE 9000
STOPSIGNAL SIGINT

RUN go build -o main . 
CMD ["./main"]


