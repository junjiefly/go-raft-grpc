package raft

import (
	protobuf "github.com/junjiefly/go-raft-grpc/raft/protobuf"
)

// The request sent to a server to append entries to the log.
type AppendEntriesRequest struct {
	Term         uint64
	PrevLogIndex uint64
	PrevLogTerm  uint64
	CommitIndex  uint64
	LeaderName   string
	Entries      []*protobuf.LogEntry
}

// The response returned from a server appending entries to the log.
type AppendEntriesResponse struct {
	pb     *protobuf.AppendEntriesResponse
	peer   string
	append bool
}

// Creates a new AppendEntries request.
func newAppendEntriesRequest(term uint64, prevLogIndex uint64, prevLogTerm uint64,
	commitIndex uint64, leaderName string, entries []*LogEntry) *AppendEntriesRequest {
	pbEntries := make([]*protobuf.LogEntry, len(entries))

	for i := range entries {
		pbEntries[i] = entries[i].pb
	}

	return &AppendEntriesRequest{
		Term:         term,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		CommitIndex:  commitIndex,
		LeaderName:   leaderName,
		Entries:      pbEntries,
	}
}

// Creates a new AppendEntries response.
func newAppendEntriesResponse(term uint64, success bool, index uint64, commitIndex uint64) *AppendEntriesResponse {
	pb := &protobuf.AppendEntriesResponse{
		Term:        term,
		Index:       index,
		Success:     success,
		CommitIndex: commitIndex,
	}

	return &AppendEntriesResponse{
		pb: pb,
	}
}

func (aer *AppendEntriesResponse) Index() uint64 {
	return aer.pb.GetIndex()
}

func (aer *AppendEntriesResponse) CommitIndex() uint64 {
	return aer.pb.GetCommitIndex()
}

func (aer *AppendEntriesResponse) Term() uint64 {
	return aer.pb.GetTerm()
}

func (aer *AppendEntriesResponse) Success() bool {
	return aer.pb.GetSuccess()
}
