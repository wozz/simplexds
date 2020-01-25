package mesh

import (
	"fmt"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func clusterID(node NodeInfo) string {
	return fmt.Sprintf("%s##%s", node.ID(), node.Cluster())
}

func googleTLS() *any.Any {
	upstreamContext := &auth.UpstreamTlsContext{
		Sni: "www.google.com",
	}
	upstreamContextAny, _ := ptypes.MarshalAny(upstreamContext)
	return upstreamContextAny
}

func NewGatewayClusters(nodes []NodeInfo) []*v2.Cluster {
	return append(NewClusters(nodes), NewDefaultClusters()...)
}

func NewClusters(nodes []NodeInfo) (clusters []*v2.Cluster) {
	clusterIDs := make(map[string]struct{})
	for _, node := range nodes {
		if node.NodeType() == "gateway" {
			continue
		}
		if _, ok := clusterIDs[clusterID(node)]; ok {
			continue
		}
		clusterIDs[clusterID(node)] = struct{}{}
		clusters = append(clusters, &v2.Cluster{
			Name: clusterID(node),
			ConnectTimeout: &duration.Duration{
				Seconds: 1,
			},
			ClusterDiscoveryType: &v2.Cluster_Type{
				Type: v2.Cluster_EDS,
			},
			LbPolicy: v2.Cluster_ROUND_ROBIN,
			EdsClusterConfig: &v2.Cluster_EdsClusterConfig{
				EdsConfig: &core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_Ads{
						Ads: &core.AggregatedConfigSource{},
					},
				},
			},
			UpstreamConnectionOptions: &v2.UpstreamConnectionOptions{
				TcpKeepalive: &core.TcpKeepalive{
					KeepaliveProbes: &wrappers.UInt32Value{
						Value: 3,
					},
					KeepaliveTime: &wrappers.UInt32Value{
						Value: 300,
					},
					KeepaliveInterval: &wrappers.UInt32Value{
						Value: 30,
					},
				},
			},
			HealthChecks: []*core.HealthCheck{
				{
					Timeout: &duration.Duration{
						Seconds: 3,
					},
					Interval: &duration.Duration{
						Seconds: 10,
					},
					InitialJitter: &duration.Duration{
						Seconds: 1,
					},
					IntervalJitter: &duration.Duration{
						Seconds: 1,
					},
					UnhealthyThreshold: &wrappers.UInt32Value{
						Value: 1,
					},
					HealthyThreshold: &wrappers.UInt32Value{
						Value: 2,
					},
					ReuseConnection: &wrappers.BoolValue{
						Value: true,
					},
					NoTrafficInterval: &duration.Duration{
						Seconds: 7,
					},
					UnhealthyInterval: &duration.Duration{
						Seconds: 5,
					},
					UnhealthyEdgeInterval: &duration.Duration{
						Seconds: 3,
					},
					HealthChecker: &core.HealthCheck_HttpHealthCheck_{
						HttpHealthCheck: &core.HealthCheck_HttpHealthCheck{
							Path:            "/hc",
							CodecClientType: envoy_type.CodecClientType_HTTP2,
						},
					},
				},
			},
		})
	}
	return
}

func NewDefaultClusters() []*v2.Cluster {
	return []*v2.Cluster{
		{
			Name: "www_google_com",
			ConnectTimeout: &duration.Duration{
				Seconds: 1,
			},
			ClusterDiscoveryType: &v2.Cluster_Type{
				Type: v2.Cluster_LOGICAL_DNS,
			},
			LbPolicy:        v2.Cluster_ROUND_ROBIN,
			DnsLookupFamily: v2.Cluster_V4_ONLY,
			LoadAssignment:  NewGoogleEndpoint(),
			TransportSocket: &core.TransportSocket{
				Name: "envoy.transport_sockets.tls",
				ConfigType: &core.TransportSocket_TypedConfig{
					TypedConfig: googleTLS(),
				},
			},
		},
	}
}
