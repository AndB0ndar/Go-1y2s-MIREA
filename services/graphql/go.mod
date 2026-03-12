module app/services/graphql

go 1.25.7

require (
	app/shared v0.0.0
	github.com/99designs/gqlgen v0.17.49
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/sirupsen/logrus v1.9.4
	github.com/vektah/gqlparser/v2 v2.5.16
	google.golang.org/grpc v1.79.1
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace app/shared => ../../shared
