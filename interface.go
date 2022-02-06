package tentez

type GetTargets interface {
	FetchData() (interface{}, error)
}
