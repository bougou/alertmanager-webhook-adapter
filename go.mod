module github.com/bougou/alertmanager-webhook-adapter

go 1.16

replace github.com/bougou/alertmanager-webhook-adapter v0.0.0 => ./

require (
	github.com/bougou/webhook-adapter v1.0.8
	github.com/emicklei/go-restful/v3 v3.4.0
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/spf13/cobra v1.1.3
)
