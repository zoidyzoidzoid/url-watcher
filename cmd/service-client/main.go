package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	proto "github.com/zoidbergwill/url-watcher/pkg/proto"
	grpc "google.golang.org/grpc"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:8081", "gRPC server endpoint")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*grpcServerEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println(conn)
	client := proto.NewWatcherServiceClient(conn)
	fmt.Println(client)
	resp, err := client.ListFoods(context.Background(), &proto.FoodRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.String())
	// stream, err := client.WatchFoods(context.Background(), &proto.FoodRequest{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for {
	// 	character, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("%v.ListCharacters(_) = _, %v", client, err)
	// 	}
	// 	log.Println(character.String())
	// }
}
