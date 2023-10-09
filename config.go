package tentez

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	rgt "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

type elbv2Client interface {
	ModifyListener(ctx context.Context, params *elbv2.ModifyListenerInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyListenerOutput, error)
	ModifyRule(ctx context.Context, params *elbv2.ModifyRuleInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyRuleOutput, error)
	DescribeRules(ctx context.Context, params *elbv2.DescribeRulesInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeRulesOutput, error)
	elbv2.DescribeListenersAPIClient
	elbv2.DescribeTargetGroupsAPIClient
}

type rgtClient interface {
	rgt.GetResourcesAPIClient
}

type Client struct {
	elbv2 elbv2Client
	rgt   rgtClient
}

type Config struct {
	client Client
	clock  Clock
}

type Clock interface {
	Sleep(duration time.Duration)
}

func NewConfig() (Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRetryer(func() aws.Retryer {
		return retry.AddWithMaxAttempts(retry.NewStandard(), 10)
	}))
	if err != nil {
		return Config{}, err
	}

	elbv2svc := elbv2.NewFromConfig(cfg)
	rgtsvc := rgt.NewFromConfig(cfg)

	return Config{
		client: Client{
			elbv2: elbv2svc,
			rgt:   rgtsvc,
		},
		clock: &RealClock{},
	}, nil
}

type RealClock struct{}

func (c RealClock) Sleep(duration time.Duration) {
	time.Sleep(duration)
}
