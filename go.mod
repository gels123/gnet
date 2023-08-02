module gnet

go 1.19

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/lestrrat-go/strftime v1.0.6
)

replace (
	go.uber.org/zap => ./lib/logzap/zap
	github.com/pkg/errors => ./lib/errors
	github.com/lestrrat-go/strftime => ./lib/strftime
)