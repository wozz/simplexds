package mesh

import (
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	http_connection_manager "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

func NewListener(port int, defaultRoutes bool) *v2.Listener {
	return &v2.Listener{
		Name: "default",
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: uint32(port),
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{
			{
				Filters: []*listener.Filter{
					{
						Name: "envoy.http_connection_manager",
						ConfigType: &listener.Filter_TypedConfig{
							TypedConfig: httpConnectionManager(defaultRoutes),
						},
					},
				},
			},
		},
	}
}

func routeConfigName(defaultRoutes bool) string {
	if defaultRoutes {
		return "default_route_config"
	}
	return "route_config_1"
}

func httpConnectionManager(defaultRoutes bool) *any.Any {
	hcm := &http_connection_manager.HttpConnectionManager{
		StatPrefix: "hcm1",
		RouteSpecifier: &http_connection_manager.HttpConnectionManager_Rds{
			Rds: &http_connection_manager.Rds{
				ConfigSource: &core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_Ads{
						Ads: &core.AggregatedConfigSource{},
					},
				},
				RouteConfigName: routeConfigName(defaultRoutes),
			},
		},
		HttpFilters: []*http_connection_manager.HttpFilter{
			{
				Name: "envoy.router",
			},
		},
	}
	hcmAny, _ := ptypes.MarshalAny(hcm)
	return hcmAny
}
