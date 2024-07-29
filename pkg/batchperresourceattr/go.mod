module github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchperresourceattr

go 1.21.0

require (
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/collector/consumer v0.105.1-0.20240729082905-fb5b1e6aa550
	go.opentelemetry.io/collector/consumer/consumertest v0.0.0-20240726220702-a6287aca1a43
	go.opentelemetry.io/collector/pdata v1.12.1-0.20240729082905-fb5b1e6aa550
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.11.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/collector/consumer/consumerprofiles v0.105.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.105.1-0.20240729082905-fb5b1e6aa550 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.76.2
	v0.76.1
	v0.65.0
)
