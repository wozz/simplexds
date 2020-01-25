package mesh

import "github.com/envoyproxy/go-control-plane/pkg/cache"

func GatewaySnapshot(id string, port int, nodes []NodeInfo) cache.Snapshot {
	return cache.NewSnapshot(
		id,
		toResources(endpointList(NewNodeEndpoints(nodes))),
		toResources(clusterList(NewGatewayClusters(nodes))),
		toResources(routeList(NewRoutes(nodes))),
		[]cache.Resource{NewListener(port, false)},
		nil,
	)
}
