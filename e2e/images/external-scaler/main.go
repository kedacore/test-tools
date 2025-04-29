package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/gorilla/mux"
	pb "github.com/kedacore/test-tools/external-scaler/externalscaler"
)

func setValue(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	number := vars["number"]
	ExternalScalerValue, _ = strconv.ParseInt(number, 10, 64)
	log.Printf("new int value: %d\n", ExternalScalerValue)
	if ExternalScalerValue < 0 {
		log.Print("negative -> it won't be included in GetMetricsResponse")
	}
}

func setFloatValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	number := vars["number"]
	ExternalScalerValueFloat, _ = strconv.ParseFloat(number, 64)
	log.Printf("new float value: %f\n", ExternalScalerValueFloat)
	if ExternalScalerValueFloat < 0 {
		log.Print("negative -> it won't be included in GetMetricsResponse")
	}
}

func RunManagementApi() {
	r := mux.NewRouter()
	r.HandleFunc("/api/value/{number:[-0-9]+}", setValue).Methods("POST")
	r.HandleFunc("/api/floatvalue/{number:[-\\.0-9]+}", setFloatValue).Methods("POST")
	r.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
	http.Handle("/", r)
	fmt.Printf("Running http management server on port: %d\n", 8082)
	fmt.Print("example usage:\n - POST -> localhost:8082/api/value/3\n - POST -> localhost:8082/api/floatvalue/3.14\n")
	fmt.Print("if you set one of those values as negative, it won't be sent in the payload\n")
	http.ListenAndServe(":8082", nil)
}

var ExternalScalerValue int64 = 0
var ExternalScalerValueFloat = .0

type ExternalScaler struct {
	pb.UnimplementedExternalScalerServer
}

func (es *ExternalScaler) IsActive(ctx context.Context, scaledObjectRef *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {
	log.Println("Executing method IsActive")

	return &pb.IsActiveResponse{Result: ExternalScalerValue > 0}, nil
}

func (es *ExternalScaler) GetMetricSpec(ctx context.Context, scaledObjectRef *pb.ScaledObjectRef) (*pb.GetMetricSpecResponse, error) {
	log.Println("Executing method GetMetricSpec")

	metricThreshold, err := strconv.ParseFloat(scaledObjectRef.ScalerMetadata["metricThreshold"], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value for metric threshold - %s", err)
	}

	return &pb.GetMetricSpecResponse{
		MetricSpecs: []*pb.MetricSpec{
			{MetricName: "external-scaler-e2e-test", TargetSizeFloat: metricThreshold},
		},
	}, nil
}

func (es *ExternalScaler) GetMetrics(ctx context.Context, metricRequest *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	log.Println("Executing method GetMetrics")

	if metricRequest.MetricName != "external-scaler-e2e-test" {
		return nil, fmt.Errorf("invalid metric name - %s", metricRequest.MetricName)
	}
	mv := &pb.MetricValue{MetricName: "external-scaler-e2e-test"}
	if ExternalScalerValue >= 0 {
		mv.MetricValue = ExternalScalerValue
	}
	if ExternalScalerValueFloat >= 0 {
		mv.MetricValueFloat = ExternalScalerValueFloat
	}

	return &pb.GetMetricsResponse{
		MetricValues: []*pb.MetricValue{mv},
	}, nil
}

func (es *ExternalScaler) StreamIsActive(scaledObjectRef *pb.ScaledObjectRef, epsServer pb.ExternalScaler_StreamIsActiveServer) error {
	log.Println("Executing method StreamIsActive")

	return nil
}

func main() {
	go RunManagementApi()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     30 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
			MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
			MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
			Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
			Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	lis, _ := net.Listen("tcp", ":6000")

	pb.RegisterExternalScalerServer(grpcServer, &ExternalScaler{})
	log.Println("Listening external scaler on :6000")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
