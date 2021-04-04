module github.com/finove/golibused

go 1.16

require (
	github.com/aws/aws-sdk-go v1.38.12
	github.com/beego/beego/v2 v2.0.1
	github.com/garyburd/redigo v1.6.2
	github.com/gocql/gocql v0.0.0-20210401103645-80ab1e13e309
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
)

replace github.com/beego/beego/v2 => github.com/finove/beego/v2 v2.0.2
