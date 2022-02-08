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

# available resources
- AWS
  - Listener
    - forward target group
  - Listener Rule
    - forward target group

## Why is named "Tentez"?
A `tentetsuki` is `railroad switch` in Japanese. It is a mechanical device used to guide trains from one track to another. This tool switches traffic, like a "tentesuki".

"Tentez" pronounces "ten-tets".
