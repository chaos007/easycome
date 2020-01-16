module github.com/chaos007/easycome

go 1.12

replace cloud.google.com/go v0.37.4 => github.com/googleapis/google-cloud-go v0.37.4

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/go-xorm/xorm v0.7.6
	github.com/golang/protobuf v1.3.2
	github.com/prometheus/common v0.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spaolacci/murmur3 v1.1.0
	github.com/xtaci/kcp-go v5.4.4+incompatible
	gitlab.dianchu.cc/2020008/test-lol v0.0.0-20200115112652-566ca6fea7a8
	go.etcd.io/etcd v3.3.15+incompatible
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297
	google.golang.org/grpc v1.26.0
)
