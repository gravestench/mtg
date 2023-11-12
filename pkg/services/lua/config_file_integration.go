package lua

import (
	"github.com/gravestench/mtg/pkg/services/configFile"
)

var _ configFile.HasDefaultConfig = &Service{}

func (s *Service) ConfigFileName() string {
	return "lua_environment.json"
}

func (s *Service) DefaultConfig() (cfg configFile.Config) {
	g := cfg.Group(s.Name())

	g.SetDefault("init script", "init.lua")

	return
}
