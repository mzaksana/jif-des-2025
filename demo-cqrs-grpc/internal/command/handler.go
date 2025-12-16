package command

import (
	"context"
	"log"

	pb "demo-cqrs-grpc/proto"
	"demo-cqrs-grpc/internal/store"
)

// Handler implements the CommandService gRPC server
type Handler struct {
	pb.UnimplementedCommandServiceServer
	writeStore *store.WriteStore
}

// NewHandler creates a new command handler
func NewHandler(ws *store.WriteStore) *Handler {
	return &Handler{writeStore: ws}
}

// CreatePost handles post creation
func (h *Handler) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	log.Printf("[COMMAND] CreatePost: title=%q, author=%q, tags=%v", req.Title, req.Author, req.Tags)

	id, err := h.writeStore.CreatePost(req.Title, req.Content, req.Author, req.Tags)
	if err != nil {
		log.Printf("[COMMAND] CreatePost error: %v", err)
		return &pb.CreatePostResponse{Success: false}, err
	}

	return &pb.CreatePostResponse{Id: id, Success: true}, nil
}

// UpdatePost handles post updates
func (h *Handler) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	log.Printf("[COMMAND] UpdatePost: id=%q, title=%q", req.Id, req.Title)

	err := h.writeStore.UpdatePost(req.Id, req.Title, req.Content, req.Tags)
	if err != nil {
		log.Printf("[COMMAND] UpdatePost error: %v", err)
		return &pb.UpdatePostResponse{Success: false}, err
	}

	return &pb.UpdatePostResponse{Success: true}, nil
}

// DeletePost handles post deletion
func (h *Handler) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	log.Printf("[COMMAND] DeletePost: id=%q", req.Id)

	err := h.writeStore.DeletePost(req.Id)
	if err != nil {
		log.Printf("[COMMAND] DeletePost error: %v", err)
		return &pb.DeletePostResponse{Success: false}, err
	}

	return &pb.DeletePostResponse{Success: true}, nil
}
