package tentez

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2Types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

func checkSwitchTargetGroupHealth(targetSwitch Switch, weight Weight, cfg Config) error {
	if weight.Old > 0 {
		isHealthy := false
		isNoTarget := false
		var err error
		for !isHealthy {
			fmt.Fprintf(cfg.io.out, "\rwaiting Healthy... last check: %s", time.Now().Format("15:04:05"))
			time.Sleep(1 * time.Second)

			isHealthy, isNoTarget, err = checkTargetGroupHealth(targetSwitch.Old, cfg)
			if err != nil {
				return err
			}
		}
		if isNoTarget {
			fmt.Fprintln(cfg.io.out, "\rno old targets\033[K")
		} else {
			fmt.Fprintln(cfg.io.out, "\rOld Target Healthy!\033[K")
		}
	}

	if weight.New > 0 {
		isHealthy := false
		isNoTarget := false
		var err error

		for !isHealthy {
			fmt.Fprintf(cfg.io.out, "\rwaiting Healthy... last check: %s\033[K", time.Now().Format("15:04:05"))
			time.Sleep(1 * time.Second)

			isHealthy, isNoTarget, err = checkTargetGroupHealth(targetSwitch.New, cfg)
			if err != nil {
				return err
			}
		}
		if isNoTarget {
			fmt.Fprintln(cfg.io.out, "\rno new targets\033[K")
		} else {
			fmt.Fprintln(cfg.io.out, "\rNew Target Healthy!\033[K")
		}
	}

	return nil
}

func checkTargetGroupHealth(arn string, cfg Config) (isHealthy bool, isNoTarget bool, err error) {
	out, err := cfg.client.elbv2.DescribeTargetHealth(context.TODO(), &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(arn),
	})
	if err != nil {
		return false, false, err
	}

	if len(out.TargetHealthDescriptions) == 0 {
		return true, true, nil
	}

	for _, t := range out.TargetHealthDescriptions {
		if t.TargetHealth.State != elbv2Types.TargetHealthStateEnumHealthy {
			return false, false, nil
		}
	}

	return true, false, nil
}
