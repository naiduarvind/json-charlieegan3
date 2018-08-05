FROM golang:1.10 as build

WORKDIR /go/src/github.com/charlieegan3/json-charlieegan3

RUN go get -u github.com/golang/dep/cmd/dep

COPY . .

RUN CGO_ENABLED=0 go build -o statusUpdater cmd/run.go

# run container
FROM scratch
COPY --from=build /go/src/github.com/charlieegan3/json-charlieegan3/statusUpdater /
CMD ["/statusUpdater"]
