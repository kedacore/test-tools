package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
	pb "github.com/kedacore/test-tools/external-scaler/externalscaler"
)

func setValue(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	number := vars["number"]
	ExternalScalerValue, _ = strconv.ParseInt(number, 10, 64)
	log.Printf("new value: %d\n", ExternalScalerValue)
}

func RunManagementApi() {
	r := mux.NewRouter()
	r.HandleFunc("/api/value/{number:[0-9]+}", setValue).Methods("POST")
	http.Handle("/", r)
	fmt.Printf("Running http management server on port: %d\n", 8080)
	http.ListenAndServe(":8080", nil)
}

var ExternalScalerValue int64 = 0

type ExternalScaler struct {
	pb.UnimplementedExternalScalerServer
}

func (es *ExternalScaler) IsActive(ctx context.Context, scaledObjectRef *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {
	log.Println("Executing method IsActive")

	return &pb.IsActiveResponse{Result: ExternalScalerValue > 0}, nil
}

func (es *ExternalScaler) GetMetricSpec(ctx context.Context, scaledObjectRef *pb.ScaledObjectRef) (*pb.GetMetricSpecResponse, error) {
	log.Println("Executing method GetMetricSpec")

	metricThreshold, err := strconv.ParseInt(scaledObjectRef.ScalerMetadata["metricThreshold"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value for metric threshold - %s", err)
	}

	return &pb.GetMetricSpecResponse{
		MetricSpecs: []*pb.MetricSpec{
			{MetricName: "external-scaler-e2e-test", TargetSize: metricThreshold},
		},
	}, nil
}

func (es *ExternalScaler) GetMetrics(ctx context.Context, metricRequest *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	log.Println("Executing method GetMetrics")

	if metricRequest.MetricName != "external-scaler-e2e-test" {
		return nil, fmt.Errorf("invalid metric name - %s", metricRequest.MetricName)
	}

	return &pb.GetMetricsResponse{
		MetricValues: []*pb.MetricValue{
			{MetricName: "external-scaler-e2e-test", MetricValue: ExternalScalerValue},
		},
	}, nil
}

func (es *ExternalScaler) StreamIsActive(scaledObjectRef *pb.ScaledObjectRef, epsServer pb.ExternalScaler_StreamIsActiveServer) error {
	log.Println("Executing method StreamIsActive")

	return nil
}

func main() {
	go RunManagementApi()

	grpcServer := grpc.NewServer()
	lis, _ := net.Listen("tcp", ":6000")

	pb.RegisterExternalScalerServer(grpcServer, &ExternalScaler{})
	log.Println("Listening external scaler on :6000")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
