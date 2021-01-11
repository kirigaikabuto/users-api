package main

import (
	"fmt"
	"github.com/djumanoff/amqp"
	users_store "github.com/kirigaikabuto/users-store"
	"log"
)

var cfg = amqp.Config{
	Host:        "localhost",
	VirtualHost: "",
	User:        "",
	Password:    "",
	Port:        5672,
	LogLevel:    5,
}

var srvCfg = amqp.ServerConfig{
	ResponseX: "response",
	RequestX:  "request",
}

var cfgAmqp = amqp.Config{
	Host:        "localhost",
	VirtualHost: "",
	User:        "",
	Password:    "",
	Port:        5672,
	LogLevel:    5,
}

func main() {
	sess := amqp.NewSession(cfg)

	if err := sess.Connect(); err != nil {
		fmt.Println(err)
		return
	}
	defer sess.Close()

	srv, err := sess.Server(srvCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	mongoConfig := users_store.MongoConfig{
		Host:     "localhost",
		Port:     "27017",
		Database: "recommendation_system",
	}
	mongoUserStore, err := users_store.NewMongoStore(mongoConfig)
	if err != nil {
		log.Fatal(err)
	}
	usersService := users_store.NewUserService(mongoUserStore)
	moviesAmqpEndpoints := users_store.NewAMQPEndpointFactory(usersService)

	srv.Endpoint("users.create", moviesAmqpEndpoints.CreateUserAmqpEndpoint())
	srv.Endpoint("users.getByUsername", moviesAmqpEndpoints.GetUserByUsernameAmqpEndpoint())
	fmt.Println("Start server")
	if err := srv.Start(); err != nil {
		fmt.Println(err)
		return
	}
}
