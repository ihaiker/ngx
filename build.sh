#!/bin/bash
set -e

Version=$(git describe --tags $(git rev-list --tags --max-count=1))
GitCommit=$(git rev-parse HEAD)
BuildDate=$(date +"%F %T")

debug="-w -s"
param="-X main.VERSION=${Version} -X main.GITLOG_VERSION=${GitCommit} -X 'main.BUILD_TIME=${BuildDate}'"

build() {
  echo "build $1 $2 $3"
  export CGO_ENABLED=0
  export GOOS=$1
  export GOARCH=$2
  export SUFFIX=$3
  go build -ldflags "${debug} ${param}" -o bin/ngx-${GOOS}-${GOARCH}${SUFFIX} cmd/main.go

  if [ "$GOOS" == "windows" ]; then
    zip bin/dist/ngx-${GOOS}-${GOARCH}.zip bin/ngx-${GOOS}-${GOARCH}${SUFFIX}
  else
    tar -czvf bin/dist/ngx-${GOOS}-${GOARCH}.tar.gz bin/ngx-${GOOS}-${GOARCH}${SUFFIX}
  fi
}
mkdir -p bin/dist
build windows amd64 .exe
build windows 386 .exe
build windows arm .exe
build darwin amd64
build linux amd64
build linux 386
build linux arm
build freebsd amd64
build freebsd 386
build freebsd arm
