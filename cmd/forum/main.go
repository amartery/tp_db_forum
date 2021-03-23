package main

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/amartery/tp_db_forum/internal/app/http_server"
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

	// postgresCon, err := utility.CreatePostgresConnection(config.DataBaseURL)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// statRep := postgresDB.NewStatRepository(postgresCon)
	// statUsecase := usecase.NewStatUsecase(statRep)

	s := http_server.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
