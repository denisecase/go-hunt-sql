package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Team definition
type Team struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Creator   User      `json:"creator"`
	CreatorID uint32    `gorm:"not null" json:"creator_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare new item
func (p *Team) Prepare() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Creator = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate info against schema
func (p *Team) Validate() error {
	if p.Name == "" {
		return errors.New("Required Name")
	}
	if p.CreatorID < 1 {
		return errors.New("Required Creator")
	}
	return nil
}

// SaveTeam to storage
func (p *Team) SaveTeam(db *gorm.DB) (*Team, error) {
	var err error
	err = db.Debug().Model(&Team{}).Create(&p).Error
	if err != nil {
		return &Team{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorID).Take(&p.Creator).Error
		if err != nil {
			return &Team{}, err
		}
	}
	return p, nil
}

// FindAllTeams returns all
func (p *Team) FindAllTeams(db *gorm.DB) (*[]Team, error) {
	var err error
	teams := []Team{}
	err = db.Debug().Model(&Team{}).Limit(100).Find(&teams).Error
	if err != nil {
		return &[]Team{}, err
	}
	if len(teams) > 0 {
		for i := range teams {
			err := db.Debug().Model(&User{}).Where("id = ?", teams[i].CreatorID).Take(&teams[i].Creator).Error
			if err != nil {
				return &[]Team{}, err
			}
		}
	}
	return &teams, nil
}

// FindTeamByID returns item matching ID
func (p *Team) FindTeamByID(db *gorm.DB, pid uint64) (*Team, error) {
	var err error
	err = db.Debug().Model(&Team{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Team{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorID).Take(&p.Creator).Error
		if err != nil {
			return &Team{}, err
		}
	}
	return p, nil
}

// UpdateATeam updates item with new information
func (p *Team) UpdateATeam(db *gorm.DB) (*Team, error) {
	var err error
	err = db.Debug().Model(&Team{}).Where("id = ?", p.ID).Updates(Team{Name: p.Name, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Team{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CreatorID).Take(&p.Creator).Error
		if err != nil {
			return &Team{}, err
		}
	}
	return p, nil
}

// DeleteATeam removes it from storage
func (p *Team) DeleteATeam(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Team{}).Where("id = ? and creator_id = ?", pid, uid).Take(&Team{}).Delete(&Team{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
