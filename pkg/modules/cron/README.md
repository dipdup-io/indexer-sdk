# Cron

Cron is the module which implements a cron specific parser and notifications about scheduler events.

## Usage

Usage of cron module is described by the [example](/examples/cron/).

Cron module implements interface `Module`. So you can use it like any other module. For exmaple:

```go
// create cron module
cronModule, err := cron.NewModule(cfg.Cron)
if err != nil {
    log.Panic(err)
}
// start cron module
cronModule.Start(ctx)

// your code is here

// close cron module
if err := cronModule.Close(); err != nil {
    log.Panic(err)
}
```

## Config

Default yaml config of cron module contains only one field `jobs`. It's a map of job names to cron pattern. Job names is used like subscription id in inner-message communication.

```yaml
cron:
  jobs:
    half_an_hour: "0 30 * * * *"
    every_minute: "* 1 * * * *"
    every_five_second: "@every 5s"
    every_second: "* * * * * *"
```

## Subscription on internal events

To subscribe on events from cron module you should write your custom component ([example](/examples/cron/custom.go)) inherited from `Subscriber`. Then you can subscribe on job from your config file.

```go
cronModule.Subscribe(customModule.Subscriber, "every_second")      // set job name as subscription id
cronModule.Subscribe(customModule.Subscriber, "every_five_second") // set job name as subscription id
```