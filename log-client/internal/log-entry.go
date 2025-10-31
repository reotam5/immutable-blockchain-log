package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"time"
)

type LogEntry struct {
	ID        uint `gorm:"primaryKey"`
	Content   string
	Timestamp time.Time
	Source    string
}

type DetailedLogEntry struct {
	ID        uint
	Content   string
	Timestamp time.Time
	IsValid   bool
	Source    string
}

func (l LogEntry) Hash() (string, error) {
	v := reflect.ValueOf(l)
	t := v.Type()

	// concatenate all exported field
	data := ""
	for i := 0; i < v.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}

		field := v.Field(i).Interface()

		switch val := field.(type) {
		case time.Time:
			data += val.UTC().Format(time.RFC3339Nano)
		default:
			data += fmt.Sprintf("%v", val)
		}
	}

	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (l *LogEntry) GetDetailedLogEntry(hash string) (detailedLogEntry *DetailedLogEntry, err error) {
	isValid, err := l.ValidateHash(hash)
	if err != nil {
		return nil, err
	}
	dle := DetailedLogEntry{
		ID:        l.ID,
		Content:   l.Content,
		Timestamp: l.Timestamp,
		Source:    l.Source,
		IsValid:   isValid,
	}
	return &dle, nil
}

func (l *LogEntry) ValidateHash(hash string) (bool, error) {
	err := l.LoadFromDB(l.ID)
	if err != nil {
		return false, err
	}

	actualHash, err := l.Hash()

	if err != nil {
		return false, err
	}

	return actualHash == hash, nil
}

func (l *LogEntry) LoadFromDB(id uint) error {
	db, err := InitDB()
	if err != nil {
		return err
	}

	if err := db.First(&l, id).Error; err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) WriteToDB() error {
	db, err := InitDB()
	if err != nil {
		return err
	}

	if db.Create(&l).Error != nil {
		return fmt.Errorf("failed to write log entry to database: %w", err)
	}
	return nil
}

func (l LogEntry) String() string {
	return fmt.Sprintf("LogEntry[ID=%d, Content=%s, Timestamp=%s]", l.ID, l.Content, l.Timestamp.Format(time.RFC3339Nano))
}
