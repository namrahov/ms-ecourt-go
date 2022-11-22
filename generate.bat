protoc -Igreet/proto --go_opt=module=github.com/namrahov/ms-ecourt-go --go_out=. --go-grpc_opt=module=github.com/namrahov/ms-ecourt-go --go-grpc_out=. greet/proto/*.proto
