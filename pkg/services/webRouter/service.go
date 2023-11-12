package webRouter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/gravestench/runtime"

	"github.com/gravestench/mtg/pkg/services/configFile"
	"github.com/gravestench/mtg/pkg/services/webRouter/middleware/static_assets"
)

type Service struct {
	log        *zerolog.Logger
	cfgManager configFile.Dependency

	root *gin.Engine

	boundServices map[string]*struct{} // holds 0-size entries

	config struct {
		debug bool
	}

	reloadDebounce time.Time
}

func (s *Service) BindLogger(l *zerolog.Logger) {
	s.log = l
}

func (s *Service) Logger() *zerolog.Logger {
	return s.log
}

func (s *Service) Init(rt runtime.R) {
	gin.SetMode("release")
	rt.Add(&static_assets.Middleware{})
	s.root = gin.New()
	go s.beginDynamicRouteBinding(rt)
}

func (s *Service) Name() string {
	return "Web Router"
}
