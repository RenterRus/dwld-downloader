package app

func NewApp() error {
	ReadConfig(".", "config_example")

	return nil
}
