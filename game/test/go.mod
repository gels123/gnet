module test

go 1.19

require (
	gnet v1.0.0 //indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
)

replace (
	gnet => ../../
	go.uber.org/zap => ../../lib/logzap/zap
)
