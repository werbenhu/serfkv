package cluster

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/hashicorp/logutils"
	"github.com/hashicorp/serf/serf"
	"github.com/rs/xid"
)

type Options struct {
	// Unique ID of server
	ID string
	// Local address to bind to
	Address string
	// Members in the cluster
	Members []string
}

type Server struct {
	Opts    *Options
	events  chan serf.Event
	serf    *serf.Serf
	handler *Handler
	storage sync.Map
}

func New(opts *Options) (*Server, error) {

	var err error
	var port string
	config := serf.DefaultConfig()

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("ERROR"),
		Writer:   os.Stderr,
	}
	logger := log.New(os.Stderr, "", log.LstdFlags)
	logger.SetOutput(filter)

	config.Logger = logger
	config.MemberlistConfig.Logger = logger
	config.MemberlistConfig.BindAddr, port, err = net.SplitHostPort(opts.Address)
	config.MemberlistConfig.BindPort, _ = strconv.Atoi(port)

	if err != nil {
		config.MemberlistConfig.BindPort = 0
	}

	s := &Server{
		Opts:   opts,
		events: make(chan serf.Event),
	}
	if len(opts.ID) == 0 {
		hostname, _ := os.Hostname()
		opts.ID = hostname + "-" + xid.New().String()
	}

	config.EventCh = s.events
	config.NodeName = opts.ID
	s.handler = NewHandler(s)
	s.serf, err = serf.Create(config)

	logger.Println("================================")
	logger.Printf("start serf cluster on %s\n", opts.Address)
	logger.Println("================================")

	go s.eventLoop()

	if len(opts.Members) > 0 {
		s.Join(opts.Members)
	}

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) Delete(key string, local bool) error {
	s.storage.Delete(key)
	if local {
		payload, err := (&Message{
			Action: "del",
			Key:    key,
		}).Encode()

		if err != nil {
			return errors.New("encode failed")
		}
		s.Broadcast(payload)
	}
	return nil
}

func (s *Server) Get(key string) (any, error) {
	if val, ok := s.storage.Load(key); ok {
		return val, nil
	}
	return "", errors.New("not found")
}

func (s *Server) Set(key string, val any, local bool) error {
	s.storage.Store(key, val)
	if local {
		payload, err := (&Message{
			Action: "set",
			Key:    key,
			Val:    val,
		}).Encode()

		if err != nil {
			return errors.New("encode failed")
		}
		s.Broadcast(payload)
	}
	return nil
}

func (s *Server) Join(members []string) error {
	_, err := s.serf.Join(members, true)
	return err
}

func (s *Server) Broadcast(payload []byte) {
	s.serf.UserEvent(s.Opts.ID, payload, false)
}

func (s *Server) eventLoop() {
	for e := range s.events {
		switch e.EventType() {
		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				if s.serf.LocalMember().Name == member.Name {
					continue
				}
				s.handler.HandleJoin(member)
			}

		case serf.EventMemberUpdate:
			for _, member := range e.(serf.MemberEvent).Members {
				if s.serf.LocalMember().Name == member.Name {
					continue
				}
				s.handler.HandleUpdate(member)
			}

		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, member := range e.(serf.MemberEvent).Members {
				if s.serf.LocalMember().Name == member.Name {
					continue
				}
				s.handler.HandleLeave(member)
			}

		case serf.EventMemberReap:
			for _, member := range e.(serf.MemberEvent).Members {
				if s.serf.LocalMember().Name == member.Name {
					continue
				}
				s.handler.HandleLeave(member)
			}

		case serf.EventUser:
			msg := e.(serf.UserEvent)
			if msg.Name == s.serf.LocalMember().Name {
				continue
			}
			s.handler.HandleMessage(msg.Payload)

		case serf.EventQuery:
			query := e.(*serf.Query)
			if query.SourceNode() == s.serf.LocalMember().Name {
				continue
			}
			s.handler.HandleQuery(query.Payload)
		}
	}
}
