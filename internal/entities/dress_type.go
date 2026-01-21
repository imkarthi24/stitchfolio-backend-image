package entities

type DressType struct {
	*Model `mapstructure:",squash"`

	Name         string `json:"name"`
	Description  string `json:"description"`
	Measurements string `json:"measurements"` //CSV of mesurement types Hip, Waist, Chest
}

func (DressType) TableName() string {
	return "stich.DressTypes"
}

func (DressType) TableNameForQuery() string {
	return "\"stich\".\"DressTypes\" E"
}
