package hammer

import (
	"context"
	"log"
	"time"

	"encoding/gob"

	"github.com/antihax/goesi/esi"
)

func init() {
	registerConsumer("war", warConsumer)
	gob.Register(esi.GetWarsWarIdOk{})
}

func warConsumer(s *Hammer, parameter interface{}) {
	id := parameter.(int32)

	war, _, err := s.esi.ESI.WarsApi.GetWarsWarId(context.TODO(), id, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// if the war ended, market it finished
	if war.Finished.IsZero() == false && war.Finished.Before(time.Now().UTC()) {
		err = s.inQueue.SetWorkCompleted("evedata_known_kills", int64(id))
		if err != nil {
			log.Println(err)
		}
	}

	// Send out the result
	err = s.QueueResult(war, "war")
	if err != nil {
		log.Println(err)
		return
	}

	// Add the alliance corporation for intel purposes
	if war.Aggressor.AllianceId == 0 {
		err = s.AddCorporation(war.Aggressor.CorporationId)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// Add the alliance corporation for intel purposes
	if war.Defender.AllianceId == 0 {
		err = s.AddCorporation(war.Defender.CorporationId)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
