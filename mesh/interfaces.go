package mesh

type NodeInfo interface {
	IP() string
	ID() string
	Cluster() string
	NodeType() string
	PreferredPort() int
	Priority() int
}
