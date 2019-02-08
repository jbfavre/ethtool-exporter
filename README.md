# ethtool-exporter

Export ethtool metrics to prometheus compatible format

## Usage

```./ethtool-exporter --help
Usage of ./ethtool-exporter:
  -ifaceregexp string
    	an interface name or regexp (default ".*")
  -output string
    	an existing directory to store file containing metrics (default "/prom_output/ethtool.prom")
  -sleep int
    	time in second to wait between two statistics gathering (default 20)
```

## Under the hood

Uses [github.com/safchain/ethtool](https://github.com/safchain/ethtool) to mimic ethtool binary
behaviour and gather network cards metrics.

By default, all network card metrics are polled. You can restrict polled network card using `-ifaceregexp`
option. This option supports either an interface name or a regular expression.

Metrics are then formatted to fit prometheus conventions and stored in an output file, which can be customized
using `-output` option. Please note that any directory iin the path **must** exist.
`ethtool-exporter` won't try to create them.

This file can be polled by an existing prometheus-exported which has access to the directory.

## Contributing

Fork this repo, work, open a pull request, wait for TRavis CI result
