package env

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sethvargo/go-envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/ptsypyshev/gb-golang-level3-new/internal/apigw/routes"
	v1 "github.com/ptsypyshev/gb-golang-level3-new/internal/apigw/v1"
	"github.com/ptsypyshev/gb-golang-level3-new/internal/database/links"
	"github.com/ptsypyshev/gb-golang-level3-new/internal/database/users"
	"github.com/ptsypyshev/gb-golang-level3-new/internal/env/config"
	"github.com/ptsypyshev/gb-golang-level3-new/internal/link/linkgrpc"
	"github.com/ptsypyshev/gb-golang-level3-new/internal/link/stories/linkupdater"
	"github.com/ptsypyshev/gb-golang-level3-new/internal/user/usergrpc"

	"github.com/ptsypyshev/gb-golang-level3-new/pkg/pb"
)

type Env struct {
	Config          config.Config
	APIGWHTTPServer *http.Server
	LinksGRPCServer *grpc.Server
	UsersGRPCServer *grpc.Server
	LinkUpdater     *linkupdater.Story
}

func Setup(ctx context.Context) (*Env, *Closer, error) {
	var cfg config.Config
	env := &Env{}

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, nil, fmt.Errorf("env processing: %w", err)
	}

	linksDBConn, err := mongo.Connect(
		ctx, &options.ClientOptions{
			ConnectTimeout: &cfg.LinksService.Mongo.ConnectTimeout,
			Hosts:          []string{fmt.Sprintf("%s:%d", cfg.LinksService.Mongo.Host, cfg.LinksService.Mongo.Port)},
			MaxPoolSize:    &cfg.LinksService.Mongo.MaxPoolSize,
			MinPoolSize:    &cfg.LinksService.Mongo.MinPoolSize,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("mongo.Connect: %w", err)
	}

	usersDBConn, err := pgxpool.Connect(ctx, cfg.UsersService.Postgres.ConnectionURL())
	if err != nil {
		return nil, nil, fmt.Errorf("pgxpool Connect: %w", err)
	}

	amqpConn, err := amqp.Dial(cfg.LinksService.AMQP.String())
	if err != nil {
		return nil, nil, fmt.Errorf("amqp Dial: %w", err)
	}

	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("amqp Channel: %w", err)
	}

	// задекларируйте очередь
	_, err = amqpChannel.QueueDeclare(cfg.LinksService.AMQP.QueueName, false, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("QueueDeclare: %w", err)
	}

	usersRepository := users.New(usersDBConn, 5*time.Second) // вынести в конфиг duration
	linksRepository := links.New(
		linksDBConn.Database(cfg.LinksService.Mongo.Name),
		5*time.Second, // вынести в конфиг duration
	)

	{
		handler := linkgrpc.New(linksRepository, cfg.LinksService.GRPCServer.Timeout, amqpChannel, cfg.LinksService.AMQP.QueueName)

		s := grpc.NewServer()
		reflection.Register(s) // этот код нужен для дебаггинга
		pb.RegisterLinkServiceServer(s, handler)

		// grpc server start function
		env.LinksGRPCServer = s
	}

	{
		handler := usergrpc.New(usersRepository, cfg.LinksService.GRPCServer.Timeout)

		s := grpc.NewServer()
		reflection.Register(s) // этот код нужен для дебаггинга
		pb.RegisterUserServiceServer(s, handler)

		// grpc server start function
		env.UsersGRPCServer = s
	}

	// Инициализируем клиенты GRPC

	// Клиент для осуществления запросов в users service
	usersClientConn, err := grpc.DialContext(
		ctx, cfg.APIGWService.UsersClientAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("grpc DialContext: %w", err)
	}

	usersClient := pb.NewUserServiceClient(usersClientConn)

	// Клиент для осуществления запросов в links service
	linksClientConn, err := grpc.DialContext(
		ctx, cfg.APIGWService.LinksClientAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("grpc DialContext: %w", err)
	}

	linksClient := pb.NewLinkServiceClient(linksClientConn)

	// API GW handler
	// В роуйтере пакета v1 нужно использовать клиенты и запрашивать данные с сервисов links и users
	handler := v1.New(usersClient, linksClient)
	router := routes.Router(handler)

	apiGWServer := &http.Server{
		Addr:              cfg.APIGWService.Addr,
		Handler:           router,
		ReadTimeout:       cfg.APIGWService.ReadTimeout,
		ReadHeaderTimeout: cfg.APIGWService.ReadTimeout,
		WriteTimeout:      cfg.APIGWService.WriteTimeout,
		IdleTimeout:       cfg.APIGWService.ReadTimeout,
	}

	linkUpdaterStory := linkupdater.New(linksRepository, amqpChannel, cfg.LinksService.AMQP.QueueName)

	env.APIGWHTTPServer = apiGWServer
	env.Config = cfg
	env.LinkUpdater = linkUpdaterStory

	return env, NewCloser(usersDBConn, linksDBConn, amqpConn, amqpChannel), nil
}
