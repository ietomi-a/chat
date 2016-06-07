export GOPATH=`pwd`
rm chat
go build -o chat
./chat -addr=":8080"
