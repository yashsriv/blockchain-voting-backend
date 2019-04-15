module blockchain-voting

go 1.12

require (
	ethlib v0.0.0-00010101000000-000000000000
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/jlaffaye/ftp v0.0.0-20190411155707-52d3001130a6
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/mediocregopher/radix/v3 v3.2.3
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/spf13/viper v1.3.2
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
)

replace ethlib => ../blockchain-voting-ethlib
