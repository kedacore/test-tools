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
	fmt.Printf("Running http management server on port: %d\n", 8080)
	fmt.Print("example usage:\n - POST -> localhost:8080/api/value/3\n - POST -> localhost:8080/api/floatvalue/3.14\n")
	fmt.Print("if you set one of those values as negative, it won't be sent in the payload\n")
	http.ListenAndServe(":8080", nil)
}

var ExternalScalerValue int64 = 0
var ExternalScalerValueFloat = .0

type ExternalScaler struct {
	pb.UnimplementedExternalScalerServer

	ctx context.Context
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
	for {
		tmr := time.NewTimer(time.Second)
		select {
		case <-es.ctx.Done():
			tmr.Stop()
			return nil
		case <-tmr.C:
			tmr.Stop()
			epsServer.Send(&pb.IsActiveResponse{Result: ExternalScalerValue > 0})
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go RunManagementApi()

	grpcServer := grpc.NewServer()
	lis, _ := net.Listen("tcp", ":6000")

	pb.RegisterExternalScalerServer(grpcServer, &ExternalScaler{ctx: ctx})
	log.Println("Listening external scaler on :6000")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
