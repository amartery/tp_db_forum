package http_server

import (
	"github.com/amartery/tp_db_forum/internal/app/forum"
	DeliveryForum "github.com/amartery/tp_db_forum/internal/app/forum/delivery/http"
	"github.com/amartery/tp_db_forum/internal/app/post"
	DeliveryPost "github.com/amartery/tp_db_forum/internal/app/post/delivery/http"
	"github.com/amartery/tp_db_forum/internal/app/service"
	DeliveryService "github.com/amartery/tp_db_forum/internal/app/service/delivery/http"
	"github.com/amartery/tp_db_forum/internal/app/thread"
	DeliveryThread "github.com/amartery/tp_db_forum/internal/app/thread/delivery/http"
	"github.com/amartery/tp_db_forum/internal/app/user"
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
	config         *Config
	logger         *logrus.Logger
	router         *router.Router
	usecaseForum   forum.Usecase
	usecasePost    post.Usecase
	usecaseService service.Usecase
	usecaseThread  thread.Usecase
	usecaseUser    user.Usecase
}

// New ...
func New(
	config *Config,
	forumUse forum.Usecase,
	postUse post.Usecase,
	serviceUse service.Usecase,
	threadUse thread.Usecase,
	userUse user.Usecase) *ForumServer {
	return &ForumServer{
		config:         config,
		logger:         logrus.New(),
		router:         router.New(),
		usecaseForum:   forumUse,
		usecasePost:    postUse,
		usecaseService: serviceUse,
		usecaseThread:  threadUse,
		usecaseUser:    userUse,
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

	forumHandlers := DeliveryForum.NewForumHandler(s.usecaseForum, s.usecaseUser, s.usecaseThread)
	s.router.POST("/forum/create", forumHandlers.ForumCreate)
	s.router.GET("/forum/{slug}/details", forumHandlers.ForumDetails)
	s.router.POST("/forum/{slug}/create", forumHandlers.ForumCreateBranch)
	s.router.GET("/forum/{slug}/users", forumHandlers.CurrentForumUsers)
	s.router.GET("/forum/{slug}/threads", forumHandlers.ForumBranches)

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

// DeliveryForum "github.com/amartery/tp_db_forum/internal/app/forum/delivery/http"
// DeliveryPost "github.com/amartery/tp_db_forum/internal/app/post/delivery/http"
// DeliveryService "github.com/amartery/tp_db_forum/internal/app/service/delivery/http"
// DeliveryThread "github.com/amartery/tp_db_forum/internal/app/thread/delivery/http"
// DeliveryUser "github.com/amartery/tp_db_forum/internal/app/user/delivery/http"

// UsecaseForum "github.com/amartery/tp_db_forum/internal/app/forum/usecase"
// UsecasePost "github.com/amartery/tp_db_forum/internal/app/post/usecase"
// UsecaseService "github.com/amartery/tp_db_forum/internal/app/service/usecase"
// UsecaseThread "github.com/amartery/tp_db_forum/internal/app/thread/usecase"
// UsecaseUser "github.com/amartery/tp_db_forum/internal/app/user/usecase"

// RepositoryForum "github.com/amartery/tp_db_forum/internal/app/forum/repository"
// RepositoryPost "github.com/amartery/tp_db_forum/internal/app/post/repository"
// RepositoryService "github.com/amartery/tp_db_forum/internal/app/service/repository"
// RepositoryThread "github.com/amartery/tp_db_forum/internal/app/thread/repository"
// RepositoryUser "github.com/amartery/tp_db_forum/internal/app/user/repository"
