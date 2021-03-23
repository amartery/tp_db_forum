package http_server

import (
	DeliveryForum "github.com/amartery/tp_db_forum/internal/app/forum/delivery/http"
	DeliveryPost "github.com/amartery/tp_db_forum/internal/app/post/delivery/http"
	DeliveryService "github.com/amartery/tp_db_forum/internal/app/service/delivery/http"
	DeliveryThread "github.com/amartery/tp_db_forum/internal/app/thread/delivery/http"
	DeliveryUser "github.com/amartery/tp_db_forum/internal/app/user/delivery/http"
	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

//  DeliveryForum "github.com/amartery/tp_db_forum/internal/app/forum/delivery/http"
// 	DeliveryPost "github.com/amartery/tp_db_forum/internal/app/post/delivery/http"
// 	DeliveryService "github.com/amartery/tp_db_forum/internal/app/service/delivery/http"
// 	DeliveryThread "github.com/amartery/tp_db_forum/internal/app/thread/delivery/http"
// 	DeliveryUser "github.com/amartery/tp_db_forum/internal/app/user/delivery/http"

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
	s.router.POST("/forum/create", DeliveryForum.ForumCreate)
	s.router.GET("/forum/{slug}/details", DeliveryForum.ForumDetails)
	s.router.POST("/forum/{slug}/create", DeliveryForum.ForumCreateBranch)
	s.router.GET("/forum/{slug}/users", DeliveryForum.CurrentForumUsers)
	s.router.GET("/forum/{slug}/threads", DeliveryForum.ForumBranches)

	// ??? в сваггере написано что для ветки а не для поста
	s.router.GET("/post/{id}/details", DeliveryPost.PostDetailsGet)
	s.router.POST("/post/{id}/details", DeliveryPost.PostDetailsUpdate)

	s.router.POST("/service/clear", DeliveryService.ServiceClear)
	s.router.GET("/service/status", DeliveryService.ServiceStatus)

	s.router.POST("/thread/{slug_or_id}/create", DeliveryThread.CreatePostInBranch)
	s.router.GET("/thread/{slug_or_id}/details", DeliveryThread.BranchDetailsGet)
	s.router.POST("/thread/{slug_or_id}/details", DeliveryThread.BranchDetailsUpdate)
	s.router.POST("/thread/{slug_or_id}/vote", DeliveryThread.VoteForBranch)
	s.router.GET("/thread/{slug_or_id}/posts", DeliveryThread.CurrentBranchPosts)

	s.router.POST("/user/{nickname}/create", DeliveryUser.CreateUser)
	s.router.GET("/user/{nickname}/profile", DeliveryUser.AboutUserGet)
	s.router.POST("/user/{nickname}/profile", DeliveryUser.AboutUserUpdate)
}
