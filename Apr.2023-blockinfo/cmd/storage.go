package main

import (
	"blockinfo/apps/blockinfoapp"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	log "github.com/sirupsen/logrus"
)

var (
	blockinfoDB *gorm.DB
)

type BlockInfo struct {
	gorm.Model
	Timestamp int64  `gorm:"index:timestamp, unique, sort:desc"`
	Height    *int64 `gorm:"index:height"`
	Addresses *int64
}

func createStorage() gen.ServerBehavior {
	return &Storage{}
}

type Storage struct {
	gen.Server
	log *logrus.Entry
}

// Init invoked on a start this process.
func (s *Storage) Init(process *gen.ServerProcess, args ...etf.Term) error {
	db, err := gorm.Open(sqlite.Open("blockinfo.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&BlockInfo{})
	blockinfoDB = db

	s.log = log.WithFields(log.Fields{
		process.Name(): process.Self(),
	})
	s.log.Info("Storage process is started")
	return nil
}

// HandleInfo invoked if this process received message sent with Process.Send(...).
func (s *Storage) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	var data BlockInfo

	switch m := message.(type) {
	case blockinfoapp.MessageBlockData:
		s.log.Infof("received block data: %v", m)
		data.Timestamp = m.T
		data.Height = &m.V
		if statement := blockinfoDB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "timestamp"}},
			DoUpdates: clause.AssignmentColumns([]string{"height"}),
		}).Create(&data); statement.Error != nil {
			s.log.Error(statement.Error)
		}
	case blockinfoapp.MessageAddressData:
		s.log.Infof("received address data: %v", m)
		data.Timestamp = m.T
		data.Addresses = &m.V
		if statement := blockinfoDB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "timestamp"}},
			DoUpdates: clause.AssignmentColumns([]string{"addresses"}),
		}).Create(&data); statement.Error != nil {
			s.log.Error(statement.Error)
		}
	default:
		s.log.Warningf("unknown message %#v", message)
		return gen.ServerStatusOK
	}

	return gen.ServerStatusOK
}
