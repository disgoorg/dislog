module github.com/DisgoOrg/dislog/example

go 1.16

replace (
	github.com/DisgoOrg/dislog => ../
)

require (
	github.com/DisgoOrg/dislog v1.0.4
	github.com/sirupsen/logrus v1.8.1
)
