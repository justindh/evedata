package nail

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("character", spawnCharacterConsumer)
	AddHandler("alliance", spawnAllianceConsumer)
	AddHandler("corporation", spawnCorporationConsumer)
	AddHandler("corporationHistory", spawnCorporationHistoryConsumer)
	AddHandler("allianceHistory", spawnAllianceHistoryConsumer)
	AddHandler("loyaltyStore", spawnLoyaltyStoreConsumer)
	AddHandler("characterAuthOwner", spawnCharacterAuthOwnerConsumer)
}

func spawnCharacterConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterHandler)))
}

func spawnAllianceConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.allianceHandler)))
}

func spawnCorporationConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.corporationHandler)))
}

func spawnCorporationHistoryConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.corporationHistoryHandler)))
}

func spawnAllianceHistoryConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.allianceHistoryHandler)))
}
func spawnLoyaltyStoreConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.loyaltyStoreHandler)))
}

func spawnCharacterAuthOwnerConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterAuthOwnerHandler)))
}

func (s *Nail) corporationHandler(message *nsq.Message) error {
	c := datapackages.Corporation{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	cacheUntil := time.Now().UTC().Add(time.Hour * 24 * 31)

	err = s.doSQL(`INSERT INTO evedata.corporations
		(corporationID,name,ticker,ceoID,allianceID,factionID,memberCount,updated,cacheUntil)
		VALUES(?,?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
		ON DUPLICATE KEY UPDATE ceoID=VALUES(ceoID), name=VALUES(name), ticker=VALUES(ticker), allianceID=VALUES(allianceID), 
		factionID=VALUES(factionID), memberCount=VALUES(memberCount),  
		updated=UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)
	`, c.CorporationID, c.Corporation.Name, c.Corporation.Ticker, c.Corporation.CeoId, c.Corporation.AllianceId, c.Corporation.FactionId, c.Corporation.MemberCount, cacheUntil)
	if err != nil {
		log.Println(err)
		return err
	}

	return s.addEntity(c.CorporationID, "corporation")
}

func (s *Nail) allianceHandler(message *nsq.Message) error {
	c := datapackages.Alliance{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	cacheUntil := time.Now().UTC().Add(time.Hour * 24 * 31)

	err = s.doSQL(`
		INSERT INTO evedata.alliances 
			(
				allianceID,
				name,
				shortName,
				executorCorpID,
				startDate,
				corporationsCount,
				updated,
				cacheUntil
			)
			VALUES(?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE executorCorpID = VALUES(executorCorpID),
				corporationsCount = VALUES(corporationsCount), 
				updated = UTC_TIMESTAMP(), 
				cacheUntil=VALUES(cacheUntil)
	`, c.AllianceID, c.Alliance.Name, c.Alliance.Ticker, c.Alliance.ExecutorCorporationId, c.Alliance.DateFounded.UTC().Format("2006-01-02 15:04:05"), len(c.AllianceCorporations), cacheUntil)
	if err != nil {
		log.Println(err)
		return err
	}

	return s.addEntity(c.AllianceID, "alliance")
}

func (s *Nail) corporationHistoryHandler(message *nsq.Message) error {
	c := datapackages.CorporationHistory{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	var values []string
	for _, e := range c.CorporationHistory {
		values = append(values, fmt.Sprintf("(%d,%q,%d,%d)",
			c.CharacterID, e.StartDate.UTC().Format("2006-01-02 15:04:05"), e.RecordId, e.CorporationId))
	}

	if len(values) > 0 {
		stmt := fmt.Sprintf("INSERT INTO evedata.corporationHistory (characterID,startDate,recordID,corporationID) VALUES %s ON DUPLICATE KEY UPDATE characterID=VALUES(characterID);", strings.Join(values, ",\n"))
		err = s.doSQL(stmt)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (s *Nail) allianceHistoryHandler(message *nsq.Message) error {
	c := datapackages.AllianceHistory{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	var values []string
	for _, e := range c.AllianceHistory {
		deleted := 0
		if e.IsDeleted {
			deleted = 1
		}

		values = append(values, fmt.Sprintf("(%d,%q,%d,%d,%d)",
			c.CorporationID, e.StartDate.UTC().Format("2006-01-02 15:04:05"), e.RecordId, e.AllianceId, deleted))
	}

	if len(values) > 0 {
		stmt := fmt.Sprintf("INSERT INTO evedata.allianceHistory (corporationID,startDate,recordID,allianceID,deleted) VALUES %s ON DUPLICATE KEY UPDATE corporationID=VALUES(corporationID);", strings.Join(values, ",\n"))
		err = s.doSQL(stmt)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (s *Nail) characterHandler(message *nsq.Message) error {
	c := datapackages.Character{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	cacheUntil := time.Now().UTC().Add(time.Hour * 24 * 31)

	err = s.doSQL(`
		INSERT INTO evedata.characters (characterID,name,bloodlineID,ancestryID,corporationID,allianceID,race,gender,securityStatus,updated,cacheUntil,birthDate,etag)
			VALUES(?,?,?,?,?,?,evedata.raceByID(?),?,?,UTC_TIMESTAMP(),?,?,?) 
			ON DUPLICATE KEY UPDATE corporationID=VALUES(corporationID), gender=VALUES(gender), allianceID=VALUES(allianceID),birthDate=VALUES(birthDate),
			securityStatus=VALUES(securityStatus), updated = UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil), etag=VALUES(etag)
	`, c.CharacterID, c.Character.Name, c.Character.BloodlineId, c.Character.AncestryId, c.Character.CorporationId, c.Character.AllianceId, c.Character.RaceId, c.Character.Gender, c.Character.SecurityStatus, cacheUntil, c.Character.Birthday, c.ETag)
	if err != nil {
		log.Println(err)
		return err
	}

	return s.addEntity(c.CharacterID, "character")
}

func (s *Nail) loyaltyStoreHandler(message *nsq.Message) error {
	c := datapackages.Store{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	var offers, requirements []string
	for _, offer := range c.Store {
		offers = append(offers, fmt.Sprintf("(%d,%d,%d,%d,%d,%d,%d)",
			offer.OfferId, c.CorporationID, offer.TypeId, offer.Quantity, offer.LpCost, 0, int(offer.IskCost)))
		for _, requirement := range offer.RequiredItems {
			requirements = append(requirements, fmt.Sprintf("(%d,%d,%d)",
				offer.OfferId, requirement.TypeId, requirement.Quantity))
		}
	}

	stmt := fmt.Sprintf("INSERT INTO evedata.lpOffers (offerID,corporationID,typeID,quantity,lpCost,akCost,iskCost) VALUES %s ON DUPLICATE KEY UPDATE akCost=VALUES(akCost), iskCost=VALUES(iskCost), lpCost=VALUES(lpCost);", strings.Join(offers, ",\n"))
	err = s.doSQL(stmt)
	if err != nil {
		log.Println(err)
		return err
	}

	stmt = fmt.Sprintf("INSERT IGNORE INTO evedata.lpOfferRequirements (offerID,typeID,quantity) VALUES %s ON DUPLICATE KEY UPDATE quantity=VALUES(quantity);", strings.Join(requirements, ",\n"))
	err = s.doSQL(stmt)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

func (s *Nail) addEntity(id int32, entityType string) error {
	return s.doSQL("INSERT INTO evedata.entities (id, type) VALUES (?,?) ON DUPLICATE KEY UPDATE id = id;", id, entityType)
}

func (s *Nail) characterAuthOwnerHandler(message *nsq.Message) error {
	c := datapackages.CharacterRoles{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.doSQL(`
		UPDATE evedata.crestTokens SET roles = ?
		WHERE characterID = ? AND tokenCharacterID = ?
	`, strings.Join(c.Roles.Roles, ","), c.CharacterID, c.TokenCharacterID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
