package simplexds

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/log"
	"github.com/wozz/simplexds/mesh"
	"google.golang.org/grpc/peer"
)

func (c *cb) update() {
	version := fmt.Sprintf("%d", time.Now().Unix())
	c.store.shufflePriorities()
	for _, node := range c.store.allNodes() {
		switch node.NodeType() {
		case "gateway":
			c.cache.SetSnapshot(node.ID(), mesh.GatewaySnapshot(version, node.PreferredPort(), c.store.allNodes()))
		case "service":
			c.cache.SetSnapshot(node.ID(), mesh.NodeSnapshot(version, node.PreferredPort()))
		}
	}
}

func (c *cb) run(ctx context.Context) {
	c.l.Infof("Initialize callbacks")
	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			c.l.Infof("Callbacks exiting")
			break
		case <-c.updateChan:
			c.update()
		case <-t.C:
			c.update()
		}
	}
}

type cb struct {
	l          log.Logger
	cache      cache.SnapshotCache
	store      *store
	updateChan chan struct{}
}

func (c *cb) OnStreamOpen(ctx context.Context, id int64, typeURL string) error {
	c.l.Debugf("OnStreamOpen[%d] %s", id, typeURL)
	p, ok := peer.FromContext(ctx)
	if !ok {
		return errors.New("failed to get peer info")
	}
	hostAddr, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return fmt.Errorf("could not parse peer address: %w", err)
	}
	c.store.addNode(id, hostAddr)
	c.l.Debugf("OnStreamOpen[%d]: %s", id, hostAddr)
	return nil
}

func (c *cb) OnStreamClosed(id int64) {
	c.l.Debugf("OnStreamClosed[%d]", id)
	c.store.removeNode(id)
}

func (c *cb) OnStreamRequest(id int64, req *envoy_api_v2.DiscoveryRequest) error {
	c.l.Debugf("OnStreamRequest[%d]", id)
	if c.store.updateNode(id, req.GetNode()) {
		c.updateChan <- struct{}{}
	}
	return nil
}

func (c *cb) OnStreamResponse(id int64, req *envoy_api_v2.DiscoveryRequest, resp *envoy_api_v2.DiscoveryResponse) {
	c.l.Debugf("OnStreamResponse[%d]", id)
}

func (c *cb) OnFetchRequest(ctx context.Context, req *envoy_api_v2.DiscoveryRequest) error {
	c.l.Debugf("OnFetchRequest")
	return errors.New("not implemented")
}

func (c *cb) OnFetchResponse(req *envoy_api_v2.DiscoveryRequest, resp *envoy_api_v2.DiscoveryResponse) {
	c.l.Debugf("OnFetchResponse")
}
