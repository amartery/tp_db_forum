package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/amartery/tp_db_forum/internal/app/http_server"
	"github.com/jackc/pgx/v4/pgxpool"

	UsecaseForum "github.com/amartery/tp_db_forum/internal/app/forum/usecase"
	UsecasePost "github.com/amartery/tp_db_forum/internal/app/post/usecase"
	UsecaseService "github.com/amartery/tp_db_forum/internal/app/service/usecase"
	UsecaseThread "github.com/amartery/tp_db_forum/internal/app/thread/usecase"
	UsecaseUser "github.com/amartery/tp_db_forum/internal/app/user/usecase"

	RepositoryForum "github.com/amartery/tp_db_forum/internal/app/forum/repository"
	RepositoryPost "github.com/amartery/tp_db_forum/internal/app/post/repository"
	RepositoryService "github.com/amartery/tp_db_forum/internal/app/service/repository"
	RepositoryThread "github.com/amartery/tp_db_forum/internal/app/thread/repository"
	RepositoryUser "github.com/amartery/tp_db_forum/internal/app/user/repository"
)

var (
	configPath string = "./configs/forum.toml"
)

func main() {
	config := http_server.NewConfig()

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	// "host=localhost dbname=db_tp user=admin password=admin sslmode=disable"
	DBcon, err := pgxpool.Connect(
		context.Background(),
		"host=localhost user=docker password=docker dbname=forum sslmode=disable",
	)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer DBcon.Close()

	forumRepo := RepositoryForum.NewForumRepository(DBcon)
	postRepo := RepositoryPost.NewPostRepository(DBcon)
	serviceRepo := RepositoryService.NewServiceRepository(DBcon)
	threadRepo := RepositoryThread.NewThreadRepository(DBcon)
	userRepo := RepositoryUser.NewUserRepository(DBcon)

	forumUse := UsecaseForum.NewForumUsecase(forumRepo)
	postUse := UsecasePost.NewPostUsecase(postRepo, userRepo, forumRepo, threadRepo)
	serviceUse := UsecaseService.NewServiceUsecase(serviceRepo)
	threadUse := UsecaseThread.NewThreadUsecase(threadRepo, userRepo)
	userUse := UsecaseUser.NewUserUsecase(userRepo)

	s := http_server.New(config, forumUse, postUse, serviceUse, threadUse, userUse)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
