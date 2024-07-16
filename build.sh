basedir=`pwd`
cd $basedir

export NAME="gotools"

#x64
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o $NAME-$GOOS-$GOARCH

#aarch64
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm64
go build -o $NAME-$GOOS-$GOARCH