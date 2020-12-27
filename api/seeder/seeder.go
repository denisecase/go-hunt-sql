package seeder

import (
	"github.com/denisecase/go-hunt-sql/api/models"
	"gorm.io/gorm"
)

var users = []models.User{
	{ID: 1, Email: "dcase@nwmissouri.edu", Password: "password"},
	{ID: 2, Email: "hoot@nwmissouri.edu", Password: "password"},
	{ID: 3, Email: "alex154590@gmail.com", Password: "password"},
	{ID: 4, Email: "samarpita.chandolu@outlook.com", Password: "password"},
	{ID: 5, Email: "bhanu1994@gmail.com", Password: "password"},
	{ID: 6, Email: "chandu131198@gmail.com", Password: "password"},
	{ID: 7, Email: "chanduhvg@gmail.com", Password: "password"},
	{ID: 8, Email: "p.harichandraprasad@gmail.com", Password: "password"},
	{ID: 9, Email: "krishna.ksk1996@gmail.com", Password: "password"},
	{ID: 10, Email: "mohansai03@outlook.com", Password: "password"},
	{ID: 11, Email: "prasad.gd@gmail.com", Password: "password"},
	{ID: 12, Email: "pruthvunaskanti@hotmail.com", Password: "password"},
	{ID: 14, Email: "raviteja.pagidoju@gmail.com", Password: "password"},
	{ID: 15, Email: "saikrish1545@gmail.com", Password: "password"},
	{ID: 16, Email: "teja2004@outlook.com", Password: "password"},
	{ID: 17, Email: "srkvodnala@gmail.com", Password: "password"},
	{ID: 18, Email: "csrisudheera@gmail.com", Password: "password"},
	{ID: 19, Email: "swaroopreddy.g@gmail.com", Password: "password"},
	{ID: 20, Email: "swaroopat@hotmail.com", Password: "password"},
	{ID: 21, Email: "kiran021997@gmail.com", Password: "password"},
	{ID: 22, Email: "yashwanthrocks@gmail.com", Password: "password"},
	{ID: 23, Email: "vishal041197@outlook.com", Password: "password"},
	{ID: 24, Email: "default@email.com", Password: "password"},
}

var teams = []models.Team{
	{ID: 1, Name: "Thunder Thinkers", CreatorID: 18},
	{ID: 2, Name: "Mavericks", CreatorID: 19},
	{ID: 3, Name: "Sunrisers Horizons", CreatorID: 15},
	{ID: 4, Name: "Barbies", CreatorID: 21},
	{ID: 5, Name: "Hunters", CreatorID: 5},
}

// Load the database with sample data
func Load(db *gorm.DB) {

	db.AutoMigrate(&models.User{}, &models.Team{})

	for i := range users {
		db.Model(&models.User{}).Create(&users[i])
	}

	for i := range teams {
		db.Model(&models.Team{}).Create(&teams[i])
	}

}
