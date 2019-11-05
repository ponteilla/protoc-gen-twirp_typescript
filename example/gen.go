package example

//go:generate retool do protoc -I ../vendor -I .--twirp_out=paths=source_relative:. --go_out=paths=source_relative:. --twirp_typescript_out=. service.proto
