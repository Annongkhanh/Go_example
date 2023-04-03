package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Annongkhanh/Simple_bank/api"
	db "github.com/Annongkhanh/Simple_bank/db/sqlc"
	_ "github.com/Annongkhanh/Simple_bank/doc/statik"
	"github.com/Annongkhanh/Simple_bank/gapi"
	"github.com/Annongkhanh/Simple_bank/pb"
	"github.com/Annongkhanh/Simple_bank/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".")

	if config.Environment == "development" {
		// Human friendly logging
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if err != nil {
		log.Error().Err(err).Msg("Can not load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Error().Err(err).Msg("Can not connect to database")
	}

	store := db.NewStore(conn)

	go runGatewayServer(config, store)

	runGrpcServer(config, store)

}

func runGrpcServer(config util.Config, store db.Store) {

	grpcLogger := grpc.UnaryInterceptor(gapi.GRPCLogger)

	grpcServer := grpc.NewServer(grpcLogger)
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Error().Err(err).Msg("Can not initialize server")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	pb.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)

	if err != nil {
		log.Error().Err(err).Msg("can not create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Error().Err(err).Msg("can not start gRPC server")
	}

}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Error().Err(err).Msg("can not initialize server")
	}

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Error().Err(err).Msg("can not register handler server")
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Error().Err(err).Msg("failed to create statik filesystem")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)

	if err != nil {
		log.Error().Err(err).Msg("can not create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handler := gapi.HTTPLogger(mux)
	err = http.Serve(listener, handler)

	if err != nil {
		log.Error().Err(err).Msg("can not start HTTP gateway server")
	}

}

func runGinServer(config util.Config, store db.Store) {

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Error().Err(err).Msg("Can not initialize server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("Can not start server")
	}
}

// Run DB migration
func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Error().Err(err).Msg("Can not create new migrate instance")
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("Can not migrate database")
	}

	log.Info().Msg("Success to migrate database")

}
