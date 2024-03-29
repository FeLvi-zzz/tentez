# Tentez

Tentez helps you switching traffic.

## Installation

If you don't want to build from source go grab a [binary release](https://github.com/FeLvi-zzz/tentez/releases).

or use `go install`

```
$ go install github.com/FeLvi-zzz/tentez/cmd/tentez@latest
```

## Usage

```console
# show plan
$ tentez -f ./examples/example.yaml plan
Plan
1. pause
2. switch old:new = 70:30
  - tentez-web
  - tentez-api
  - tentez-foo
3. sleep 600s
4. pause
5. switch old:new = 30:70
  - tentez-web
  - tentez-api
  - tentez-foo
6. sleep 600s
7. pause
8. switch old:new = 0:100
  - tentez-web
  - tentez-api
  - tentez-foo
9. sleep 600s
```

```console
# show plan and apply
$ tentez -f ./examples/example.yaml apply
Plan
1. pause
2. switch old:new = 70:30
  1. tentez-web
  2. tentez-api
  3. tentez-foo
3. sleep 600s
4. pause
5. switch old:new = 30:70
  1. tentez-web
  2. tentez-api
  3. tentez-foo
6. sleep 600s
7. pause
8. switch old:new = 0:100
  1. tentez-web
  2. tentez-api
  3. tentez-foo
9. sleep 600s

1 / 9 steps
Pause
enter "yes", continue steps.
If you'd like to interrupt steps, enter "quit".
> yes

2 / 9 steps
Switch old:new = 70:30
1. tentez-web switched!
2. tentez-api switched!
3. tentez-foo switched!

3 / 9 steps
Sleep 600s
Resume at 2022-02-05 15:10:03
Remain: 600s
...
Remain: 1s
Resume

(...snip)

Apply complete!
```

```console
# get target resources' current states.
$ tentez -f ./examples/example.yaml get
aws_listeners:
- target: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:listener/app/my-lb/0123456789abcdef/0123456789abcdef
  weights:
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef
    weight: 0
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210
    weight: 100
aws_listener_rules:
- target: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:listener/app/my-lb/0123456789abcdef/0123456789abcdef/0123456789abcdef
  weights:
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef
    weight: 0
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210
    weight: 100
- target: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:listener/app/my-lb/0123456789abcdef/0123456789abcdef/0123456789abcdef
  weights:
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef
    weight: 0
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210
    weight: 100
```

```console
# rollback
# if you want to skip the pause step, add `--no-pause` flag.
$ tentez -f ./examples/example.yaml rollback
1. pause
2. switch old:new = 100:0
  1. tentez-web
  2. tentez-api
  3. tentez-foo

1 / 2 steps
Pause
enter "yes", continue steps.
If you'd like to interrupt steps, enter "quit".
> yes
continue step

2 / 2 steps
Switch old:new = 100:0
1. tentez-web switched!
2. tentez-api switched!
3. tentez-foo switched!
Switched at 2022-03-12 12:05:30

Apply complete!
```

```console
# switch weights
# this command overrides steps of the config file.
# if you want to skip the pause step, add `--no-pause` flag.
$ tentez -f ./examples/example.yaml switch --weights 70,30
1. pause
2. switch old:new = 70:30
  1. tentez-web
  2. tentez-api
  3. tentez-foo

1 / 2 steps
Pause
enter "yes", continue steps.
If you'd like to interrupt steps, enter "quit".
> yes
continue step

2 / 2 steps
Switch old:new = 70:30
1. tentez-web switched!
2. tentez-api switched!
3. tentez-foo switched!
Switched at 2023-05-25 13:18:25

Apply complete!
```

```console
# show version
$ tentez version
tentez version: x.x.x (rev: xxxxxxx)
```

```console
# generate config from terraform plan json
$ terraform plan -out tfplan && terraform show -json tfplan > tfplan.json
$ tentez generate-config tfplanjson -f ./tfplan.json -o tentez.yaml
```

For instance, you can generate a config from the below terraform diff.

```diff
 resource "aws_lb_listener" "example" {
   ...

   default_action {
     type             = "forward"
-    target_group_arn = aws_lb_target_group.old.arn
+    target_group_arn = aws_lb_target_group.new.arn
   }
 }
```

```console
# generate config from tagged AWS resouces
$ tentez generate-config resource-tag -f examples/tentez.ResourceTag.v1beta1.yaml
```

Refer examples/tentez.ResourceTag.v1beta1.yaml.

### Assume other IAM Role

```console
# set `AWS_ASSUME_ROLE_ARN` environment variable
$ AWS_ASSUME_ROLE_ARN=[IAM_ROLE_ARN] tentez -f ./examples/example.yaml get
```

## Available resources

- AWS
  - Listener
    - forward target group. for default LB listener rule.
  - Listener Rule
    - forward target group. for except default LB listner rule.

## Why is named "Tentez"?

A `tentetsuki` is `railroad switch` in Japanese. It is a mechanical device used to guide trains from one track to another. This tool switches traffic, like a "tentetsuki".

"Tentez" pronounces "ten-tets".
