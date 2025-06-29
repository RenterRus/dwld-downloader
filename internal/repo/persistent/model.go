package persistent

type LinkModel struct {
	link, filename, path, text, stage_config, message string
	retry                                             int
}

type SQLRepo interface {
	Select(q string) ([]LinkModel, error)
	Upsert(LinkModel) ([]LinkModel, error)
}
