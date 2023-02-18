package main

import (
	"flag"

	"github.com/zxq97/gokit/pkg/config"
	"github.com/zxq97/gokit/pkg/database/xmysql"
	"github.com/zxq97/relation/app/relationship/pkg/dal/method"
	"gorm.io/gen"
)

var (
	flagConf string
	genConf  generateConfig
)

type generateConfig struct {
	Mysql *xmysql.Config `yaml:"mysql"`
}

func init() {
	flag.StringVar(&flagConf, "conf", "app/relationship/generate/generate.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	err := config.LoadYaml(flagConf, &genConf)
	if err != nil {
		panic(err)
	}
	db, err := xmysql.NewMysqlDB(genConf.Mysql)
	if err != nil {
		panic(err)
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:      "app/relationship/pkg/dal/query",
		ModelPkgPath: "app/relationship/pkg/dal/model",
	})
	g.UseDB(db)

	g.ApplyInterface(func(method.UserFollow) {}, g.GenerateModel("user_follows"))
	g.ApplyInterface(func(method.UserFollower) {}, g.GenerateModel("user_followers"))
	g.ApplyInterface(func(method.UserRelationCount) {}, g.GenerateModel("user_relation_counts"))
	g.ApplyInterface(func(method.ExtraFollower) {}, g.GenerateModel("extra_followers"))

	g.Execute()
}
