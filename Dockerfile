FROM golang:1.19-alpine3.16 as build

ENV CGO_ENABLED=0 \
    GOOS=linux

RUN apk add --no-cache curl git

# Download and install go-swagger tool and swagger-ui.
RUN curl -L https://github.com/go-swagger/go-swagger/releases/download/v0.27.0/swagger_linux_amd64 -o /go/bin/swagger && \
    chmod a+x /go/bin/swagger && \
    curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v3.47.1.tar.gz | \
    tar -xz --strip-components 1 -C /go/bin swagger-ui-3.47.1/dist && \
    mv /go/bin/dist /go/bin/static && \
    sed -i 's,https://petstore.swagger.io/v2/swagger.json,swagger.yml,g' /go/bin/static/index.html

# Download generating tool.
RUN go install github.com/vektra/mockery/v2@v2.14.0

WORKDIR /go/src/bitbucket.org/creativeadvtech/project-template

# Download all project dependencies.
# This trick exploits how the docker uses the layers created on RUN statements.
# The dependencies will be re-downloaded only if the `go.mod` or `go.sum` was changed.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code, perform testing, swagger file and build binaries.
COPY . .
RUN go generate ./... && \
    go test -short ./... && \
    go install ./cmd/app

RUN mv ./api/swagger.yml /go/bin/static/swagger.yml

FROM alpine:3.16

RUN apk add --no-cache curl bash

# Copy swagger, migrations, application binary, and runnable bash script.
WORKDIR /
COPY --from=build /go/src/bitbucket.org/creativeadvtech/project-template/scripts/wait-for-it.sh ./wait-for-it.sh
COPY --from=build /go/src/bitbucket.org/creativeadvtech/project-template/pkg/database/*.yml ./
COPY --from=build /go/src/bitbucket.org/creativeadvtech/project-template/migrations ./migrations
COPY --from=build /go/bin/app ./bin/
COPY --from=build /go/bin/static ./static
RUN chmod +x ./wait-for-it.sh

CMD app