# Overview

With the combination of

![overview](/img/overview.png)

- switch
- a wrapper bash script
- timewarrior & taskwarrior
- leapp
- iterm2

You get an automated switch with the aws profile to:

1) Start a time tracker for this project
2) Switch to the project directory
3) Set a badge in terminal
4) populate AWS environment variables
5) Open a project url

![After](/img/after-switch.png)

**4** Populated AWS environment variables
```bash
env |grep AWS
AWS_DEFAULT_REGION=eu-central-1
AWS_ACCESS_KEY_ID=ASIA******B
AWS_SECRET_ACCESS_KEY=n******s
AWS_SESSION_TOKEN=IQoJ********AmCDxKA=
AWS_DEFAULT_PROFILE=letsbuild
AWS_REGION=eu-central-1
```

The switch is done on two steps:

1) Start session in leap
2) call `switch profilename`

## AWS SSO & Cross Account access

With [Leapp](https://docs.leapp.cloud/0.9.0/) you can use AWS SSO, IAM credentials or cross account roles.

Leapp saves the credentials not unsafe as text on your filesystem, but in an encrypted vault.


## Installation

### Copy binary in executable path

```bash
cp -pr dist/switchaws /usr/local/bin/switchaws
```
This is the call if you build switch yourself.

Or use the precompiled binaries from the release.

### Copy Bash script

```bash
cp switchawswrapper.sh /usr/local/bin
```

Or in an other directory which is in you $PATH.

### Set alias

```bash
alias switch='source /usr/local/bin/switchawswrapper.sh'
```

Because environment variables from the calling process has to be set, this extra step is neccesary.

## Start a session with leapp

### Configure a session in leapp

![Session](/img/leapp-session.png)

In this example I have a AWS-SSO integration called "ggadmin". Within this SSO login, there is a role "AWSPowerUserAccess". That session is calles "letsbuild" and uses the profile "default".

The other example has a direct ACCESS/SECRET authenticated AWS IAM user. This goes to the named profile "letsbuild".

Before starting a session with leapp, your credentials file is empty:

```bash
l ~/.aws/credentials
-rw-r--r--  1 pparker  staff  0  3 Mär 16:45 /Users/pparker/.aws/credentials
```

No credentials saved in clear text!

#### Start a session

![Start Session](/img/start-session.png)

```bash
l ~/.aws/credentials
-rw-r--r--  1 pparker  staff  835  3 Mär 16:48 /Users/pparker/.aws/credentials
```

Now *only* the started session is created in the `credentials` file.

### Switch to session

#### AWS credentials

If you use the "default" session, all cli and sdk commands would use the "letsbuild" credentials.
With the named profile "letsbuild", you have to populate the AWS environment variables.

You may name the profiles whatever you want, "letsbuild" is just for this example.

Before:

![before](/img/before-switch.png)

After:

```bash
env |grep AWS
AWS_DEFAULT_REGION=eu-central-1
AWS_ACCESS_KEY_ID=ASIA******B
AWS_SECRET_ACCESS_KEY=n******s
AWS_SESSION_TOKEN=IQoJ********AmCDxKA=
AWSUME_PROFILE=letsbuild
AWS_DEFAULT_PROFILE=letsbuild
AWS_REGION=eu-central-1
```
![After](/img/after-switch.png)


#### Timetracking

With [timewarrior](https://timewarrior.net/docs/) and [taskwarrior](https://taskwarrior.org/docs/) I have a complete cli based time tracking system:

If letsbuild is added as todo:

```bash
ID Active Age   Project   Tag               Due   Description                  Urg
27        58s                                     letsbuild                       0
```

And the `ID` is configured in `~/.aws/config`:

```bash
[profile letsbuild]
...
taskwarrior=27
```



Taskwarrior Project

```bash
timew su
```


```bash
Wk Date       Day Tags                                           Start      End    Time   Total
W9 2022-03-03 Thu letsbuild                                      16:55:46        - 0:01:45 7:10:48
```



## Extra config entrys

Example:

```bash
  1 [profile letsbuild]
  2 region=eu-central-1
  3 workdir=/Users/peterp/letsbuild/lambdaproject
  4 itermbadge=letsbuild
  5 taskwarrior=27
  6 url=/Users/peterp/letsbuild/lambdaproject/draw/fancydiagramm.drawio
```


The following entries in `~/.aws/config`are parsed:

## workdir

If current directory is outside `workdir`, then path is changed to workdir. If you are already inside workdir, path is not switched.

## itermbadge

If you use [iterm2](https://iterm2.com/documentation.html) then the screen badge is set to this string


## taskwarrior

If you use [taskwarrior](https://taskwarrior.org/docs/), the task is started. So you can track to working time per project.


## url

If `url` is set, then the url will be opened. If this is supported from you OS, not only websites but other apps are opened also.