package mesh

import (
	"fmt"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
)

func makeRoutes(nodes []NodeInfo) []*route.Route {
	serviceRoutes := []*route.Route{}
	serviceIDs := map[string]struct{}{}
	for _, n := range nodes {
		if n.NodeType() == "gateway" {
			continue
		}
		if _, ok := serviceIDs[n.ID()]; ok {
			continue
		}
		serviceIDs[n.ID()] = struct{}{}
		// TODO: support multiple clusters for a service
		serviceRoutes = append(serviceRoutes, &route.Route{
			Name: n.ID(),
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: fmt.Sprintf("/%s/", n.ID()),
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: clusterID(n),
					},
					PrefixRewrite: "/",
				},
			},
		})
	}
	return append(serviceRoutes, DefaultRoutes()...)
}

func DefaultRoutes() []*route.Route {
	// TODO: route to service on node for sidecar configs
	return []*route.Route{
		{
			Name: "google",
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/google/",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: "www_google_com",
					},
					PrefixRewrite: "/",
					HostRewriteSpecifier: &route.RouteAction_HostRewrite{
						HostRewrite: "www.google.com",
					},
				},
			},
		},
		{
			Name: "health-check",
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Path{
					Path: "/hc",
				},
			},
			Action: &route.Route_DirectResponse{
				DirectResponse: &route.DirectResponseAction{
					Status: 200,
				},
			},
		},
		{
			Name: "default",
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/",
				},
			},
			Action: &route.Route_DirectResponse{
				DirectResponse: &route.DirectResponseAction{
					Status: 404,
				},
			},
		},
	}
}

func NewDefaultRoutes() []*v2.RouteConfiguration {
	return []*v2.RouteConfiguration{
		{
			Name: "default_route_config",
			VirtualHosts: []*route.VirtualHost{
				{
					Name:    "vhost1",
					Domains: []string{"*"},
					Routes:  DefaultRoutes(),
				},
			},
		},
	}
}

func NewRoutes(nodes []NodeInfo) []*v2.RouteConfiguration {
	return []*v2.RouteConfiguration{
		{
			Name: "route_config_1",
			VirtualHosts: []*route.VirtualHost{
				{
					Name:    "vhost1",
					Domains: []string{"*"},
					Routes:  makeRoutes(nodes),
				},
			},
		},
	}
}
