# drone-artifactory-cache

Drone plugin for caching artifacts to a central artifactory server. For the
usage information and a listing of the available options please take a look at
[the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t jmccann/drone-artifactory-cache .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-artifactory-cache' not found or does not exist..
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e DRONE_REPO=octocat/hello-world \
  -e DRONE_REPO_BRANCH=master \
  -e DRONE_COMMIT_BRANCH=master \
  -e PLUGIN_MOUNT=node_modules \
  -e PLUGIN_RESTORE=false \
  -e PLUGIN_REBUILD=true \
  -e PLUGIN_PATH=repo-key/project/archive.tar \
  -e CACHE_ARTIFACTORY_URL=https://company.com/artifactory \
  -e CACHE_ARTIFACTORY_USERNAME=johndoe \
  -e CACHE_ARTIFACTORY_PASSWORD=supersecret \
  jmccann/drone-artifactory-cache:latest
```
