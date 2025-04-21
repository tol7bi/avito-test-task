package grpc

import (
	"context"
	"fmt"
	"net"
	"time"
	pb "pvz-backend/internal/grpc/gen"
	"pvz-backend/internal/models"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	
)

type server struct {
	pb.UnimplementedPVZServiceServer
}

func (s *server) GetPVZList(ctx context.Context, req *pb.GetPVZListRequest) (*pb.GetPVZListResponse, error) {
	rows, err := models.DB.Query(ctx, `
		SELECT id, registration_date, city FROM pvz
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*pb.PVZ
	for rows.Next() {
		var id, city string
		var regDate time.Time
	
		if err := rows.Scan(&id, &regDate, &city); err != nil {
			return nil, err
		}
	
		result = append(result, &pb.PVZ{
			Id:               id,
			City:             city,
			RegistrationDate: timestamppb.New(regDate),
		})
	}

	return &pb.GetPVZListResponse{Pvzs: result}, nil
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	s := grpc.NewServer()
	pb.RegisterPVZServiceServer(s, &server{})
	reflection.Register(s)

	fmt.Println("🛰  gRPC сервер запущен на порту 3000")
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
