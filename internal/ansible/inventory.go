package ansible

type Inventory struct {
	All struct {
		Children map[string]interface{} `yaml:"children"`
	} `yaml:"all"`
}
