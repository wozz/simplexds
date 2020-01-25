package simplexds

import (
	"context"
	"net"
	"os"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	return logger
}

func newCallbacks(ctx context.Context, logger *logrus.Logger, snapshotCache cache.SnapshotCache) *cb {
	callbacks := &cb{
		l:     logger,
		cache: snapshotCache,
		store: &store{
			nodeMap: make(map[int64]*node),
		},
		updateChan: make(chan struct{}, 100),
	}
	go callbacks.run(ctx)
	return callbacks
}

func Run(ctx context.Context) {
	logger := newLogger()
	snapshotCache := cache.NewSnapshotCache(true, cache.IDHash{}, logger)
	callbacks := newCallbacks(ctx, logger, snapshotCache)
	server := xds.NewServer(ctx, snapshotCache, callbacks)

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		logger.Panicf("could not listen: %v", err)
	}

	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	api.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	api.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	api.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	api.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	go func() {
		logger.Println("Starting server")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Panicf("fatal error: %v", err)
		}
	}()
}
