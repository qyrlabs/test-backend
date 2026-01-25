module github.com/qyrlabs/test-backend/payment

go 1.25.6

replace github.com/qyrlabs/test-backend/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/qyrlabs/test-backend/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.78.0
)

require (
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120174246-409b4a993575 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
