package main

import (
	"context"

	"github.com/boyanivskyy/toll-calculator/types"
)

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

func (s *GRPSAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.OBUID),
		Value: req.Value,
		Unix:  req.Unix,
	}

	return &types.None{}, s.service.AggregateDistance(distance)
}
