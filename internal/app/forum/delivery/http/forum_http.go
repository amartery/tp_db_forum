package http

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// ForumServer ...
type ForumServer struct {
	config *Config
	logger *logrus.Logger
	router *router.Router
}

// New ...
func New(config *Config) *ForumServer {
	return &ForumServer{
		config: config,
		logger: logrus.New(),
		router: router.New(),
	}
}

func (s *ForumServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()

	s.logger.Info("starting statistics server" + s.config.BindAddr)
	return fasthttp.ListenAndServe(s.config.BindAddr, s.router.Handler)

}

func (s *ForumServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *ForumServer) configureRouter() {
	s.router.POST("/forum/create", s.ForumCreate)

	// TODO: s.router.GET("/forum/{slug}/details", s.ForumDetails)
	// TODO: s.router.GET("/forum/{slug}/details", s.ForumDetails)
	// TODO: s.router.POST("/forum/{slug}/create", s.ForumCreateBranch)
	// TODO: s.router.GET("/forum/{slug}/users", s.CurrentForumUsers)
	// TODO: s.router.GET("​/forum​/{slug}​/threads", s.ForumBranches)
	// TODO: s.router.GET("​/post​/{id}​/details", s.PostDetailsGet)  // ??? в сваггере написано что для ветки а не для поста
	// TODO: s.router.POST("​/post​/{id}​/details", s.PostDetailsUpdate)
	// TODO: s.router.POST("​/service​/clear", s.ServiceClear)
	// TODO: s.router.GET("​/service​/status", s.ServiceStatus)
	// TODO: s.router.POST("​/thread​/{slug_or_id}​/create", s.CreatePostInBranch)
	// TODO: s.router.GET("/thread​/{slug_or_id}​/details", s.BranchDetailsGet)
	// TODO: s.router.POST("/thread​/{slug_or_id}​/details", s.BranchDetailsUpdate)
	// TODO: s.router.GET("/thread/{slug_or_id}/posts", s.CurrentBranchPosts)
	// TODO: s.router.POST("/thread​/{slug_or_id}​/vote", s.VoteForBranch)

}

func (f *ForumServer) ForumCreate(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumCreate")
	fmt.Fprintf(ctx, "hello world!\n")
}
