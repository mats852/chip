# chip 🪵🪒

instead of logging, distill your thoughts into bit flags that you can easily query

## disclaimer 💬

this project is a work in progress, and is not yet ready for production use!

## concept

chip is a novel observability tool that aims to replace traditional verbose logging
with compact 64-bit flag arrays. rather than generating kilobytes of text logs,
chip allows you to flip bits to represent application states, drastically reducing
storage requirements, parsing and querying compute and network traffic.

## why use chip?

- **space efficient**: 8 bytes per report instead of kilobytes to megabytes of traditional logs
- **cost effective**: significant savings for large deployments
- **security centric**: reduces sensitive data leakage by avoiding free-form text
- **queryable**: easily filter and analyze application states using bitwise operations

## how it works

1. define constant bit flags for your application states
2. use your configured flags to set the state of your application
3. periodically export chip batches to a sink (e.g., a database, file, etc.)
4. use your favorite tool to query the exported chips

## example usage

```go
import "github.com/mats852/chip"

type myAppContextValueKey string

const (
    Chip_Memory_Threshold_80  uint64 = 0b0001 // first position
    Chip_Cpu_Threshold_50     uint64 = 1 << 2 // second position
    Chip_Some_Positional_Flag uint8  = 63
)

var myMainContext = context.Background()

myFirstChip := chip.New("well-known-uuid")

chipExporter := chip.NewExporter(chip, chip.ExporterOpts{/* ... */})

chipExporter.Add(myFirstChip) // chip is added and periodically exported

go chipExporter.Serve(context.TODO()) // github.com/thejerf/suture supervisable interface

// optional: pass the chip around in your application
context.WithValue(myMainContext, myAppContextValueKey("chip"), myFirstChip)

// use in your application
myFirstChip.Set(Chip_Memory_Threshold_80)
myFirstChip.Set(Chip_Cpu_Threshold_50)
myFirstChip.SetPositions(Chip_Some_Positional_Flag)
```
