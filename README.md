# Tentez

This tool helps you switching traffic.

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
You enter "yes", Tentez will continue steps.
If you'd like to interrupt steps, Ctrl+C or "quit".
> yes

2 / 9 steps
Switch old:new = 70:30
1. tentez-web switched!
2. tentez-api switched!
3. tentez-foo switched!
switch all targets!

3 / 9 steps
Sleep 600s
Resume at 2022-02-05 15:10:03
remaining: 600s
...
remaining: 1s

Resume

4 / 9 steps
Pause
You enter "yes", Tentez will continue steps.
If you'd like to interrupt steps, Ctrl+C or "quit".
> yes

(...snip)

apply complete!
```

```console
# get target resources information.
$ tentez -f ./examples/example.yaml get
tentez-web:
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef
    weight: 0
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210
    weight: 100
tentez-api:
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef
    weight: 0
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210
    weight: 100
tentez-foo:
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef
    weight: 0
  - arn: arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210
    weight: 100
```

## Why is named "Tentez"?
A `tentetsuki` is `railroad switch` in Japanese. It is a mechanical device used to guide trains from one track to another. This tool switches traffic, like a "tentesuki".

"Tentez" pronounces "ten-tets".
