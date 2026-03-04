# Cron

Module that implements a cron scheduler. Emits signals to outputs on schedule, allowing connected modules to react to timed events.

## Usage

```go
import "github.com/dipdup-net/indexer-sdk/pkg/modules/cron"

cronModule, err := cron.NewModule(cfg.Cron)
if err != nil {
    log.Panic(err)
}

cronModule.Start(ctx)
defer cronModule.Close()
```

## Configuration

YAML config with a `jobs` map — keys are job names (used as output names), values are cron patterns:

```yaml
cron:
  jobs:
    half_an_hour: "0 30 * * * *"
    every_minute: "* 1 * * * *"
    every_five_second: "@every 5s"
    every_second: "* * * * * *"
```

Uses [robfig/cron](https://github.com/robfig/cron) syntax with seconds precision (6 fields).

## Outputs

Each job creates an output with the same name as the job key. The output sends an empty struct on each trigger.

```go
// Connect cron job output to your module's input
modules.Connect(cronModule, customModule, "every_second", "my_input")
```

## Handling Events

```go
func (m *MyModule) listen(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case <-m.everySecond.Listen():
            log.Info().Msg("triggered every second")
        case <-m.everyFiveSecond.Listen():
            log.Info().Msg("triggered every 5 seconds")
        }
    }
}
```

Full example: [`examples/cron/`](/examples/cron/)
