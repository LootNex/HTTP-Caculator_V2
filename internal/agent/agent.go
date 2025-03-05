package agent

import (
	calculator "Calculator_V2/pkg"
	config "Calculator_V2/pkg/config"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
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
	port := config.New()
	for {
		resp, err := http.Get("http://localhost:" + port.Port + "/internal/task")
		if err != nil {
			time.Sleep(20 * time.Millisecond)
			continue
		}
		start := time.Now()
		task := new(calculator.NewTask)
		json.NewDecoder(resp.Body).Decode(&task)
		if task.Arg1 == 0 && task.Arg2 == 0 {
			continue
		}
		solvedtask := new(Solved_Task)
		solvedtask.Id = task.Id
		solvedtask.Result, err = Count(task.Arg1, task.Arg2, task.Operation)
		log.Println("FUNC COUNT AGENT", solvedtask.Result, err, task.Arg1, task.Arg2, task.Operation)
		if err != nil {
			solvedtask.Status = err.Error()
		} else {
			solvedtask.Status = "success"
		}
		solvedtask.Operation_time = time.Since(start)
		json_solved_task, err := json.Marshal(solvedtask)
		if err != nil {
			log.Print("wrong json")
		}

		_, err = http.Post("http://localhost:"+port.Port+"/internal/task", "application/json", bytes.NewBufferString(string(json_solved_task)))
		if err != nil {
			log.Print("something wrong with agent post")
		}
		resp.Body.Close()
		time.Sleep(20 * time.Millisecond)

	}
}
