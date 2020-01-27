package network

// NODE ROLES MASTER AND SLAVE

// ServerNode comment
type ServerNode struct {
	Address       string
	Port          string
	IsLeader      bool
	MasterAddress string
}
