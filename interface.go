package tentez

type Target interface {
	getName() string
	execSwitch(weight Weight) error
}
type Targets interface {
	fetchData() (interface{}, error)
}
