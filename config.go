package tentez

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type elbv2Client interface {
	ModifyListener(ctx context.Context, params *elbv2.ModifyListenerInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyListenerOutput, error)
	ModifyRule(ctx context.Context, params *elbv2.ModifyRuleInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyRuleOutput, error)
	DescribeRules(ctx context.Context, params *elbv2.DescribeRulesInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeRulesOutput, error)
	elbv2.DescribeListenersAPIClient
	elbv2.DescribeTargetGroupsAPIClient
}

type Client struct {
	elbv2 elbv2Client
}

type IOStreams struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

type Config struct {
	client Client
	io     IOStreams
	clock  Clock
}

type Clock interface {
	Sleep(duration time.Duration)
}

func NewConfig() (Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return Config{}, err
	}

	elbv2svc := elbv2.NewFromConfig(cfg)

	return Config{
		client: Client{
			elbv2: elbv2svc,
		},
		io: IOStreams{
			in:  os.Stdin,
			out: os.Stdout,
			err: os.Stderr,
		},
		clock: &RealClock{},
	}, nil
}

type RealClock struct{}

func (c RealClock) Sleep(duration time.Duration) {
	time.Sleep(duration)
}
