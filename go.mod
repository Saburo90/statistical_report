module gitee.com/NotOnlyBooks/statistical_report

require (
	cloud.google.com/go v0.33.1 // indirect
	github.com/360EntSecGroup-Skylar/excelize v1.4.0 // indirect
	github.com/BurntSushi/toml v0.3.1
	github.com/Chronokeeper/anyxml v0.0.0-20160530174208-54457d8e98c6 // indirect
	github.com/CloudyKit/fastprinter v0.0.0-20170127035650-74b38d55f37a // indirect
	github.com/CloudyKit/jet v2.1.2+incompatible // indirect
	github.com/Luxurioust/excelize v1.4.0
	github.com/agrison/go-tablib v0.0.0-20160310143025-4930582c22ee // indirect
	github.com/agrison/mxj v0.0.0-20160310142625-1269f8afb3b4 // indirect
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf // indirect
	github.com/bndr/gotabulate v1.1.2 // indirect
	github.com/clbanning/mxj v1.8.3 // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20181014144952-4e0d7dc8888f // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51 // indirect
	github.com/facebookgo/freeport v0.0.0-20150612182905-d4adf43b75b9 // indirect
	github.com/facebookgo/grace v0.0.0-20180706040059-75cf19382434
	github.com/facebookgo/httpdown v0.0.0-20180706035922-5979d39b15c2 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/stats v0.0.0-20151006221625-1b76add642e4 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/gin-contrib/cors v0.0.0-20181008113111-488de3ec974f
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/go-openapi/jsonreference v0.17.2 // indirect
	github.com/go-openapi/spec v0.17.2 // indirect
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible
	github.com/go-redis/redis v6.14.2+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/builder v0.3.3 // indirect
	github.com/go-xorm/core v0.6.0 // indirect
	github.com/go-xorm/sqlfiddle v0.0.0-20180821085327-62ce714f951a // indirect
	github.com/golang/protobuf v1.2.0
	github.com/google/go-cmp v0.2.0 // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgx v3.2.0+incompatible // indirect
	github.com/json-iterator/go v1.1.5 // indirect
	github.com/lib/pq v1.0.0 // indirect
	github.com/mattn/go-isatty v0.0.4
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/robfig/cron v0.0.0-20180505203441-b41be1df6967
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/swaggo/gin-swagger v1.0.0
	github.com/swaggo/swag v1.4.0
	github.com/tealeg/xlsx v1.0.3 // indirect
	github.com/ugorji/go/codec v0.0.0-20181127175209-856da096dbdf // indirect
	github.com/xormplus/core v0.0.0-20181016121923-6bfce2eb8867
	github.com/xormplus/xorm v0.0.0-20181105071520-4fd8a981b629
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/crypto v0.0.0-20181127143415-eb0de9b17e85 // indirect
	golang.org/x/net v0.0.0-20181129055619-fae4c4e3ad76
	golang.org/x/sys v0.0.0-20181128092732-4ed8d59d0b35 // indirect
	golang.org/x/tools v0.0.0-20181128225727-c5b00d9557fd // indirect
	google.golang.org/appengine v1.3.0 // indirect
	google.golang.org/grpc v1.16.0
	gopkg.in/flosch/pongo2.v3 v3.0.0-20141028000813-5e81b817a0c4 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
)

replace (
	cloud.google.com/go => github.com/GoogleCloudPlatform/google-cloud-go v0.27.0
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20180910181607-0e37d006457b
	golang.org/x/lint/golint => github.com/golang/lint v0.0.0-20180702182130-06c8688daad7
	golang.org/x/net => github.com/golang/net v0.0.0-20180911220305-26e67e76b6c3
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/sync => github.com/golang/sync v0.0.0-20180314180146-1d60e4601c6f
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180909124046-d0be0721c37e
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20180820150726-fbb02b2291d28
	golang.org/x/tools => github.com/golang/tools v0.0.0-20181128225727-c5b00d9557fd
	google.golang.org/appengine => github.com/golang/appengine v1.1.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20180817151627-c66870c02cf8
	google.golang.org/grpc => github.com/grpc/grpc-go v1.16.0
	protos v0.0.0 => ../protos
)
