package conf

type DBConfig struct {
	Debug    bool
	Connects map[string]string
}
