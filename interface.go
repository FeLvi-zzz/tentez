package tentez

import "context"

type Target interface {
	getName() string
	execSwitch(ctx context.Context, weight Weight, isForce bool, cfg Config) error
}
type Targets interface {
	targetsSlice() []Target
	fetchData(ctx context.Context, cfg Config) (TargetsData, error)
}
type TargetsData interface{}
