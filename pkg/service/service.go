package service

import (
	"context" // Use "golang.org/x/net/context" for Golang version <= 1.6
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"log"
	"net"

	gw "github.com/zoidbergwill/url-watcher/pkg/proto"
	pb "github.com/zoidbergwill/url-watcher/pkg/proto"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint        = flag.String("grpc-server-endpoint", "localhost:8081", "gRPC server endpoint")
	grpcGatewayServerEndpoint = flag.String("grpc-gateway-server-endpoint", "localhost:8080", "gRPC server endpoint")
)

func actuallyRunGRPCGatewayServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterWatcherServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	fmt.Printf("Starting http GRPC gateway server on %s\n", *grpcGatewayServerEndpoint)
	return http.ListenAndServe(*grpcGatewayServerEndpoint, mux)
}

// RunGRPCGatewayServer should have a comment
func RunGRPCGatewayServer() {
	flag.Parse()
	defer glog.Flush()

	if err := actuallyRunGRPCGatewayServer(); err != nil {
		glog.Fatal(err)
	}
}

// server is used to implement service.StarWarsServer.
type server struct {
	storage watchedURL
}

// ListFoods should have a comment
func (s *server) ListFoods(ctx context.Context, req *pb.FoodRequest) (*pb.FoodResponse, error) {
	version := s.storage.latestVersion
	foods := s.storage.snapshots[version]
	return &pb.FoodResponse{Version: version, Items: foods}, nil
}

// WatchFoods should have a comment
func (s *server) WatchFoods(req *pb.FoodRequest, stream pb.WatcherService_WatchFoodsServer) error {
	// log.Printf("Received: %v", in.GetValue())
	version := s.storage.latestVersion
	foods := s.storage.snapshots[version]
	if err := stream.Send(&pb.FoodResponse{Version: version, Items: foods}); err != nil {
		return err
	}
	return nil
	// for _, food := range foods {
	// 	if err := stream.Send(food); err != nil {
	// 		return err
	// 	}
	// }
	// return nil
}

// FoodsResponse needs a comment
type FoodsResponse []*pb.Food

// UnmarshalFoodsResponse needs a comment
func UnmarshalFoodsResponse(data []byte) (FoodsResponse, error) {
	var r FoodsResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal needs a comment
func (r *FoodsResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// WE EVENT SOURCING NOW
type watchedURL struct {
	latestVersion uint64
	deltas        map[uint64][]*pb.Food
	snapshots     map[uint64][]*pb.Food
}

// KeepFoodsFresh need a comment
func KeepFoodsFresh(state *server) {
	for {
		log.Print("Fetching foods")
		// Issue an HTTP GET request to a server. `http.Get` is a
		// convenient shortcut around creating an `http.Client`
		// object and calling its `Get` method; it uses the
		// `http.DefaultClient` object which has useful default
		// settings.
		resp, err := http.Get("http://127.0.0.1/foods")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// Print the HTTP response status.
		log.Print("Response status:", resp.Status)

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Print(string(bytes))

		data, err := UnmarshalFoodsResponse(bytes)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v\n", data)

		latestVersion := state.storage.latestVersion
		if latestVersion == 0 {
			state.storage.deltas = make(map[uint64][]*pb.Food)
			state.storage.snapshots = make(map[uint64][]*pb.Food)
			state.storage.snapshots[1] = data
			state.storage.latestVersion = 1
		} else {
			old := state.storage.snapshots[latestVersion]
			delta := make([]*pb.Food, 0)
			for _, food2 := range data {
				found := false
				for _, food := range old {
					if reflect.DeepEqual(food, food2) {
						found = true
						break
					}
				}
				if !found {
					delta = append(delta, food2)
				}
			}
			log.Printf("Delta %+v", delta)
			if len(delta) == 0 {
				log.Print("No change to data")
			} else {
				state.storage.snapshots[latestVersion+1] = data
				state.storage.deltas[latestVersion+1] = delta
				state.storage.latestVersion = latestVersion + 1
				log.Printf("Updated to version %d", latestVersion+1)
			}
		}
		log.Print("Sleeping for 10s...")
		time.Sleep(time.Second * 10)
	}
}

// RunGRPCServer should have a comment
func RunGRPCServer() {
	state := &server{}
	go KeepFoodsFresh(state)
	fmt.Printf("Starting grpc server on %s\n", *grpcServerEndpoint)
	lis, err := net.Listen("tcp", *grpcServerEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWatcherServiceServer(s, state)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
