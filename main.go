package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"time"

	"github.com/bodocoder/grpc-calci/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultInputValue = 0
)

var (
	addr = flag.String("addr", "localhost:3000", "the address to connect to")
	x    = flag.Float64("x", defaultInputValue, "first input")
	y    = flag.Float64("y", defaultInputValue, "second input")
	op   = flag.String("op", "+", "operation to be performed")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCalculatorServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var r *pb.CalculateResponse
	switch *op {
	case "+":
		r, err = c.PerformAddition(ctx, &pb.CalculateRequest{X: *x, Y: *y})
	case "-":
		r, err = c.PerformSubtraction(ctx, &pb.CalculateRequest{X: *x, Y: *y})
	case "*":
		r, err = c.PerformMultiplication(ctx, &pb.CalculateRequest{X: *x, Y: *y})
	case "/":
		r, err = c.PerformDivision(ctx, &pb.CalculateRequest{X: *x, Y: *y})
	default:
		r, err = nil, errors.New("this operation is not available!")
	}

	if err != nil {
		log.Fatalf("could not perform: %v", err)
	}
	log.Printf("got result: %s", r.GetRes())
}
