package configuration

type Configuration struct {
	Application Application `yaml:"application"`
}

type Application struct {
	Port string `yaml:"port"`
}
