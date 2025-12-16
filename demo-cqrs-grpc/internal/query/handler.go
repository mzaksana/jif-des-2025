package query

import (
	"context"
	"log"

	pb "demo-cqrs-grpc/proto"
	"demo-cqrs-grpc/internal/store"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler implements the QueryService gRPC server
type Handler struct {
	pb.UnimplementedQueryServiceServer
	readStore *store.ReadStore
}

// NewHandler creates a new query handler
func NewHandler(rs *store.ReadStore) *Handler {
	return &Handler{readStore: rs}
}

// GetPost retrieves a single post by ID
func (h *Handler) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	log.Printf("[QUERY] GetPost: id=%q", req.Id)

	post, exists := h.readStore.GetPost(req.Id)
	if !exists {
		return nil, status.Errorf(codes.NotFound, "post not found: %s", req.Id)
	}

	return &pb.GetPostResponse{
		Post: &pb.Post{
			Id:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			Author:    post.Author,
			Tags:      post.Tags,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		},
	}, nil
}

// ListPosts returns paginated posts
func (h *Handler) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := int(req.Offset)

	log.Printf("[QUERY] ListPosts: limit=%d, offset=%d", limit, offset)

	posts, total := h.readStore.ListPosts(limit, offset)

	pbPosts := make([]*pb.Post, len(posts))
	for i, p := range posts {
		pbPosts[i] = &pb.Post{
			Id:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			Author:    p.Author,
			Tags:      p.Tags,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	return &pb.ListPostsResponse{
		Posts: pbPosts,
		Total: int32(total),
	}, nil
}

// SearchPosts searches posts by query and tags
func (h *Handler) SearchPosts(ctx context.Context, req *pb.SearchPostsRequest) (*pb.SearchPostsResponse, error) {
	log.Printf("[QUERY] SearchPosts: query=%q, tags=%v", req.Query, req.Tags)

	posts := h.readStore.SearchPosts(req.Query, req.Tags)

	pbPosts := make([]*pb.Post, len(posts))
	for i, p := range posts {
		pbPosts[i] = &pb.Post{
			Id:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			Author:    p.Author,
			Tags:      p.Tags,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	return &pb.SearchPostsResponse{Posts: pbPosts}, nil
}
