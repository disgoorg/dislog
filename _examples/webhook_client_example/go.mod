module github.com/DisgoOrg/dislog/_examples/webhook_client_example

go 1.17

replace github.com/DisgoOrg/dislog => ../../

require (
	github.com/DisgoOrg/disgo v0.6.8
	github.com/DisgoOrg/dislog v1.0.2
	github.com/sirupsen/logrus v1.8.1
)

require (
	github.com/DisgoOrg/log v1.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sasha-s/go-csync v0.0.0-20210812194225-61421b77c44b // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
)
