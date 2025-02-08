package inventory

type InventoryTemplater struct {
	publicBaseURL string
}

type InventoryTemplaterConfig struct {
	PublicBaseURL string
}

func NewInventoryTemplater(config InventoryTemplaterConfig) *InventoryTemplater {
	return &InventoryTemplater{
		publicBaseURL: config.PublicBaseURL,
	}
}
