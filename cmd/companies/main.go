package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/itzurabhi/companies-micro/internal/handlers"
	"github.com/itzurabhi/companies-micro/internal/logic"
	"github.com/itzurabhi/companies-micro/internal/repositories"
	kafkaRepos "github.com/itzurabhi/companies-micro/internal/repositories/kafka"
	pgRepos "github.com/itzurabhi/companies-micro/internal/repositories/postgres"
	"github.com/itzurabhi/companies-micro/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	postgresInternals "github.com/itzurabhi/companies-micro/internal/database/postgres"
)

var ListenPort = utils.EnvOrDefault("PORT", "8081")
var ListenHost = utils.EnvOrDefault("HOST", "0.0.0.0")

type server struct {
	// conections
	pgdb          *gorm.DB
	kafkaProducer *kafka.Producer

	// repositories
	companiesRepo repositories.Companies
	companyEvents repositories.EventBus

	// logic layers
	companyLogic *logic.CompanyLogic

	// handlers
	companyHandler handlers.CompanyHandler
}

func (srv *server) createPostgresConn() *server {
	dsn := utils.EnvOrDefault("POSTGRES_DSN", "")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("creating postgres connection failed", err)
	}
	srv.pgdb = db
	return srv
}

func (srv *server) migratePostgres() *server {
	if srv.pgdb == nil {
		logrus.Fatal("database connection must be created before migrations")
	}

	if err := postgresInternals.Migrate(srv.pgdb); err != nil {
		log.Fatal(err)
	}
	return srv
}

func (srv *server) createKafkaConnn() *server {
	prod, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": utils.EnvOrDefault("KAFKA_CLIENT_BOOTSTRAP_SERVERS", ""),
		"client.id":         utils.EnvOrDefault("KAFKA_CLIENT_ID", ""),
		"acks":              "all"})
	if err != nil {
		logrus.Fatal("creating kafka connection failed", err)
	}
	srv.kafkaProducer = prod
	return srv
}

func (srv *server) createPostgresRepos() *server {
	if srv.pgdb == nil {
		logrus.Fatal("database connection must be created before repositories")
	}
	srv.companiesRepo = pgRepos.CreateCompaniesRepo(srv.pgdb)
	return srv
}

func (srv *server) createKafkaRepos() *server {
	if srv.kafkaProducer == nil {
		logrus.Fatal("kafka connection must be created before repositories")
	}
	srv.companyEvents = kafkaRepos.CreateCompaniesEventBus(srv.kafkaProducer, "CompanyEvents")
	return srv
}

func (srv *server) createLogicLayers() *server {
	if srv.companiesRepo == nil || srv.companyEvents == nil {
		logrus.Fatal("repositories must be initialized before logic layers")
	}

	srv.companyLogic = logic.CreateCompanyLogic(srv.companiesRepo, srv.companyEvents)

	return srv
}

func (srv *server) createFiberHandlers() *server {
	if srv.companyLogic == nil {
		logrus.Fatal("logic must be initialized before handlers")
	}
	srv.companyHandler = *handlers.CreateCompanyHandler(srv.companyLogic)
	return srv
}

func (srv *server) cleanupResources() {
	if srv.pgdb != nil {
		logrus.Println("closing postgres connection")
		conn, err := srv.pgdb.DB()

	}
}

func (srv *server) AddRoutes(app *fiber.App) *fiber.App {

	// health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(map[string]string{
			"status": "ok",
		})
	})

	// companies route

	companiesRoute := app.Group("companies")

	companiesRoute.Get("/:id", srv.companyHandler.Get)
	companiesRoute.Post("/", srv.companyHandler.Create)
	companiesRoute.Patch("/:id", srv.companyHandler.Patch)
	companiesRoute.Delete("/:id", srv.companyHandler.Delete)

	return app
}

func main() {

	server := new(server)

	server.createPostgresConn()
	server.createKafkaConnn()

	server.migratePostgres()

	server.createPostgresRepos()
	server.createKafkaRepos()

	server.createLogicLayers()

	server.createFiberHandlers()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	app := fiber.New(fiber.Config{
		ServerHeader: "Companies Serverice",
	})

	app = server.AddRoutes(app)

	go func() {
		<-sigCh

		logrus.Println("interrupt recieved. shutting down...")

		if err := app.Shutdown(); err != nil {
			log.Fatal("App shutdown error", err)
		}

	}()

	hostPort := ListenHost + ":" + ListenPort
	logrus.Info("listening on :", hostPort)
	if err := app.Listen(hostPort); err != nil {
		logrus.Fatal(err)
	}

	logrus.Println("cleanup started")
	logrus.Println("cleanup completed")
}
