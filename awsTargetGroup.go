package tentez

import (
	"context"
	"errors"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2Types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

// checkTargetGroupsExistense returns the target groups which are not exist.
func checkTargetGroupsExistense(c elbv2Client, tgArns []string) ([]string, error) {
	// dedup
	m := map[string]struct{}{}
	for _, tgArn := range tgArns {
		m[tgArn] = struct{}{}
	}

	errTgs := []string{}
	for tgArn := range m {
		if _, err := c.DescribeTargetGroups(context.TODO(), &elbv2.DescribeTargetGroupsInput{
			TargetGroupArns: []string{tgArn},
		}); err != nil {
			var tgnf *elbv2Types.TargetGroupNotFoundException
			if errors.As(err, &tgnf) {
				errTgs = append(errTgs, tgArn)
			} else {
				return nil, err
			}
		}
	}

	return errTgs, nil
}
