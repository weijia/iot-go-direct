module iot_go

go 1.22.2

require (
	// github.com/carlmjohnson/versioninfo v0.22.5 // not used anymore
	github.com/coreos/go-systemd v0.0.0-20230601205942-d843340ab4bd
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/kraken-hpc/go-fork v0.1.1
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/pkg/sftp v1.13.6
	github.com/spf13/cobra v1.8.1
	github.com/spf13/viper v1.16.0
	github.com/weijia/supervisord/supervisord_main v0.0.0-20240704061853-d431807ae74a
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/crypto v0.18.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/rpc v1.2.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/go-envparse v0.1.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jessevdk/go-flags v1.6.1 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kardianos/service v1.2.2 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/ochinchina/filechangemonitor v0.3.1 // indirect
	github.com/ochinchina/go-daemon v0.1.5 // indirect
	github.com/ochinchina/go-ini v1.0.1 // indirect
	github.com/ochinchina/go-reaper v0.0.0-20181016012355-6b11389e79fc // indirect
	github.com/ochinchina/gorilla-xmlrpc v0.0.0-20171012055324-ecf2fe693a2c // indirect
	github.com/ochinchina/supervisord/config v0.0.0-20230902082938-c2cae38b7454 // indirect
	github.com/ochinchina/supervisord/events v0.0.0-20230902082938-c2cae38b7454 // indirect
	github.com/ochinchina/supervisord/faults v0.0.0-20230902082938-c2cae38b7454 // indirect
	github.com/ochinchina/supervisord/logger v0.0.0-20230902082938-c2cae38b7454 // indirect
	github.com/ochinchina/supervisord/process v0.0.0-20230902082938-c2cae38b7454 // indirect
	github.com/ochinchina/supervisord/signals v0.0.0-20230902082938-c2cae38b7454 // indirect
	github.com/ochinchina/supervisord/types v0.0.0-20230902082938-c2cae38b7454 // indirect
	github.com/ochinchina/supervisord/util v0.0.0-20230902082938-c2cae38b7454 // indirect
	github.com/ochinchina/supervisord/xmlrpcclient v0.0.0-20230902082938-c2cae38b7454 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-charset v0.0.0-20190617161244-0dc95cdf6f31 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.1.0

replace github.com/ochinchina/supervisord/process => github.com/weijia/supervisord/process v0.0.0-20240709105315-49e6c872ebd4
