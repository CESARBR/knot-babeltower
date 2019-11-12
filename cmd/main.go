package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/CESARBR/knot-babeltower/internal/config"
	"github.com/CESARBR/knot-babeltower/pkg/controllers"
	"github.com/CESARBR/knot-babeltower/pkg/interactors"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/server"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

func monitorSignals(sigs chan os.Signal, quit chan bool, logger logging.Logger) {
	signal := <-sigs
	logger.Infof("Signal %s received", signal)
	quit <- true
}

func main() {
	config := config.Load()
	logrus := logging.NewLogrus(config.Logger.Level)

	logger := logrus.Get("Main")
	logger.Info("Starting KNoT Babeltower")

	// Signal Handler
	sigs := make(chan os.Signal, 1)
	quit := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go monitorSignals(sigs, quit, logger)

	// AMQP
	amqpChan := make(chan bool, 1)
	amqp := network.NewAmqp(config.RabbitMQ.URL, logrus.Get("Amqp"))

	// Services
	userProxy := network.NewUserProxy(logrus.Get("UserProxy"), config.Users.Hostname, config.Users.Port)

	// Interactors
	createUser := interactors.NewCreateUser(logrus.Get("CreateUser"), userProxy)

	// Controllers
	userController := controllers.NewUserController(logrus.Get("Controller"), createUser)

	// Server
	serverChan := make(chan bool, 1)
	server := server.NewServer(config.Server.Port, logrus.Get("Server"), userController)

	// Start goroutines
	go amqp.Start(amqpChan)
	go server.Start(serverChan)

	// Main loop
	for {
		select {
		case started := <-serverChan:
			if started {
				logger.Info("Server started")
			}
		case started := <-amqpChan:
			if started {
				logger.Info("AMQP connection started")
			}
		case <-quit:
			amqp.Stop()
			server.Stop()
		}
	}
}
