package relationtask

import "github.com/zxq97/gotool/config"

type RelationTaskConfig struct {
	Bind     string `yaml:"bind"`
	HTTPBind string `yaml:"http_bind"`
	Name     string `yaml:"name"`

	IncludeMySQL string `yaml:"include_mysql"`
	IncludeMC    string `yaml:"include_mc"`
	IncludeRedis string `yaml:"include_redis"`
	IncludeETCD  string `yaml:"include_etcd"`
	IncludeKafka string `yaml:"include_kafka"`

	Mysql map[string]*config.MysqlConf
	MC    map[string]*config.MCConf
	Redis map[string]*config.RedisConf
	Etcd  map[string]*config.EtcdConf
	Kafka map[string]*config.KafkaConf

	LogPath *config.LogConf `yaml:"log_path"`
}

func (conf *RelationTaskConfig) Initialize() {
	if conf.IncludeMySQL != "" {
		if err := config.LoadYaml(conf.IncludeMySQL, &conf.Mysql); err != nil {
			panic(err)
		}
	}
	if conf.IncludeMC != "" {
		if err := config.LoadYaml(conf.IncludeMC, &conf.MC); err != nil {
			panic(err)
		}
	}
	if conf.IncludeRedis != "" {
		if err := config.LoadYaml(conf.IncludeRedis, &conf.Redis); err != nil {
			panic(err)
		}
	}
	if conf.IncludeETCD != "" {
		if err := config.LoadYaml(conf.IncludeETCD, &conf.Etcd); err != nil {
			panic(err)
		}
	}
	if conf.IncludeKafka != "" {
		if err := config.LoadYaml(conf.IncludeKafka, &conf.Kafka); err != nil {
			panic(err)
		}
	}
}
