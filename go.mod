module marsgo

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/PuerkitoBio/goquery v1.5.1 // indirect
	github.com/Shopify/sarama v1.23.1
	github.com/antchfx/htmlquery v1.2.2 // indirect
	github.com/antchfx/xmlquery v1.2.3 // indirect
	github.com/antchfx/xpath v1.1.4 // indirect
	github.com/bilibili/kratos v0.3.3
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/coreos/bbolt v1.3.4 // indirect
	github.com/coreos/etcd v3.3.20+incompatible
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-redis/redis/v7 v7.2.0
	github.com/go-resty/resty/v2 v2.2.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-zookeeper/zk v1.0.1
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gocolly/colly v1.2.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/google/wire v0.4.0
	github.com/json-iterator/go v1.1.9
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/marsofsnow/zkgroup v1.0.0 // indirect
	github.com/montanaflynn/stats v0.6.3
	github.com/nanu-c/zkgroup v0.8.7 // indirect
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/philchia/agollo v2.1.0+incompatible
	github.com/philchia/agollo/v3 v3.1.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.4.1
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/shirou/gopsutil v2.20.1+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.6.1
	github.com/temoto/robotstxt v1.1.1 // indirect
	github.com/tidwall/gjson v1.5.0
	github.com/tsuna/gohbase v0.0.0-20190823190353-a66bcc9075db
	go.etcd.io/etcd v3.3.20+incompatible
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200731012542-8145dea6a485
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v2 v2.2.8

)

replace (
	google.golang.org/grpc v1.27.1 => google.golang.org/grpc v1.26.0
	google.golang.org/grpc v1.29.1 => google.golang.org/grpc v1.26.0
)

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
