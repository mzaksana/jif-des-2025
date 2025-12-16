package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"demo-cqrs-grpc/internal/command"
	"demo-cqrs-grpc/internal/event"
	"demo-cqrs-grpc/internal/query"
	"demo-cqrs-grpc/internal/store"
	pb "demo-cqrs-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("========================================")
	log.Println("  CQRS + gRPC Blog Service Demo")
	log.Println("  JIF USK x Twibbonize Workshop")
	log.Println("========================================")

	// Initialize event bus for syncing write -> read
	bus := event.NewBus()

	// Initialize write store (SQLite)
	writeStore, err := store.NewWriteStore(bus)
	if err != nil {
		log.Fatalf("Failed to initialize write store: %v", err)
	}
	defer writeStore.Close()

	// Initialize read store (In-memory)
	readStore := store.NewReadStore(bus)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	commandHandler := command.NewHandler(writeStore)
	queryHandler := query.NewHandler(readStore)

	pb.RegisterCommandServiceServer(grpcServer, commandHandler)
	pb.RegisterQueryServiceServer(grpcServer, queryHandler)

	// Enable reflection for grpcurl
	reflection.Register(grpcServer)

	// Start listening
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("")
	log.Println("Server started on :50051")
	log.Println("")
	log.Println("Architecture:")
	log.Println("  [Client] --gRPC--> [CommandService] --write--> [SQLite]")
	log.Println("                            |")
	log.Println("                       event sync")
	log.Println("                            |")
	log.Println("                            v")
	log.Println("  [Client] --gRPC--> [QueryService] <--read--- [In-Memory]")
	log.Println("")
	log.Println("Ready to accept requests...")
	log.Println("========================================")

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("\nShutting down gracefully...")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
