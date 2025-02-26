package client

import (
	"context"
	raft "github.com/junjiefly/go-raft-grpc/raft/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/keepalive"
	"net"
	"errors"
	"time"
)

var KeepaliveTime = 10 * time.Second
var PermitWithoutStream = true

var defaultRaftTimeout = time.Second * 2

type RaftGrpcClient struct {
	Client         raft.RaftServerClient
	Conn           *grpc.ClientConn
	Addr           string
	RequestTimeout time.Duration
}

func NewGrpcClient(addr, lAddr string) (*grpc.ClientConn, error) {
	connParams := grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  1.0 * time.Second,
			Multiplier: 1.6,
			Jitter:     0.2,
			MaxDelay:   3 * time.Second,
		},
		MinConnectTimeout: 3 * time.Second,
	}
	keepaliveParams := keepalive.ClientParameters{
		Time:                KeepaliveTime,
		Timeout:             3 * time.Second,
		PermitWithoutStream: PermitWithoutStream}
	localAddr, _ := net.ResolveTCPAddr("tcp", lAddr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(),
		grpc.WithConnectParams(connParams), grpc.WithKeepaliveParams(keepaliveParams),
		grpc.WithReadBufferSize(1<<20), grpc.WithWriteBufferSize(1<<20),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(16<<20)),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			dial := &net.Dialer{LocalAddr: localAddr}
			return dial.DialContext(ctx, "tcp", addr)
		}))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewRaftGrpcClient(addr, lAddr string, timeout time.Duration) (*RaftGrpcClient, error) {
	conn, err := NewGrpcClient(addr, lAddr)
	if err != nil {
		return nil, err
	}
	cli := raft.NewRaftServerClient(conn)
	client := &RaftGrpcClient{
		Client:         cli,
		Conn:           conn,
		Addr:           addr,
		RequestTimeout: timeout,
	}
	if client.RequestTimeout == 0 {
		client.RequestTimeout = defaultRaftTimeout
	}
	return client, nil
}

func (grpcClient *RaftGrpcClient) Close() {
	grpcClient.Conn.Close()
}

func (grpcClient *RaftGrpcClient) SendClusterJoinRequest(request *raft.ClusterJoinRequest) (*raft.ClusterJoinResponse, error) {
	if grpcClient == nil || grpcClient.Client == nil {
		return nil, errors.New("dest grpc client null")
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcClient.RequestTimeout)
	reply, err := grpcClient.Client.ClusterJoinHandler(ctx, request)
	cancel()
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (grpcClient *RaftGrpcClient) SendClusterLeaveRequest(request *raft.ClusterLeaveRequest) (*raft.ClusterLeaveResponse, error) {
	if grpcClient == nil || grpcClient.Client == nil {
		return nil, errors.New("dest grpc client null")
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcClient.RequestTimeout)
	reply, err := grpcClient.Client.ClusterLeaveHandler(ctx, request)
	cancel()
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (grpcClient *RaftGrpcClient) SendClusterStatusRequest(request *raft.ClusterStatusRequest) (*raft.ClusterStatusResponse, error) {
	if grpcClient == nil || grpcClient.Client == nil {
		return nil, errors.New("dest grpc client null")
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcClient.RequestTimeout)
	reply, err := grpcClient.Client.ClusterStatusHandler(ctx, request)
	cancel()
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (grpcClient *RaftGrpcClient) SendAppendEntriesRequest(request *raft.AppendEntriesRequest) (*raft.AppendEntriesResponse, error) {
	if grpcClient == nil || grpcClient.Client == nil {
		return nil, errors.New("dest grpc client null")
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcClient.RequestTimeout)
	reply, err := grpcClient.Client.AppendEntriesRequestHandler(ctx, request)
	cancel()
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (grpcClient *RaftGrpcClient) SendVoteRequest(request *raft.RequestVoteRequest) (*raft.RequestVoteResponse, error) {
	if grpcClient == nil || grpcClient.Client == nil {
		return nil, errors.New("dest grpc client null")
	}			       
	ctx, cancel := context.WithTimeout(context.Background(), grpcClient.RequestTimeout)
	reply, err := grpcClient.Client.VoteRequestHandler(ctx, request)
	cancel()
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (grpcClient *RaftGrpcClient) SendNewVoteRequest(request *raft.NewVoteRequest) (*raft.NewVoteResponse, error) {
	if grpcClient == nil || grpcClient.Client == nil {
		return nil, errors.New("dest grpc client null")
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcClient.RequestTimeout)
	reply, err := grpcClient.Client.NewVoteHandler(ctx, request)
	cancel()
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (grpcClient *RaftGrpcClient) SendSnapshotRequest(request *raft.SnapshotRequest) (*raft.SnapshotResponse, error) {
	if grpcClient == nil || grpcClient.Client == nil {
		return nil, errors.New("dest grpc client null")
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcClient.RequestTimeout)
	reply, err := grpcClient.Client.SnapshotRequestHandler(ctx, request)
	cancel()
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (grpcClient *RaftGrpcClient) SendSnapshotRecoveryRequest(request *raft.SnapshotRecoveryRequest) (*raft.SnapshotRecoveryResponse, error) {
	if grpcClient == nil || grpcClient.Client == nil {
		return nil, errors.New("dest grpc client null")
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcClient.RequestTimeout)
	reply, err := grpcClient.Client.SnapshotRecoveryRequestHandler(ctx, request)
	cancel()
	if err != nil {
		return reply, err
	}
	return reply, nil
}
