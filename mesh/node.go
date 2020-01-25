package mesh

import "github.com/envoyproxy/go-control-plane/pkg/cache"

func NodeSnapshot(id string, port int) cache.Snapshot {
	return cache.NewSnapshot(
		id,
		nil,
		toResources(clusterList(NewDefaultClusters())),
		toResources(routeList(NewDefaultRoutes())),
		[]cache.Resource{NewListener(port, true)},
		nil,
	)
}
