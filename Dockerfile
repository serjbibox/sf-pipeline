FROM golang:latest AS build_stage
RUN mkdir -p go/src/sf-pipeline
WORKDIR /go/src/sf-pipeline
COPY ./ ./
RUN go env -w GO111MODULE=auto 
RUN go install .

FROM alpine:latest
WORKDIR /
COPY --from=build_stage /go/bin .
#RUN apk add libc6-compat
ENTRYPOINT ./SF20.2.1
