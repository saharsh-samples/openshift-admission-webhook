# build binary
FROM golang:1.15.3 AS build

WORKDIR /go/src/maw
ADD src .

RUN ./testall.sh

RUN GOOS=linux CGO_ENABLED=0 go build

# create runnable image
FROM scratch

COPY --from=build /go/src/maw/maw /maw

ENTRYPOINT [ "/maw" ]