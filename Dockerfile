FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

ENV GOOGLE_APPLICATION_CREDENTIALS /go/src/app/river-direction-210022-245643451dee.json

CMD ["app"]
