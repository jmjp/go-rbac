#Start by building the application.
FROM golang:1.22.2-bullseye as build

WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app -ldflags "-s -w" cmd/main.go


# Now copy it into our base image.
FROM gcr.io/distroless/static-debian12
# FROM scratch
COPY --from=build /go/bin/app /

CMD ["/app"]
