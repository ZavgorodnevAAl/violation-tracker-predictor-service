package config

type Service struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Logger struct {
	LogInFile bool   `yaml:"logInFile"`
	OutputDir string `yaml:"outputDir"`
}

type Cors struct {
	Target string `yaml:"target"`
	Path   string `yaml:"path"`
}
