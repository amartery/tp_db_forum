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

	"github.com/amartery/tp_db_forum/internal/app/middleware"
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
	s.router.POST("/api/forum/create", middleware.ContentTypeJson(forumHandlers.ForumCreate))
	s.router.GET("/api/forum/{slug}/details", middleware.ContentTypeJson(forumHandlers.ForumDetails))
	s.router.POST("/api/forum/{slug}/create", middleware.ContentTypeJson(forumHandlers.ForumCreateBranch))
	s.router.GET("/api/forum/{slug}/users", middleware.ContentTypeJson(forumHandlers.CurrentForumUsers))
	s.router.GET("/api/forum/{slug}/threads", middleware.ContentTypeJson(forumHandlers.ForumBranches))

	postHandlers := DeliveryPost.NewPostHandler(s.usecasePost, s.usecaseUser, s.usecaseForum, s.usecaseThread)
	s.router.GET("/api/post/{id}/details", middleware.ContentTypeJson(postHandlers.PostDetailsGet))
	s.router.POST("/api/post/{id}/details", middleware.ContentTypeJson(postHandlers.PostDetailsUpdate))

	serviceHandlers := DeliveryService.NewServiceHandler(s.usecaseService)
	s.router.POST("/api/service/clear", middleware.ContentTypeJson(serviceHandlers.ServiceClear))
	s.router.GET("/api/service/status", middleware.ContentTypeJson(serviceHandlers.ServiceStatus))

	threadHandlers := DeliveryThread.NewThreadHandler(s.usecaseThread, s.usecaseUser)
	s.router.POST("/api/thread/{slug_or_id}/create", middleware.ContentTypeJson(threadHandlers.CreatePostInBranch))
	s.router.GET("/api/thread/{slug_or_id}/details", middleware.ContentTypeJson(threadHandlers.BranchDetailsGet))
	s.router.POST("/api/thread/{slug_or_id}/details", middleware.ContentTypeJson(threadHandlers.BranchDetailsUpdate))
	s.router.POST("/api/thread/{slug_or_id}/vote", middleware.ContentTypeJson(threadHandlers.VoteForBranch))
	s.router.GET("/api/thread/{slug_or_id}/posts", middleware.ContentTypeJson(threadHandlers.CurrentBranchPosts))

	userHandlers := DeliveryUser.NewUserHandler(s.usecaseUser, s.usecaseForum)
	s.router.POST("/api/user/{nickname}/create", middleware.ContentTypeJson(userHandlers.CreateUser))
	s.router.GET("/api/user/{nickname}/profile", middleware.ContentTypeJson(userHandlers.AboutUserGet))
	s.router.POST("/api/user/{nickname}/profile", middleware.ContentTypeJson(userHandlers.AboutUserUpdate))
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
