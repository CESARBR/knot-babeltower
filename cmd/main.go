package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/CESARBR/knot-babeltower/internal/config"
	thingInteractors "github.com/CESARBR/knot-babeltower/pkg/interactors"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/server"
	"github.com/CESARBR/knot-babeltower/pkg/user/controllers"
	userDeliveryHTTP "github.com/CESARBR/knot-babeltower/pkg/user/delivery/http"
	userInteractors "github.com/CESARBR/knot-babeltower/pkg/user/interactors"

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

	// AMQP Publisher
	msgPublisher := network.NewMsgPublisher(logrus.Get("MsgPublisher"), amqp)

	// Services
	userProxy := userDeliveryHTTP.NewUserProxy(logrus.Get("UserProxy"), config.Users.Hostname, config.Users.Port)
	thingProxy := network.NewThingProxy(logrus.Get("ThingProxy"), config.Things.Hostname, config.Things.Port)
	connector := network.NewConnector(logrus.Get("Connector"), amqp)

	// Interactors
	createUser := userInteractors.NewCreateUser(logrus.Get("CreateUser"), userProxy)
	createToken := userInteractors.NewCreateToken(logrus.Get("CreateToken"), userProxy)
	registerThing := thingInteractors.NewRegisterThing(logrus.Get("RegisterThing"), msgPublisher, thingProxy, connector)

	// Controllers
	userController := controllers.NewUserController(logrus.Get("Controller"), createUser, createToken)

	// Server
	serverChan := make(chan bool, 1)
	server := server.NewServer(config.Server.Port, logrus.Get("Server"), userController)

	// AMQP Handler
	msgChan := make(chan bool, 1)
	msgHandler := network.NewMsgHandler(logrus.Get("MsgHandler"), amqp, registerThing)

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
				go msgHandler.Start(msgChan)
			}
		case started := <-msgChan:
			if started {
				logger.Info("Msg handler started")
			} else {
				quit <- true
			}
		case <-quit:
			msgHandler.Stop()
			amqp.Stop()
			server.Stop()
		}
	}
}
