//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
	"github.com/zxq97/relation/app/relationship/job/internal/data"
	"github.com/zxq97/relation/app/relationship/job/internal/service"
)

func initAPP() {
	wire.Build(service.ProviderSet, biz.ProviderSet, data.ProviderSet)
}
