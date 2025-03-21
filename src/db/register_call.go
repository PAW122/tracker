package db

import (
	"encoding/json"
	"fmt"
	"time"

	dbclient "github.com/PAW122/TsunamiDB/lib/dbclient"
	uuid "github.com/google/uuid"
)

type IpRecord struct {
	Ip           string    `json:"ip"`
	Id           string    `json:"id"`
	RegisterTime time.Time `json:"register_time"`
}

type RaportRecord map[string]Record

type Record struct {
	TotalTime int `json:"total_time"`
}

func GetRegisterId(ip string) (IpRecord, error) {

	isIpRegistered := false

	data, err := dbclient.Read(ip, "ip_table")
	if err != nil {
		isIpRegistered = false
	}

	var IpRecord IpRecord
	err = json.Unmarshal(data, &IpRecord)
	if err != nil {
		isIpRegistered = false
	}
	isIpRegistered = true

	if isIpRegistered && IpRecord.Ip == ip {
		return IpRecord, nil
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return IpRecord, fmt.Errorf("err generowania UUID: %v ", err)
		}

		IpRecord.Ip = ip
		IpRecord.Id = id.String()
		IpRecord.RegisterTime = time.Now()

		strJson, err := json.Marshal(IpRecord)
		if err != nil {
			return IpRecord, fmt.Errorf("err serializacji danych: %v ", err)
		}

		err = dbclient.Save(ip, "ip_table", strJson)
		if err != nil {
			return IpRecord, fmt.Errorf("err zapisu danych: %v ", err)
		}

		// save first record in raport table
		var RaportRecord = make(RaportRecord)
		RaportRecord[time.Now().Format("0-0-0")] = Record{TotalTime: 0}

		strJson, err = json.Marshal(RaportRecord)
		if err != nil {
			return IpRecord, fmt.Errorf("err serializacji danych: %v ", err)
		}

		err = dbclient.Save(IpRecord.Id, "raport_table", strJson)
		if err != nil {
			return IpRecord, fmt.Errorf("err zapisu danych: %v ", err)
		}

		return IpRecord, nil

	}
}
