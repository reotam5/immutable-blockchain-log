#/bin/bash!

# run backend api gateway
pushd ./log-client
go mod tidy
go build -o gateway ./cmd/gateway/main.go
popd

# run frontend
pushd ./log-dashboard
rm -rf node_modules
npm install
npm run build
npm run run-prod
popd
