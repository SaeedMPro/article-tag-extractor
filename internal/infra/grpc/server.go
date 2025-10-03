package grpc

import (
	"context"
	"net"

	"github.com/SaeedMPro/article-tag-extractor/internal/app"
	pb "github.com/SaeedMPro/article-tag-extractor/internal/proto"
	"google.golang.org/grpc"
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

func (s *Server) ProcessArticles(context.Context, *pb.ProcessArticlesRequest) (*pb.ProcessArticlesResponse, error) {
	//TODO
	return nil, nil
}

func (s *Server) GetTopTags(context.Context, *pb.GetTopTagsRequest) (*pb.GetTopTagsResponse, error) {
	//TODO
	return nil, nil
}

func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}
