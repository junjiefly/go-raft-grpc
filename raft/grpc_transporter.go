package raft

import (
	"errors"
	"github.com/junjiefly/jlog"
	"github.com/junjiefly/go-raft-grpc/raft/client"
	protobuf "github.com/junjiefly/go-raft-grpc/raft/protobuf"
	"sync"
	"time"
)

type GRPCTransporter struct {
	grpcClient map[string]*client.RaftGrpcClient
	sync.RWMutex
}

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

// Creates a new grpc transporter
func NewGRPCTransporter(timeout time.Duration, peers []string) *GRPCTransporter {
	t := &GRPCTransporter{
		grpcClient: make(map[string]*client.RaftGrpcClient),
	}
	for _, v := range peers {
		var err error
		t.grpcClient[v], err = client.NewRaftGrpcClient(v, "0.0.0.0", timeout)
		if err != nil {
			jlog.Fatalln("create grpc client for:", v, "err:", err)
		}
	}
	return t
}

func (t *GRPCTransporter) SendClusterJoinRequest(server Server, peer string, req *DefaultJoinCommand) error {
	var request = &protobuf.ClusterJoinRequest{
		Name:             req.Name,
		ConnectionString: req.ConnectionString,
	}
	traceln(server.Name(), "POST", "clusterJoin")
	reply, err := t.grpcClient[peer].SendClusterJoinRequest(request)
	if err != nil {
		traceln("transporter.ae.response.error:", err)
		return err
	}
	if reply.ErrMsg != "" {
		return errors.New(reply.ErrMsg)
	}
	return nil
}

func (t *GRPCTransporter) SendClusterLeaveRequest(server Server, peer string, req *DefaultLeaveCommand) error {
	var request = &protobuf.ClusterLeaveRequest{
		Name: req.Name,
	}
	traceln(server.Name(), "POST", "clusterLeave")
	reply, err := t.grpcClient[peer].SendClusterLeaveRequest(request)
	if err != nil {
		traceln("transporter.ae.response.error:", err)
		return err
	}
	if reply.ErrMsg != "" {
		return errors.New(reply.ErrMsg)
	}
	return nil
}

// Sends an AppendEntries RPC to a peer.
func (t *GRPCTransporter) SendAppendEntriesRequest(server Server, peer *Peer, req *AppendEntriesRequest) *AppendEntriesResponse {
	var request = &protobuf.AppendEntriesRequest{
		Term:         req.Term,
		PrevLogIndex: req.PrevLogIndex,
		PrevLogTerm:  req.PrevLogTerm,
		CommitIndex:  req.CommitIndex,
		LeaderName:   req.LeaderName,
		Entries:      req.Entries,
	}
	traceln(server.Name(), "POST", "appendEntries")
	reply, err := t.grpcClient[peer.Name].SendAppendEntriesRequest(request)
	if err != nil {
		traceln("transporter.ae.response.error:", err)
		return nil
	}
	resp := &AppendEntriesResponse{
		pb: reply,
	}
	return resp
}

// Sends a RequestVote RPC to a peer.
func (t *GRPCTransporter) SendVoteRequest(server Server, peer *Peer, req *RequestVoteRequest) (*RequestVoteResponse, error) {
	var request = &protobuf.RequestVoteRequest{
		Term:          req.Term,
		LastLogIndex:  req.LastLogIndex,
		LastLogTerm:   req.LastLogTerm,
		CandidateName: req.CandidateName,
	}
	reply, err := t.grpcClient[peer.Name].SendVoteRequest(request)
	if err != nil {
		traceln("transporter.ae.response.error:", err)
		return nil, err
	}
	resp := &RequestVoteResponse{
		Term:        reply.Term,
		VoteGranted: reply.VoteGranted,
	}
	return resp, nil
}

// Sends a RequestVote RPC to a peer to select a new leader.
func (t *GRPCTransporter) SendNewVoteRequest(server Server, peer string, req *NewVoteRequest) (*NewVoteResponse, error) {
	resp := &NewVoteResponse{}
	var request = &protobuf.NewVoteRequest{}
	traceln(server.Name(), "POST", "SendNewVoteRequest")
	reply, err := t.grpcClient[peer].SendNewVoteRequest(request)
	if err != nil {
		traceln("transporter.ae.response.error:", err)
		return nil, err
	}
	resp.Ok = reply.Ok
	return resp, nil
}

// Sends a SnapshotRequest RPC to a peer.
func (t *GRPCTransporter) SendSnapshotRequest(server Server, peer *Peer, req *SnapshotRequest) *SnapshotResponse {
	var request = &protobuf.SnapshotRequest{
		LeaderName: req.LeaderName,
		LastIndex:  req.LastIndex,
		LastTerm:   req.LastTerm,
	}
	traceln(server.Name(), "POST", "SendSnapshotRequest")
	reply, err := t.grpcClient[peer.Name].SendSnapshotRequest(request)
	if err != nil {
		traceln("transporter.ae.response.error:", err)
		return nil
	}
	resp := &SnapshotResponse{
		Success: reply.Success,
	}
	return resp
}

// Sends a SnapshotRequest RPC to a peer.
func (t *GRPCTransporter) SendSnapshotRecoveryRequest(server Server, peer *Peer, req *SnapshotRecoveryRequest) *SnapshotRecoveryResponse {
	var request = &protobuf.SnapshotRecoveryRequest{
		LeaderName: req.LeaderName,
		LastIndex:  req.LastIndex,
		LastTerm:   req.LastTerm,
		State:      req.State,
	}
	for k := range req.Peers {
		request.Peers = append(request.Peers, &protobuf.Peer{Name: req.Peers[k].Name, ConnectionString: req.Peers[k].ConnectionString})
	}
	traceln(server.Name(), "POST", "SendSnapshotRecoveryRequest")
	reply, err := t.grpcClient[peer.Name].SendSnapshotRecoveryRequest(request)
	if err != nil {
		traceln("transporter.ae.response.error:", err)
		return nil
	}
	resp := &SnapshotRecoveryResponse{
		Term:        reply.Term,
		Success:     reply.Success,
		CommitIndex: reply.CommitIndex,
	}
	return resp
}
