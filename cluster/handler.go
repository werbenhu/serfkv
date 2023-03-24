package cluster

import (
	"fmt"

	"github.com/hashicorp/serf/serf"
)

type Handler struct {
	server *Server
}

func NewHandler(s *Server) *Handler {
	return &Handler{
		server: s,
	}
}

func (h *Handler) HandleJoin(member serf.Member) error {
	fmt.Println("A node has joined: " + member.Name)
	return nil
}

func (h *Handler) HandleUpdate(member serf.Member) error {
	fmt.Println("A node was updated: " + member.Name)
	return nil
}

func (s *Handler) HandleLeave(member serf.Member) error {
	fmt.Println("A node has left: " + member.Name)
	return nil
}

func (s *Handler) HandleMessage(payload []byte) error {
	fmt.Printf("HandleMessage payload:%s\n", string(payload))
	m := &Message{}
	if err := m.Decode(payload); err != nil {
		return err
	}

	switch m.Action {
	case "set":
		s.server.Set(m.Key, m.Val, false)
	case "del":
		s.server.Delete(m.Key, false)
	}

	return nil
}

func (s *Handler) HandleQuery(payload []byte) error {
	return nil
}
