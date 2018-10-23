package models

func MaintContactSync() error {
	if _, err := RetryExec(`
        DELETE S.* FROM evedata.contactSyncs S
        LEFT OUTER JOIN evedata.crestTokens T ON S.destination = T.tokenCharacterID
        WHERE tokenCharacterID IS NULL;`); err != nil {
		return err
	}
	if _, err := RetryExec(`
        DELETE S.* FROM evedata.contactSyncs S
        LEFT OUTER JOIN evedata.crestTokens T ON S.source = T.tokenCharacterID
        WHERE tokenCharacterID IS NULL;`); err != nil {
		return err
	}

	return nil
}

func MaintOrphanCharacters() ([]int32, error) {
	ret := []int32{}
	err := database.Select(&ret, `
        SELECT A.corporationID from evedata.characters A
            LEFT OUTER JOIN evedata.corporations C ON C.corporationID = A.corporationID
            WHERE C.name IS NULL
        `)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func MaintMarket() error {
	if _, err := RetryExec(`
        UPDATE evedata.alliances A SET memberCount = 
            IFNULL(
                    (SELECT sum(memberCount) AS memberCount FROM evedata.corporations  C
                    WHERE C.allianceID = A.allianceID
                    GROUP BY allianceID LIMIT 1),
                    0
            );
            `); err != nil {
		return err
	}

	if _, err := RetryExec(`
        INSERT INTO evedata.discoveredAssets 
            SELECT 
                A.corporationID, 
                C.allianceID, 
                typeID, 
                K.solarSystemID, 
                K.x, 
                K.y, 
                K.z, 
                evedata.closestCelestial(K.solarSystemID, K.x, K.y, K.z) AS locationID, 
                MAX(killTime) as lastSeen 
            FROM evedata.killmailAttackers A
            INNER JOIN invTypes T ON shipType = typeID
            INNER JOIN evedata.corporations C ON C.corporationID = A.corporationID
            INNER JOIN evedata.killmails K ON K.id = A.id
            INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
            WHERE characterID = 0 AND groupID IN (365, 549, 1023, 1537, 1652, 1653, 1657, 2233)
            GROUP BY A.corporationID, solarSystemID, typeID
        ON DUPLICATE KEY UPDATE lastSeen = lastSeen;
            `); err != nil {
		return err
	}

	if _, err := RetryExec(`
        INSERT INTO evedata.discoveredAssets 
            SELECT 
                K.victimCorporationID AS corporationID, 
                C.allianceID, 
                typeID, 
                K.solarSystemID, 
                K.x, 
                K.y, 
                K.z, 
                evedata.closestCelestial(K.solarSystemID, K.x, K.y, K.z) AS locationID, 
                MAX(killTime) as lastSeen 
            FROM evedata.killmails K
            INNER JOIN invTypes T ON K.shipType = typeID
            INNER JOIN evedata.corporations C ON C.corporationID = K.victimCorporationID
            INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
            WHERE victimCharacterID = 0 AND groupID IN (365, 549, 1023, 1537, 1652, 1653, 1657, 2233)
            GROUP BY K.victimCorporationID, solarSystemID, typeID
        ON DUPLICATE KEY UPDATE lastSeen = lastSeen;
            `); err != nil {
		return err
	}

	regions, err := GetMarketRegions()
	if err != nil {
		return err
	}

	if err := RetryExecTillNoRows(`
        DELETE LOW_PRIORITY FROM evedata.market 
            WHERE date_add(issued, INTERVAL duration DAY) < UTC_TIMESTAMP() OR 
            reported < DATE_SUB(utc_timestamp(), INTERVAL 3 HOUR)
            ORDER BY regionID, typeID ASC LIMIT 50000;
            `); err != nil {
		return err
	}

	if _, err := RetryExec(`
        DELETE LOW_PRIORITY FROM evedata.marketStations ORDER BY stationName;
             `); err != nil {
		return err
	}

	if _, err := RetryExec(`
        INSERT IGNORE INTO evedata.marketStations SELECT  stationName, M.stationID, Count(*) as Count
        FROM    evedata.market M
                INNER JOIN staStations S ON M.stationID = S.stationID
        WHERE   reported >= DATE_SUB(UTC_TIMESTAMP(), INTERVAL 5 DAY)
        GROUP BY M.stationID 
        HAVING count(*) > 2000
        ORDER BY stationName;
            `); err != nil {
		return err
	}

	if _, err := RetryExec(`
       UPDATE evedata.market_vol SET quantity = 0;
             `); err != nil {
		return err
	}

	rows, err := database.Query(`
        SELECT  itemID AS itemID,
            regionID AS regionID,
            AVG(low) AS low,
            AVG(mean) AS mean,
            AVG(high) AS high,
            SUM(quantity) AS quantity,
            SUM(orders) AS orders
        FROM
            evedata.market_history
        WHERE
            date > UTC_TIMESTAMP() - INTERVAL 5 DAY
        GROUP BY regionID , itemID`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var (
			itemID, regionID, quantity, orders int64
			low, mean, high                    float64
		)
		rows.Scan(&itemID, &regionID, &low, &mean, &high, &quantity, &orders)
		if _, err := RetryExec(`REPLACE INTO evedata.marketHistoryStatistics VALUES(?,?,?,?,?,?,?)`, itemID, regionID, low, mean, high, quantity, orders); err != nil {
			return err
		}
	}
	rows.Close()

	for _, region := range regions {
		if _, err := RetryExec(`
        REPLACE INTO evedata.market_vol (
            SELECT count(*) as number,sum(quantity)/7 as quantity, regionID, itemID 
                FROM evedata.market_history 
                WHERE date > DATE_SUB(UTC_TIMESTAMP(),INTERVAL 7 DAY) 
                AND regionID = ?
                GROUP BY regionID, itemID);
            `, region.RegionID); err != nil {
			return err
		}
	}

	if _, err := RetryExec(`
       DELETE FROM evedata.jitaPrice ORDER BY itemID;
             `); err != nil {
		return err
	}
	if _, err := RetryExec(`
        INSERT IGNORE INTO evedata.jitaPrice (
        SELECT S.typeID as itemID, buy, sell, high, low, mean, quantity FROM
            (SELECT typeID, min(price) AS sell FROM evedata.market WHERE regionID = 10000002 AND bid = 0 GROUP BY typeID) S
            INNER JOIN (SELECT typeID, max(price) AS buy FROM evedata.market WHERE regionID = 10000002 AND bid = 1 GROUP BY typeID) B ON S.typeID = B.typeID
            LEFT OUTER JOIN (SELECT itemID, max(high) AS high, avg(mean) AS mean, min(low) AS low, sum(quantity) AS quantity FROM evedata.market_history WHERE regionID = 10000002 AND date > DATE_SUB(UTC_DATE(), INTERVAL 4 DAY) GROUP BY itemID) H on H.itemID = S.typeID
        HAVING mean IS NOT NULL
        ) ORDER BY itemID;
            `); err != nil {
		return err
	}

	if _, err := RetryExec(`
       DELETE FROM evedata.iskPerLp ORDER BY typeID;
             `); err != nil {
		return err
	}

	if _, err := RetryExec(`
        INSERT IGNORE INTO evedata.iskPerLp (
        SELECT
                N.itemName,
                S.typeID,
                T.typeName,
                MIN(lpCost) AS lpCost,
                MIN(iskCost) AS iskCost,
                ROUND(MIN(C.buy),0) AS JitaPrice,
                ROUND(MIN(C.quantity),0) AS JitaVolume,
                ROUND(COALESCE(MIN(P.price),0) + iskCost, 0)  AS itemCost,
                ROUND(
                        (
                                ( MIN(S.quantity) * AVG(C.buy) ) -
                                ( COALESCE( MIN(P.price), 0) + iskCost )
                        )
                        / MIN(lpCost)
                , 0) AS ISKperLP,
                P.offerID
        FROM evedata.lpOffers S

        INNER JOIN invNames N ON S.corporationID = N.itemID
        INNER JOIN invTypes T ON S.typeID = T.typeID
        INNER JOIN evedata.jitaPrice C ON C.itemID = S.typeID

        LEFT OUTER JOIN         (
                                SELECT offerID, sum(H.sell * L.quantity) AS price
                                FROM evedata.lpOfferRequirements L
                                INNER JOIN evedata.jitaPrice H ON H.itemID = L.typeID
                                GROUP BY offerID
                        ) AS P ON S.offerID = P.offerID

        GROUP BY S.offerID, S.corporationID
        HAVING ISKperLP > 0) ORDER BY typeID;
            `); err != nil {
		return err
	}

	return nil
}

func MaintKillMails() error { // Broken into smaller chunks so we have a chance of it getting completed.
	// Delete stuff older than 90 days, we do not care...
	if err := RetryExecTillNoRows(`
				DELETE A.* FROM evedata.killmailAttackers A
		            JOIN (SELECT id FROM evedata.killmails WHERE killTime < DATE_SUB(UTC_TIMESTAMP, INTERVAL 365 DAY) LIMIT 50000) K ON A.id = K.id;
		            `); err != nil {
		return err
	}
	if err := RetryExecTillNoRows(`
				DELETE A.* FROM evedata.killmailItems A
		        JOIN (SELECT id FROM evedata.killmails WHERE killTime < DATE_SUB(UTC_TIMESTAMP, INTERVAL 365 DAY) LIMIT 50000) K ON A.id = K.id;
		            `); err != nil {
		return err
	}
	if err := RetryExecTillNoRows(`
				DELETE FROM evedata.killmails
		        WHERE killTime < DATE_SUB(UTC_TIMESTAMP, INTERVAL 365 DAY) LIMIT 50000;
		            `); err != nil {
		return err
	}

	// Remove any invalid items
	/*if err := RetryExecTillNoRows(`
		        DELETE D.* FROM evedata.killmailAttackers D
	            JOIN (SELECT A.id FROM evedata.killmailAttackers A
					 LEFT OUTER JOIN evedata.killmails K ON A.id = K.id
		             WHERE K.id IS NULL LIMIT 10) S ON D.id = S.id;
		               `); err != nil {
			return err
		}*/
	if err := RetryExecTillNoRows(`
			DELETE D.* FROM evedata.killmailItems D 
            JOIN (SELECT A.id FROM evedata.killmailItems A
				 LEFT OUTER JOIN evedata.killmails K ON A.id = K.id
	             WHERE K.id IS NULL LIMIT 10) S ON D.id = S.id;
	               `); err != nil {
		return err
	}

	// Prefill stats for known entities that may have no kills
	if _, err := RetryExec(`
        INSERT IGNORE INTO evedata.entityKillStats (id)
	    (SELECT corporationID AS id FROM evedata.corporations WHERE memberCount > 0); 
            `); err != nil {
		return err
	}

	if _, err := RetryExec(`
        INSERT IGNORE INTO evedata.entityKillStats (id)
	    (SELECT allianceID AS id FROM evedata.alliances); 
            `); err != nil {
		return err
	}

	// Build entity stats
	if _, err := RetryExec(`
        INSERT INTO evedata.entityKillStats (id, losses)
            (SELECT 
                victimCorporationID AS id,
                COUNT(DISTINCT K.id) AS losses
            FROM evedata.killmails K
            WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 180 DAY)
            GROUP BY victimCorporationID
            ) ON DUPLICATE KEY UPDATE losses = values(losses);
            `); err != nil {
		return err
	}
	if _, err := RetryExec(`
        INSERT INTO evedata.entityKillStats (id, losses)
            (SELECT 
                victimAllianceID AS id,
                COUNT(DISTINCT K.id) AS losses
            FROM evedata.killmails K
            WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 180 DAY)
            GROUP BY victimAllianceID
            ) ON DUPLICATE KEY UPDATE losses = values(losses);
            `); err != nil {
		return err
	}

	if _, err := RetryExec(`
        INSERT INTO evedata.entityKillStats (id, kills)
            (SELECT 
                corporationID AS id,
                COUNT(DISTINCT K.id) AS kills
            FROM evedata.killmails K
            INNER JOIN evedata.killmailAttackers A ON A.id = K.id
            WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 180 DAY)
            GROUP BY A.corporationID
            ) ON DUPLICATE KEY UPDATE kills = values(kills);
            `); err != nil {
		return err
	}

	if _, err := RetryExec(`
        INSERT INTO evedata.entityKillStats (id, kills)
            (SELECT 
                allianceID AS id,
                COUNT(DISTINCT K.id) AS kills
            FROM evedata.killmails K
            INNER JOIN evedata.killmailAttackers A ON A.id = K.id
            WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 180 DAY)
            GROUP BY A.allianceID
            ) ON DUPLICATE KEY UPDATE kills = values(kills);
            `); err != nil {
		return err
	}

	// Update everyone efficiency
	if _, err := RetryExec(`
        UPDATE evedata.entityKillStats SET efficiency = IF(losses+kills, (kills/(kills+losses)) , 1.0000);
            `); err != nil {
		return err
	}

	return nil
}
