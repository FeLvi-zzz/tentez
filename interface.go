package tentez

type Target interface {
	getName() string
	execSwitch(weight Weight, client Client) error
}
type Targets interface {
	targetsSlice() []Target
	fetchData(client Client) (interface{}, error)
}
