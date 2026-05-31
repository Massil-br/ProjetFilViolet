package models

type Table struct {
	Model
	SaloonID       uint   `json:"saloon_id" gorm:"foreignKey:SaloonID;references:ID"`
	SlotsAvailable uint64 `json:"slots_available"`
}
