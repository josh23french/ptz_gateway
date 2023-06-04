module github.com/josh23french/ptz_gateway

go 1.14

require (
	github.com/josh23french/visca v0.1.1
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.6.1
	go.bug.st/serial v1.1.2
	golang.org/x/sys v0.0.0-20210105210732-16f7687f5001 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210105161348-2e78108cf5f8
)

replace github.com/josh23french/visca => ../visca
