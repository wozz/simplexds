package mesh

import (
	envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

type convertResource func(cache.Resource) cache.Resource

type convertableToResource interface {
	each(handler convertResource) []cache.Resource
}

type clusterList []*envoy_api_v2.Cluster

func (cl clusterList) each(handler convertResource) (resources []cache.Resource) {
	for _, c := range cl {
		resources = append(resources, handler(c))
	}
	return
}

type routeList []*envoy_api_v2.RouteConfiguration

func (rl routeList) each(handler convertResource) (resources []cache.Resource) {
	for _, r := range rl {
		resources = append(resources, handler(r))
	}
	return
}

type endpointList []*envoy_api_v2.ClusterLoadAssignment

func (el endpointList) each(handler convertResource) (resources []cache.Resource) {
	for _, e := range el {
		resources = append(resources, handler(e))
	}
	return
}

func toResources(i convertableToResource) []cache.Resource {
	return i.each(func(r cache.Resource) cache.Resource {
		return r
	})
}
