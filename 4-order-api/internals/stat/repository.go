package stat

import (
	"go/order-api/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	Database *db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Database: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	currentDate := datatypes.Date(time.Now())
	repo.Database.Find(&stat, "link_id = ? and date = ?", linkId, currentDate)
	if stat.ID == 0 {
		repo.Database.Create(&Stat{
			LinkID: linkId,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks += 1
		repo.Database.Save(&stat)
	}
}

func (repo *StatRepository) GetStat(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string
	switch by {
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"

	}
	repo.Database.
		Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", to, from).
		Group("period").
		Order("period").
		Scan(&stats)

	return stats
}
