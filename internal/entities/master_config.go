package entities

type MasterConfig struct {
	*Model `mapstructure:",squash"`
	Name   string `json:"name,omitempty"` //eg: course , qualifiications
	Type   string `json:"type,omitempty"`

	CurrentValue  string `json:"current_value,omitempty"`
	PreviousValue string `json:"previous_value,omitempty"`
	DefaultValue  string `json:"default_value,omitempty"`
	UseDefault    bool   `json:"use_default,omitempty"`

	Description string `json:"description,omitempty"`

	Format string `json:"format,omitempty"`
}

func (MasterConfig) TableName() string {
	return "MasterConfigs"
}

func (MasterConfig) TableNameForQuery() string {
	return "\"MasterConfigs\" E"
}

// Type.Name  -> CandidateForm.Courses
