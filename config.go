package tentez

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type elbv2Client interface {
	ModifyListener(ctx context.Context, params *elbv2.ModifyListenerInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyListenerOutput, error)
	ModifyRule(ctx context.Context, params *elbv2.ModifyRuleInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyRuleOutput, error)
	DescribeRules(ctx context.Context, params *elbv2.DescribeRulesInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeRulesOutput, error)
	DescribeListeners(ctx context.Context, params *elbv2.DescribeListenersInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeListenersOutput, error)
}

type Client struct {
	elbv2 elbv2Client
}

type Config struct {
	client Client
}

func newConfig() (Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return Config{}, err
	}

	elbv2svc := elbv2.NewFromConfig(cfg)

	return Config{
		client: Client{
			elbv2: elbv2svc,
		},
	}, nil
}
