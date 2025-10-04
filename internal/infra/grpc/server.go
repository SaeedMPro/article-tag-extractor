package grpc

import (
	"context"
	"net"

	"github.com/SaeedMPro/article-tag-extractor/internal/app"
	"github.com/SaeedMPro/article-tag-extractor/internal/domain/entity"
	pb "github.com/SaeedMPro/article-tag-extractor/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedArticleServiceServer
	grpcServer *grpc.Server
	service    *app.ArticleService
}

func NewServer(articleService *app.ArticleService) *Server {
	s := &Server{
		grpcServer: grpc.NewServer(),
		service:    articleService,
	}

	pb.RegisterArticleServiceServer(s.grpcServer, s)
	return s
}

func (s *Server) Serve(lis net.Listener) error {
	return s.grpcServer.Serve(lis)
}

func (s *Server) ProcessArticles(ctx context.Context, req *pb.ProcessArticlesRequest) (*pb.ProcessArticlesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if len(req.Articles) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no articles provided")
	}

	// convert protobuf articles to domain entities
	var articles []*entity.Article
	for _, article := range req.Articles {
		articles = append(articles, &entity.Article{
			Title: article.Title,
			Body:  article.Body,
		})
	}

	// process articles into service
	totalArticleProcessed, err := s.service.ProcessArticles(ctx, articles)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed in process articles: %v", err)
	}

	res := &pb.ProcessArticlesResponse{
		TotalProcessed: int32(totalArticleProcessed),
	}

	return res, nil
}

func (s *Server) GetTopTags(ctx context.Context, req *pb.GetTopTagsRequest) (*pb.GetTopTagsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.Limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "wrong limit provided")
	}

	// get tag frequency from service
	tagFrequencies, err := s.service.GetTopTags(ctx, int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get top tags: %v", err)
	}

	// convert to protobuf
	var pbTags []*pb.TagFrequency
	for _, tf := range tagFrequencies {
		pbTags = append(pbTags, &pb.TagFrequency{
			Tag:       tf.Tag,
			Frequency: int32(tf.Frequency),
		})
	}

	return &pb.GetTopTagsResponse{
		Tags: pbTags,
	}, nil
}

func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}
