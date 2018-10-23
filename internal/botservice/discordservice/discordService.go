package discordservice

import (
	"github.com/antihax/evedata/internal/botservice"
	"github.com/bwmarrin/discordgo"
)

// BotService provides access to a discord session
// Discordgo handles rate throttling
type DiscordService struct {
	session  *discordgo.Session
	serverID string
}

// NewDiscordService sends a message to a discord channel ID
func NewDiscordService(session *discordgo.Session, serverID string) DiscordService {
	return DiscordService{session, serverID}
}

// SendMessageToChannel sends a message to a discord channel ID
func (c DiscordService) SendMessageToChannel(channel, message string) error {
	_, err := c.session.ChannelMessageSend(channel, message)
	return err
}

// SendMessageToUser sends a message to a discord user ID
func (c DiscordService) SendMessageToUser(user, message string) error {
	_, err := c.session.ChannelMessageSend(user, message)
	return err
}

// KickUser kicks a discord user ID
func (c DiscordService) KickUser(user, message string) error {
	_, err := c.session.ChannelMessageSend(user, message)
	return err
}

// Get the server name
func (c DiscordService) GetName() (string, error) {
	g, err := c.session.Guild(c.serverID)
	if err != nil {
		return "", err
	}
	return g.Name, nil
}

// Get Channels
func (c DiscordService) GetChannels() ([]botservice.Name, error) {
	g, err := c.session.GuildChannels(c.serverID)
	if err != nil {
		return nil, err
	}

	channels := []botservice.Name{}
	for _, ch := range g {
		channels = append(channels, botservice.Name{ID: ch.ID, Name: ch.Name})
	}

	return channels, nil
}

// GetRoles gets all the available roles to be assigned
func (c DiscordService) GetRoles() ([]botservice.Name, error) {
	g, err := c.session.GuildRoles(c.serverID)
	if err != nil {
		return nil, err
	}

	roles := []botservice.Name{}
	for _, role := range g {
		if !role.Managed && role.Name != "@everyone" {
			roles = append(roles, botservice.Name{ID: role.ID, Name: role.Name})
		}
	}

	return roles, nil
}
