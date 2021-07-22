package models

import (
	"errors"
	"time"

	log "github.com/gophish/gophish/logger"
)

// Blacklist contains the fields used for a Blacklist model
type Blacklist struct {
	Id           int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserId       int64     `json:"-" gorm:"column:user_id"`
	Name         string    `json:"name"`
	Ips          string    `json:"ips" gorm:"column:ips"`
	ModifiedDate time.Time `json:"modified_date"`
}

// ErrBlacklistNameNotSpecified is thrown if the name of the blacklist is blank.
var ErrBlacklistNameNotSpecified = errors.New("Blacklist Name not specified")

// [TODO] Check whether content is ok i.e. list of IPs
func (b *Blacklist) parseIPS() error {
	return nil
}

// Validate ensures that a blacklist contains the appropriate details
func (b *Blacklist) Validate() error {
	if b.Name == "" {
		return ErrBlacklistNameNotSpecified
	}
	return b.parseIPS()
}

// GetBlacklists returns the blacklists owned by the given user
func GetBlacklists(uid int64) ([]Blacklist, error) {
	bs := []Blacklist{}
	err := db.Where("user_id=?", uid).Find(&bs).Error
	if err != nil {
		log.Error(err)
		return bs, err
	}
	return bs, err
}

// GetBlacklist returns the blacklist, if it exists, specified by the given id and user_id.
func GetBlacklist(id int64, uid int64) (Blacklist, error) {
	b := Blacklist{}
	err := db.Where("user_id=? and id=?", uid, id).Find(&b).Error
	if err != nil {
		log.Error(err)
	}
	return b, err
}

// GetBlacklistByName returns the blacklist, if it exists, specified by the given name and user_id.
func GetBlacklistByName(n string, uid int64) (Blacklist, error) {
	b := Blacklist{}
	err := db.Where("user_id=? and name=?", uid, n).Find(&b).Error
	if err != nil {
		log.Error(err)
	}
	return b, err
}

// PostBlacklist creates a new blacklist in the database.
func PostBlacklist(b *Blacklist) error {
	err := b.Validate()
	if err != nil {
		log.Error(err)
		return err
	}
	// Insert into the DB
	err = db.Save(b).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

// PutBlacklist edits an existing Blacklist in the database.
// Per the PUT Method RFC, it presumes all data for a blacklist is provided.
func PutBlacklist(b *Blacklist) error {
	err := b.Validate()
	if err != nil {
		return err
	}
	err = db.Where("id=?", b.Id).Save(b).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

// DeleteBlacklist deletes an existing blacklist in the database.
// An error is returned if a blacklist with the given user id and blacklist id is not found.
func DeleteBlacklist(id int64, uid int64) error {
	err := db.Where("user_id=?", uid).Delete(Blacklist{Id: id}).Error
	if err != nil {
		log.Error(err)
	}
	return err
}
