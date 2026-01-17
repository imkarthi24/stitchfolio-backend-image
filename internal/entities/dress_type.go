package entities

type DressType struct {
	*Model `mapstructure:",squash"`

	Name         string `json:"name"`
	Measurements string `json:"measurements"`
}

func (DressType) TableName() string {
	return "DressTypes"
}

func (DressType) TableNameForQuery() string {
	return "\"DressTypes\" E"
}
