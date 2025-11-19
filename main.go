package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"

	pb "github.com/ssuji15/wolf-worker/agent"

	"google.golang.org/grpc"
)

const (
	socketPath = "/socket/socket.sock"
	outputPath = "/output/output.log"
)

type WorkerAgent struct {
	pb.UnimplementedWorkerAgentServer
}

func main() {

	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWorkerAgentServer(grpcServer, &WorkerAgent{})

	fmt.Println("worker listening on", socketPath)
	grpcServer.Serve(lis)
}

func (w *WorkerAgent) StartJob(ctx context.Context, req *pb.JobRequest) (*pb.Ack, error) {
	go func() {
		os.Remove(socketPath)
		runJob(req.Engine, req.Code)
		os.Exit(0)
	}()
	return &pb.Ack{Message: "ok"}, nil
}

func runJob(engine, code string) {
	src := getFileName(engine)
	os.WriteFile(src, []byte(code), 0644)

	f, _ := os.Create(outputPath)
	f.Close()

	switch engine {
	case "c++":
		run("g++", "-O1", "-pipe", "-g0", src, "-o", "/tmp/prog")
		run("/tmp/prog")
	case "java":
		run("javac", src)
		run("java", "-cp", "/tmp", "Program")
	default:
		os.WriteFile(outputPath, []byte("unknown engine"), 0644)
	}
}

func getFileName(engine string) string {
	switch engine {
	case "c++":
		return "/tmp/p.cpp"
	case "java":
		return "/tmp/p.java"
	default:
		return "/tmp/p.txt"
	}
}

func run(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()

	f, _ := os.OpenFile(outputPath, os.O_WRONLY|os.O_APPEND, 0644)
	defer f.Close()
	f.Write(out)
	if err != nil {
		f.Write([]byte(err.Error()))
	}
}
