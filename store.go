package simplexds

import (
	"math/rand"
	"sync"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/wozz/simplexds/mesh"
)

func randMaxPriorityOrZero(max int) int {
	if rand.Intn(10) == 0 {
		return 0
	}
	return rand.Intn(max)
}

// emulate nodes moving between priorities randomly
func setRandomPrioritiesForNodesPerCluster(nodes []*node) (newNodeList []*node) {
	nodesPerCluster := make(map[string][]*node)
	for _, n := range nodes {
		nodesPerCluster[n.Cluster()] = append(nodesPerCluster[n.Cluster()], n)
	}

	for _, clusterNodes := range nodesPerCluster {
		// randomly set all nodes to priority zero to emulate panic mode where localities are unknown
		// or too many local nodes are unhealthy
		maxPriority := randMaxPriorityOrZero(len(clusterNodes))
		nextPriority := 0
		for _, nodeInCluster := range clusterNodes {
			if nextPriority > maxPriority {
				if maxPriority == 0 {
					nodeInCluster.priority = 0
				} else {
					nodeInCluster.priority = rand.Intn(maxPriority)
				}
			} else {
				nodeInCluster.priority = nextPriority
				nextPriority++
			}
		}
		newNodeList = append(newNodeList, clusterNodes...)
	}
	return
}

type node struct {
	nodeInfo  *core.Node
	ipAddress string
	priority  int
	streamID  int64
}

func (n *node) ID() string {
	return n.nodeInfo.GetId()
}

func (n *node) IP() string {
	return n.ipAddress
}

func (n *node) Cluster() string {
	return n.nodeInfo.GetCluster()
}

func (n *node) NodeType() string {
	return n.nodeInfo.GetMetadata().GetFields()["node_type"].GetStringValue()
}

func (n *node) PreferredPort() int {
	return int(n.nodeInfo.GetMetadata().GetFields()["preferred_port"].GetNumberValue())
}

func (n *node) Priority() int {
	return n.priority
}

type store struct {
	mu      sync.Mutex
	nodeMap map[int64]*node
}

func (s *store) shufflePriorities() {
	s.mu.Lock()
	defer s.mu.Unlock()
	var nodeList []*node
	for _, n := range s.nodeMap {
		nodeList = append(nodeList, n)
	}
	newNodeList := setRandomPrioritiesForNodesPerCluster(nodeList)
	for _, n := range newNodeList {
		s.nodeMap[n.streamID] = n
	}
}

func (s *store) allNodes() (nodeInfos []mesh.NodeInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, node := range s.nodeMap {
		nodeInfos = append(nodeInfos, node)
	}
	return
}

func (s *store) updateNode(id int64, nodeInfo *core.Node) (updated bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	node, ok := s.nodeMap[id]
	if !ok {
		return
	}
	if node.nodeInfo == nil {
		updated = true
	}
	node.nodeInfo = nodeInfo
	return
}

func (s *store) getGWNodeIDs() (ids []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, node := range s.nodeMap {
		if node == nil {
			continue
		}
		if node.NodeType() == "gateway" {
			ids = append(ids, node.ID())
		}
	}
	return
}

func (s *store) getNode(id int64) *node {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.nodeMap[id]
}

func (s *store) addNode(id int64, ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodeMap[id] = &node{
		ipAddress: ip,
		streamID:  id,
	}
}

func (s *store) removeNode(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.nodeMap, id)
}
