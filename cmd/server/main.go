package main

import (
	"context"
	"log"
	"net"

	"github.com/rdrumond33/go_grpc_example/cmd/server/http"
	"github.com/rdrumond33/go_grpc_example/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	pb.UnimplementedSendEventServer
}

type Event struct {
	gorm.Model
	TypeEvent string
	Context   string
	Price     float64
}

const (
	dsn = "host=pglsb user=rodrigo password=root dbname=events port=5432"
)

func (s *Server) RequestMessage(ctx context.Context, req *pb.Request) (*pb.Status, error) {
	log.Print("Mensagem recebida")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Create
	db.Create(&Event{
		TypeEvent: req.TypeEvent,
		Context:   req.Context,
		Price:     float64(req.Price),
	})
	response := &pb.Status{
		Status: "1",
	}
	return response, nil
}

func (s *Server) FindEvents(ctx context.Context, req *pb.Empty) (*pb.FindResponse, error) {
	log.Print("Buscando todos eventos")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var events []Event
	var response []*pb.Event

	db.Find(&events)
	for _, v := range events {
		response = append(response, &pb.Event{
			ID:        uint32(v.ID),
			TypeEvent: v.TypeEvent,
			Context:   v.Context,
			Price:     float32(v.Price),
		})
	}
	retorno := &pb.FindResponse{Events: response}
	return retorno, nil
}

func (s Server) mustEmbedUnimplementedSendEventServer() {}

func main() {
	webService := http.NewWebServer()
	go webService.Serve()
	log.Print("Iniciando o servidor")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Event{})

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterSendEventServer(grpcServer, &Server{})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}

}
