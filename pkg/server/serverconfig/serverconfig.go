package serverconfig

// This is placed in a sub-pkg to avoid an import cycle since the config pkg needs to reference on the server config but the server also load a new configuration file which refers back to the config pkg. Placing the server config in a sub pkg solves the problem.
type Config struct {
	Port int `default:"8088" json:"port" yaml:"port"`
}
