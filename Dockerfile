# build stage
FROM golang:1.20-alpine AS build
WORKDIR /build
COPY . .
RUN go get -d -v ./...
RUN go build -o app -v

# final stage
FROM alpine:latest
ARG SOURCE_BRANCH
ARG SOURCE_COMMIT
ARG IMAGE_NAME
ENV SOURCE_BRANCH $SOURCE_BRANCH
ENV SOURCE_COMMIT $SOURCE_COMMIT
ENV IMAGE_NAME $IMAGE_NAME
WORKDIR /dist
COPY --from=build /build .
EXPOSE 3000
CMD ["/dist/app", "app"]
