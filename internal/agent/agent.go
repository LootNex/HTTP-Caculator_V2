package agent

import (
	"context"
	"errors"
	"log"
	"time"

	pb "github.com/LootNex/HTTP-Caculator_V2/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Solved_Task struct {
	Id             int           `json:"id"`
	Result         float64       `json:"result"`
	Status         string        `json:"status"`
	Operation_time time.Duration `json:"operation_time"`
}

func Count(a, b float64, oper string) (float64, error) {

	switch oper {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	}

	return 0, errors.New("something wrong")

}

func AgentRun() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect %v", err)

	}

	defer conn.Close()

	client := pb.NewCalcServiceClient(conn)
	for {
		resp, err := client.GetTask(context.Background(), &pb.GetTaskRequest{})
		if err != nil || (resp.Arg1 == 0 && resp.Arg2 == 0) {
			time.Sleep(time.Second)
			continue
		}
		start := time.Now()
		result, err := Count(resp.Arg1, resp.Arg2, resp.Operation)
		if err != nil {
			log.Printf("failed to count %v", err)
		}

		ResultRequst := pb.SendResultRequest{
			Id:            resp.Id,
			Result:        result,
			Status:        "counted",
			OperationTime: int64(time.Since(start).Microseconds()),
		}

		_, err = client.SendResult(context.Background(), &ResultRequst)
		if err != nil {
			log.Fatalf("failed to send result %v", err)
		}

	}

}
