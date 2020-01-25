package mesh

import (
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
)

func endpointForNode(node NodeInfo) *endpoint.LocalityLbEndpoints {
	return &endpoint.LocalityLbEndpoints{
		Priority: uint32(node.Priority()),
		LbEndpoints: []*endpoint.LbEndpoint{
			{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  node.IP(),
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: uint32(node.PreferredPort()),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func NewNodeEndpoints(nodes []NodeInfo) (endpoints []*v2.ClusterLoadAssignment) {
	endpointsPerCluster := map[string][]*endpoint.LocalityLbEndpoints{}
	for _, node := range nodes {
		if node.NodeType() == "gateway" {
			continue
		}
		endpointsPerCluster[clusterID(node)] = append(endpointsPerCluster[clusterID(node)], endpointForNode(node))
	}
	for cluster, endpointList := range endpointsPerCluster {
		endpoints = append(endpoints, &v2.ClusterLoadAssignment{
			ClusterName: cluster,
			Endpoints:   endpointList,
		})
	}
	return
}

func NewGoogleEndpoint() *v2.ClusterLoadAssignment {
	return &v2.ClusterLoadAssignment{
		ClusterName: "www_google_com",
		Endpoints: []*endpoint.LocalityLbEndpoints{
			{
				LbEndpoints: []*endpoint.LbEndpoint{
					{
						HostIdentifier: &endpoint.LbEndpoint_Endpoint{
							Endpoint: &endpoint.Endpoint{
								Address: &core.Address{
									Address: &core.Address_SocketAddress{
										SocketAddress: &core.SocketAddress{
											Protocol: core.SocketAddress_TCP,
											Address:  "www.google.com",
											PortSpecifier: &core.SocketAddress_PortValue{
												PortValue: 443,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
