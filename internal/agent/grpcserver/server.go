package server

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	pb "github.com/LootNex/HTTP-Caculator_V2/internal/proto"
	calculator "github.com/LootNex/HTTP-Caculator_V2/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedCalcServiceServer
}

func (s *Server) GetTask(ctx context.Context, _ *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {

	if len(calculator.Tasks) == 0 {
		return nil, errors.New("no tasks")
	}

	return &pb.GetTaskResponse{
		Id:        int64(calculator.Tasks[0].Id),
		Arg1:      calculator.Tasks[0].Arg1,
		Arg2:      calculator.Tasks[0].Arg2,
		Operation: calculator.Tasks[0].Operation,
	}, nil

}

func (s *Server) SendResult(ctx context.Context, res *pb.SendResultRequest) (*pb.SendResultResponse, error) {

	calculator.Tasks[0].Operation_time = time.Duration(res.OperationTime)

	calculator.Task_Ready <- res.Result

	return &pb.SendResultResponse{}, nil
}

func StartServer() error {

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		return err
	}

	GrpcServer := grpc.NewServer()
	reflection.Register(GrpcServer)

	pb.RegisterCalcServiceServer(GrpcServer, &Server{})

	log.Println("grpc server is running")

	err = GrpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil

}
