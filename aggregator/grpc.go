package main

import "github.com/boyanivskyy/toll-calculator/types"

// NOTE: aggregator service for GRPC, maybe it should be renamed to some more meaningful file name

type GRPSAggregatorServer struct {
	types.UnimplementedAggregatorServer
	service Aggregator
}

func NewGRPSServer(service Aggregator) *GRPSAggregatorServer {
	return &GRPSAggregatorServer{
		service: service,
	}

}

func (s *GRPSAggregatorServer) AggregateDistance(req types.AggregateRequest) error {
	distance := types.Distance{
		OBUID: int(req.ObuId),
		Value: req.Value,
		Unix:  req.Unix,
	}

	return s.service.AggregateDistance(distance)
}
