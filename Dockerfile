# Docker image for the Drone build runner
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t jmccann/drone-artifactory-cache .

FROM alpine:3.4

# Install required packages
RUN apk update && \
  apk add \
    ca-certificates && \
  rm -rf /var/cache/apk/*

ADD drone-artifactory-cache /bin/
ENTRYPOINT ["/bin/drone-artifactory-cache"]
