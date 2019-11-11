package protobuf

//go:generate retool do protoc -I ../vendor -I . --twirp_out=paths=source_relative:. --go_out=paths=source_relative:. --twirp_typescript_out=. drawer/feather.proto drawer/glitter.proto
//go:generate retool do protoc -I ../vendor -I . --twirp_out=paths=source_relative:. --go_out=paths=source_relative:. --twirp_typescript_out=. example/service.proto
