package config

// Hooks a list of commands to execute after the templating is done.
type Hooks struct {
	Post []string `yaml:"post"`
}
