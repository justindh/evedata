package conservator

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func (s *Conservator) runRPC() error {
	l, err := net.Listen("tcp", "0.0.0.0:3001")
	if err != nil {
		return err
	}

	err = rpc.Register(s)
	if err != nil {
		return err
	}
	rpc.HandleHTTP()
	go http.Serve(l, nil)
	return nil
}

func (s *Conservator) VerifyDiscord(args *string, reply *bool) error {
	*reply = false
	g, err := s.discord.Guild(*args)
	if err != nil {
		log.Println(err)
		return nil
	}
	if g.Name != "" {
		*reply = true
	}
	return nil
}

func (s *Conservator) VerifyDiscordChannel(args []string, reply *bool) error {
	*reply = false
	g, err := s.discord.GuildChannels(args[0])
	if err != nil {
		return err
	}

	for _, c := range g {
		if c.ID == args[1] {
			*reply = true
		}
	}

	return nil
}

func (s *Conservator) GetChannels(botServiceID *int32, reply *[][]string) error {
	// Get the service
	si, ok := s.services.Load(*botServiceID)
	if !ok {
		return errors.New("missing botServiceID")
	}
	service := si.(Service)

	channels, err := service.Server.GetChannels()
	if err != nil {
		return err
	}

	for _, ch := range channels {
		*reply = append(*reply, []string{ch.ID, ch.Name})
	}

	return nil
}

func (s *Conservator) GetRoles(botServiceID *int32, reply *[][]string) error {
	// Get the service
	si, ok := s.services.Load(*botServiceID)
	if !ok {
		return errors.New("missing botServiceID")
	}
	service := si.(Service)

	roles, err := service.Server.GetRoles()
	if err != nil {
		return err
	}

	for _, ch := range roles {
		*reply = append(*reply, []string{ch.ID, ch.Name})
	}

	return nil
}
