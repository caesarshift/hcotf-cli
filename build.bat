mkdir .\release\windows
mkdir .\release\linux
mkdir .\release\darwin

env GOOS=windows GOARCH=amd64 go build -o .\release\windows github.com/caesarshift/hcotf-cli
env GOOS=linux GOARCH=amd64 go build -o .\release\linux github.com/caesarshift/hcotf-cli
env GOOS=darwin GOARCH=amd64 go build -o .\release\darwin github.com/caesarshift/hcotf-cli
