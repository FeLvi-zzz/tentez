package tentez

type Target interface {
	getName() string
	execSwitch(weight Weight, isForce bool, cfg Config) error
}
type Targets interface {
	targetsSlice() []Target
	fetchData(cfg Config) (interface{}, error)
}
