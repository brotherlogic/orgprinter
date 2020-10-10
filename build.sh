protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/orgprinter.proto
mv proto/github.com/brotherlogic/orgprinter/proto/* ./proto
