package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	pb "github.com/kedacore/test-tools/external-scaler-e2e/externalscaler"
)

type ExternalScaler struct {
	pb.UnimplementedExternalScalerServer
}

func (es *ExternalScaler) IsActive(ctx context.Context, scaledObjectRef *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {
	log.Println("Executing method IsActive")

	metricValue, err := strconv.ParseInt(scaledObjectRef.ScalerMetadata["metricValue"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value for metric value - %s", err)
	}

	return &pb.IsActiveResponse{Result: metricValue > 0}, nil
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

	metricValue, err := strconv.ParseInt(metricRequest.ScaledObjectRef.ScalerMetadata["metricValue"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value for metric value - %s", err)
	}

	return &pb.GetMetricsResponse{
		MetricValues: []*pb.MetricValue{
			{MetricName: "external-scaler-e2e-test", MetricValue: metricValue},
		},
	}, nil
}

func (es *ExternalScaler) StreamIsActive(scaledObjectRef *pb.ScaledObjectRef, epsServer pb.ExternalScaler_StreamIsActiveServer) error {
	log.Println("Executing method StreamIsActive")

	return nil
}

func main() {
	grpcServer := grpc.NewServer()
	lis, _ := net.Listen("tcp", ":6000")

	pb.RegisterExternalScalerServer(grpcServer, &ExternalScaler{})

	log.Println("Listening on :6000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
