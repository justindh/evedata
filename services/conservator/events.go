package conservator

import (
	"database/sql"
	"log"
	"strings"
)

func (s *Conservator) checkAllUsers() {
	s.services.Range(func(ki, vi interface{}) bool {
		service := vi.(Service)
		members, err := service.Server.GetMembers()
		if err != nil {
			return false
		}
		for _, m := range members {
			if err := s.checkUser(m.ID, m.Name, service.BotServiceID, m.Roles); err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func (c *Conservator) handleNewMember(memberID, memberName, serverID string) {

}

func (c *Conservator) handleMessage(memberID, memberName, serverID string) {

}

func (c *Conservator) checkUser(memberID, memberName string, botServiceID int32, roles []string) error {
	server, err := c.getService(botServiceID)
	if err != nil {
		return err
	}
	if inSlice("auth", strings.Split(server.Services, ",")) {
		if server.Options.Auth.Members != "" {
			if characterName, err := c.getMemberStatus(memberID, server.EntityID); err != nil {
				return err
			} else if characterName != "" { // Found them
				server.checkAddRoles(memberID, server.Options.Auth.Members, roles)
				return nil
			} else {
				server.checkRemoveRoles(memberID, server.Options.Auth.Members, roles)
			}
		}

		if server.Options.Auth.PlusTen != "" {
			if characterName, err := c.getPlusTenStatus(memberID, server.EntityID); err != nil {
				return err
			} else if characterName != "" { // Found them
				server.checkAddRoles(memberID, server.Options.Auth.PlusTen, roles)
				return nil
			} else {
				server.checkRemoveRoles(memberID, server.Options.Auth.PlusTen, roles)
			}
		}

		if server.Options.Auth.PlusFive != "" {
			if characterName, err := c.getPlusFiveStatus(memberID, server.EntityID); err != nil {
				return err
			} else if characterName != "" { // Found them
				server.checkAddRoles(memberID, server.Options.Auth.PlusFive, roles)
				return nil
			} else {
				server.checkRemoveRoles(memberID, server.Options.Auth.PlusFive, roles)
			}
		}

		if server.Options.Auth.Militia != "" && server.FactionID > 0 {
			if characterName, err := c.getMilitiaStatus(memberID, server.FactionID); err != nil {
				return err
			} else if characterName != "" { // Found them
				server.checkAddRoles(memberID, server.Options.Auth.Militia, roles)
				return nil
			} else {
				server.checkRemoveRoles(memberID, server.Options.Auth.Militia, roles)
			}
		}

		if server.Options.Auth.AlliedMilitia != "" && server.FactionID > 0 {
			if characterName, err := c.getMilitiaStatus(memberID, FactionAllies[server.FactionID]); err != nil {
				return err
			} else if characterName != "" { // Found them
				server.checkAddRoles(memberID, server.Options.Auth.AlliedMilitia, roles)
				return nil
			} else {
				server.checkRemoveRoles(memberID, server.Options.Auth.AlliedMilitia, roles)
			}
		}
	}
	return nil
}

func (s *Conservator) getMemberStatus(memberID string, entity int32) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT characterName
			FROM evedata.botCharacters C
			INNER JOIN evedata.crestTokens T ON T.characterID = C.characterID
			WHERE T.authCharacter = 1 AND botUserID = ? AND (allianceID = ? OR corporationID = ?) LIMIT 1;`, memberID, entity, entity).Scan(&ref); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return ref, nil
}

func (s *Conservator) getPlusFiveStatus(memberID string, entity int32) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT characterName
			FROM evedata.botCharacters C
			INNER JOIN evedata.crestTokens T ON T.characterID = C.characterID
			INNER JOIN evedata.entityContacts E ON E.contactID = T.allianceID OR E.contactID = T.corporationID OR E.contactID = T.tokenCharacterID
			WHERE T.authCharacter = 1 AND botUserID = ? AND entityID = ? AND standing = 10 LIMIT 1;`, memberID, entity).Scan(&ref); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return ref, nil
}

func (s *Conservator) getPlusTenStatus(memberID string, entity int32) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT characterName
			FROM evedata.botCharacters C
			INNER JOIN evedata.crestTokens T ON T.characterID = C.characterID
			INNER JOIN evedata.entityContacts E ON E.contactID = T.allianceID OR E.contactID = T.corporationID OR E.contactID = T.tokenCharacterID
			WHERE T.authCharacter = 1 AND botUserID = ?  AND entityID = ? AND standing = 10 LIMIT 1;`, memberID, entity).Scan(&ref); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return ref, nil
}

func (s *Conservator) getMilitiaStatus(memberID string, militia int32) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT characterName
			FROM evedata.botCharacters C
			INNER JOIN evedata.crestTokens T ON T.characterID = C.characterID
			WHERE T.authCharacter = 1 AND botUserID = ? AND factionID = ? LIMIT 1;`, memberID, militia).Scan(&ref); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return ref, nil
}

// FactionAllies
var FactionAllies = map[int32]int32{
	500001: 500003,
	500003: 500001,
	500002: 500004,
	500004: 500002,
}
