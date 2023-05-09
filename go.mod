module github.com/traefik/traefik/v2

go 1.16

// github.com/docker/docker v17.12.0-ce-rc1.0.20200204220554-5f6d6f3f2203+incompatible => v19.03.6
require (
	github.com/Azure/azure-sdk-for-go v40.3.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/BurntSushi/toml v0.3.1
	github.com/ExpediaDotCom/haystack-client-go v0.0.0-20190315171017-e7edbdf53a61
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/Microsoft/hcsshim v0.8.7 // indirect
	github.com/NYTimes/gziphandler v1.0.1 // indirect
	github.com/Shopify/sarama v1.23.1 // indirect
	github.com/abbot/go-http-auth v0.0.0-00010101000000-000000000000
	github.com/abronan/valkeyrie v0.0.0-20200127174252-ef4277a138cd
	github.com/armon/go-metrics v0.3.8 // indirect
	github.com/aws/aws-sdk-go v1.37.27
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/containerd/containerd v1.3.2 // indirect
	github.com/containous/alice v0.0.0-20181107144136-d83ebdd94cbd
	github.com/coredns/coredns v1.1.2 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/davecgh/go-spew v1.1.1
	github.com/denverdino/aliyungo v0.0.0-20170926055100-d3308649c661 // indirect
	github.com/digitalocean/godo v1.10.0 // indirect
	github.com/docker/cli v0.0.0-20200221155518-740919cc7fc0
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v17.12.0-ce-rc1.0.20200204220554-5f6d6f3f2203+incompatible
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-metrics v0.0.0-20181218153428-b84716841b82 // indirect
	github.com/docker/libcompose v0.0.0-20190805081528-eac9fe1b8b03 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/donovanhide/eventsource v0.0.0-20170630084216-b8f31a59085e // indirect
	github.com/eapache/channels v1.1.0
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/envoyproxy/go-control-plane v0.9.5 // indirect
	github.com/fatih/structs v1.1.0
	github.com/frankban/quicktest v1.11.0 // indirect
	github.com/gambol99/go-marathon v0.0.0-20180614232016-99a156b96fb2
	github.com/go-acme/lego/v4 v4.4.0
	github.com/go-check/check v0.0.0-00010101000000-000000000000
	github.com/go-kit/kit v0.10.1-0.20200915143503-439c4d2ed3ea
	github.com/go-resty/resty/v2 v2.1.1-0.20191201195748-d7b97669fe48
	github.com/golang/protobuf v1.5.2
	github.com/google/go-github/v28 v28.1.1
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20210601050228-01bbb1931b22 // indirect
	github.com/google/tcpproxy v0.0.0-20180808230851-dfa16c61dad2 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/consul v1.4.5
	github.com/hashicorp/consul/api v1.9.1
	github.com/hashicorp/go-bexpr v0.1.2 // indirect
	github.com/hashicorp/go-checkpoint v0.5.0 // indirect
	github.com/hashicorp/go-connlimit v0.3.0 // indirect
	github.com/hashicorp/go-hclog v0.16.1
	github.com/hashicorp/go-memdb v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-raftchunking v0.6.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.7 // indirect
	github.com/hashicorp/go-version v1.2.1
	github.com/hashicorp/hil v0.0.0-20200423225030-a18a1cd20038 // indirect
	github.com/hashicorp/mdns v1.0.4 // indirect
	github.com/hashicorp/memberlist v0.2.4 // indirect
	github.com/hashicorp/net-rpc-msgpackrpc v0.0.0-20151116020338-a14192a58a69 // indirect
	github.com/hashicorp/raft v1.3.1 // indirect
	github.com/hashicorp/raft-autopilot v0.1.5 // indirect
	github.com/hashicorp/vault/api v1.0.5-0.20200717191844-f687267c8086 // indirect
	github.com/hashicorp/vic v1.5.1-0.20190403131502-bbfe86ec9443 // indirect
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20191209144304-8bf82d3c094d
	github.com/instana/go-sensor v1.5.1
	github.com/joyent/triton-go v1.7.1-0.20200416154420-6801d15b779f // indirect
	github.com/klauspost/compress v1.13.0
	github.com/libkermit/compose v0.0.0-20171122111507-c04e39c026ad
	github.com/libkermit/docker v0.0.0-20171122101128-e6674d32b807
	github.com/libkermit/docker-check v0.0.0-20171122104347-1113af38e591
	github.com/lucas-clemente/quic-go v0.23.0
	github.com/mailgun/ttlmap v0.0.0-20170619185759-c1c17f74874f
	github.com/miekg/dns v1.1.43
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/go-testing-interface v1.14.0 // indirect
	github.com/mitchellh/hashstructure v1.0.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/mitchellh/pointerstructure v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
	github.com/nicolai86/scaleway-sdk v1.10.2-0.20180628010248-798f60e20bb2 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v1.0.0-rc10 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/packethost/packngo v0.1.1-0.20180711074735-b9cb5096f54c // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pires/go-proxyproto v0.5.0
	github.com/pmezard/go-difflib v1.0.0
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/client_model v0.2.0
	github.com/rancher/go-rancher-metadata v0.0.0-20200311180630-7f4c936a06ac
	github.com/rboyer/safeio v0.2.1 // indirect
	github.com/renier/xmlrpc v0.0.0-20170708154548-ce4a1a486c03 // indirect
	github.com/shirou/gopsutil/v3 v3.20.10 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/softlayer/softlayer-go v0.0.0-20180806151055-260589d94c7d // indirect
	github.com/stretchr/testify v1.7.0
	github.com/stvp/go-udp-testing v0.0.0-20191102171040-06b61409b154
	github.com/tencentcloud/tencentcloud-sdk-go v1.0.175 // indirect
	github.com/tent/http-link-go v0.0.0-20130702225549-ac974c61c2f9 // indirect
	github.com/tinylib/msgp v1.0.2 // indirect
	github.com/traefik/paerser v0.1.4
	github.com/traefik/yaegi v0.10.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	github.com/unrolled/render v1.0.2
	github.com/unrolled/secure v1.0.9
	github.com/vdemeester/shakers v0.1.0
	github.com/vmware/govmomi v0.18.0 // indirect
	github.com/vulcand/oxy v1.3.0
	github.com/vulcand/predicate v1.1.0
	go.elastic.co/apm v1.13.1
	go.elastic.co/apm/module/apmot v1.13.1
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/mod v0.4.2
	golang.org/x/net v0.0.0-20220114011407-0dd24b26b47d
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	golang.org/x/tools v0.1.1
	google.golang.org/grpc v1.43.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.19.0
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.21.0
	k8s.io/apiextensions-apiserver v0.20.2
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20210709001253-0e1f9d693477
	mvdan.cc/xurls/v2 v2.1.0
	sigs.k8s.io/gateway-api v0.3.0
	sigs.k8s.io/structured-merge-diff v0.0.0-20190525122527-15d366b2352e // indirect
)

// Containous forks
replace (
	github.com/abbot/go-http-auth => github.com/containous/go-http-auth v0.4.1-0.20200324110947-a37a7636d23e
	github.com/go-check/check => github.com/containous/check v0.0.0-20170915194414-ca0bf163426a
	github.com/gorilla/mux => github.com/containous/mux v0.0.0-20181024131434-c33f32e26898
	github.com/mailgun/minheap => github.com/containous/minheap v0.0.0-20190809180810-6e71eb837595
	github.com/mailgun/multibuf => github.com/containous/multibuf v0.0.0-20190809014333-8b6c9a7e6bba
)

replace google.golang.org/grpc => google.golang.org/grpc v1.27.1

replace github.com/tencentcloud/tencentcloud-sdk-go => github.com/tencentcloud/tencentcloud-sdk-go v1.0.175
