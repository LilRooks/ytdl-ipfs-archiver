# ytdl-ipfs-archiver
An IPFS cache backed youtube-dl

Usage: 

```
ytdl-ipfs-archiver [-cfg <config.edn>] [-tab <table.sqlite3>] [-bin </path/to/youtube-dl>] -- [youtube-dl arguments]
```

Currently uses web3.storage, and requires a token either in a configuration or as a `TOKEN` environment variable.

Configuration format is provided in [examples](./examples).

The linked youtube-dl binary must support the `--get-*` flags, and expects only one download per call.

Once the database is sufficiently populated, you can share it with others to have them download your cached version.

