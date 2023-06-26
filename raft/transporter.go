package raft

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// Transporter is the interface for allowing the host application to transport
// requests to other nodes.
type Transporter interface {
	SendClusterJoinRequest(server Server, node string, req *DefaultJoinCommand) error
	SendClusterLeaveRequest(server Server, node string, req *DefaultLeaveCommand) error
	SendVoteRequest(server Server, peer *Peer, req *RequestVoteRequest) (*RequestVoteResponse, error)
	SendNewVoteRequest(server Server, peer string, req *NewVoteRequest) (*NewVoteResponse, error)
	SendAppendEntriesRequest(server Server, peer *Peer, req *AppendEntriesRequest) *AppendEntriesResponse
	SendSnapshotRequest(server Server, peer *Peer, req *SnapshotRequest) *SnapshotResponse
	SendSnapshotRecoveryRequest(server Server, peer *Peer, req *SnapshotRecoveryRequest) *SnapshotRecoveryResponse
}
