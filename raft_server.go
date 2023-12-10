package go_raft_grpc

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/junjiefly/go-raft-grpc/raft"
	protobuf "github.com/junjiefly/go-raft-grpc/raft/protobuf"
	"google.golang.org/grpc"
	"net/url"
	"strings"
	"time"
)

/*
 	go-raft-grpc is based on go-raft[https://github.com/goraft/raft/],
	using grpc other than http as the transporter, nothing changed.
*/

type RaftServer struct {
	peers      []string
	isLeader   bool
	dataDir    string
	httpAddr   string
	raftServer raft.Server

	gRpc *grpc.Server
	protobuf.UnimplementedRaftServerServer
}

func NewRaftServer(peers []string, httpAddr string, dataDir string, pulseSeconds int) *RaftServer {
	var s = &RaftServer{
		peers:    peers,
		httpAddr: httpAddr,
		dataDir:  dataDir,
	}
	var err error
	transporter := raft.NewGRPCTransporter(5*time.Second, peers)
	s.raftServer, err = raft.NewServer(s.httpAddr, s.dataDir, transporter, nil, s, "")
	if err != nil {
		glog.V(0).Infoln(err)
		return nil
	}
	s.raftServer.SetHeartbeatInterval(1 * time.Second)
	s.raftServer.SetElectionTimeout(time.Duration(pulseSeconds) * 1150 * time.Millisecond)
	s.raftServer.Start()

	// Join to leader if specified.
	if len(s.peers) > 0 {
		if !s.raftServer.IsLogEmpty() {
			glog.V(0).Infoln("starting existing cluster")
		} else {
			glog.V(0).Infoln("joining cluster with peers:", strings.Join(s.peers, ","))
			firstJoinError := s.Join(s.peers)
			if firstJoinError != nil {
				glog.V(0).Infoln("no existing server found. starting as leader in the new cluster.")
				_, err := s.raftServer.Do(&raft.DefaultJoinCommand{
					Name:             s.raftServer.Name(),
					ConnectionString: "http://" + s.httpAddr,
				})
				if err != nil {
					glog.V(0).Infoln(err)
					return nil
				}
			}
		}
	} else if s.raftServer.IsLogEmpty() {
		glog.V(0).Infoln("initializing new cluster")
		_, err := s.raftServer.Do(&raft.DefaultJoinCommand{
			Name:             s.raftServer.Name(),
			ConnectionString: "http://" + s.httpAddr,
		})

		if err != nil {
			glog.V(0).Infoln(err)
			return nil
		}
	} else {
		glog.V(0).Infoln("recovered from log")
	}
	s.raftServer.AddEventListener(raft.LeaderChangeEventType, func(e raft.Event) {
		if s.raftServer.Leader() != "" {
			if s.IsLeader() {
				time.AfterFunc(10*time.Second, func() { s.isLeader = true })
			} else {
				s.isLeader = false
			}
			glog.V(0).Infoln(s.raftServer.Leader(), "becomes leader!")
		}
		glog.V(0).Infoln("leader changed!", "prev:", e.PrevValue(), "now:", e.Value())
	})
	if s.IsLeader() {
		glog.V(0).Infoln("i am the leader!")
		s.isLeader = true
	} else {
		if s.raftServer.Leader() != "" {
			glog.V(0).Infoln(s.raftServer.Leader(), "is the leader.")
		}
	}
	return s
}

//IsLeader  i am the leader?
func (s *RaftServer) IsLeader() bool {
	if s == nil || s.raftServer == nil {
		return false
	}
	return s.raftServer.Leader() == s.raftServer.Name()
}

//Leader  cluster leader
func (s *RaftServer) Leader() string {
	if s == nil || s.raftServer == nil {
		return ""
	}
	return s.raftServer.Leader()
}

//Peers my partner nodes
func (s *RaftServer) Peers() (members []string) {
	if s == nil || s.raftServer == nil {
		return nil
	}
	peers := s.raftServer.Peers()
	for _, p := range peers {
		members = append(members, p.Name)
	}
	return
}

//Nodes  all nodes
func (s *RaftServer) Nodes() (members []string) {
	if s == nil || s.raftServer == nil {
		return nil
	}
	peers := s.raftServer.Peers()
	for _, p := range peers {
		members = append(members, p.Name)
	}
	members = append(members, s.httpAddr)
	return
}

// Join
//joins an existing cluster.
func (s *RaftServer) Join(peers []string) error {
	command := &raft.DefaultJoinCommand{
		Name:             s.raftServer.Name(),
		ConnectionString: "http://" + s.httpAddr,
	}
	var err error
	for _, m := range peers {
		if m == s.httpAddr {
			continue
		}
		err = s.raftServer.Transporter().SendClusterJoinRequest(s.raftServer, strings.TrimSpace(m), command)
		glog.V(0).Infoln("try to connect to:", m)
		if err != nil {
			glog.V(0).Infoln("connect to peer:", m, "err:", err.Error())
			if _, ok := err.(*url.Error); ok {
				continue
			}
		} else {
			return nil
		}
	}
	return errors.New("could not connect to any peers")
}

//DeathNotify
//If a node is about to die, notify other nodes to elect a new leader as soon as possible before the node dies
func (s *RaftServer) DeathNotify() {
	members := s.Peers()
	s.raftServer.Stop()
	for _, v := range members {
		reply, err := s.raftServer.Transporter().SendNewVoteRequest(s.raftServer, v, nil)
		if err != nil {
			glog.V(0).Infoln("send new vote request to:", v, "err:", err)
		} else if reply.Ok == true {
			break
		}
	}
	//wait for a second to elect the new leader.
	time.Sleep(time.Millisecond * time.Duration(1600))
	return
}

//StartRaftServer
//start a raft server
func StartRaftServer(ip, peerIps, port, dataDir string) *RaftServer {
	var peers []string
	if peerIps != "" {
		peers = strings.Split(peerIps, ",")
		for k := range peers {
			peers[k] = fmt.Sprintf("%s:%s", peers[k], port)
		}
	}
	raftAddr := fmt.Sprintf("%s:%s", ip, port)
	raftServer := NewRaftServer(peers, raftAddr, dataDir, 3)
	if raftServer == nil {
		glog.Fatalln("fail to create raft server:", raftAddr)
	}
	return raftServer
}
