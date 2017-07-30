walle-go
===
[walle](https://github.com/Meituan-Dianping/walle/) is a tool for android written in Java.
This project re-implement its some features via [go](https://golang.org/)

[walle-cli.go](walle-cli.go) vs [walle-cli](https://github.com/Meituan-Dianping/walle/tree/master/walle-cli)

### Build ###
```
go build walle-cli.go
```
or
```
go install walle-cli.go
```

### Usage of walle-cli ###

It is a command line tool for processing channel apk

    walle-cli command [args]

The commands:  
- [`show`](#show)  print the channel info for specified apks
- [`gen`](#gen)    generate apks with specified channel info 

#### show ####
```
walle-cli show [-r] <files...>
      -h  help
        print help message of command `show`
      -r  raw
        print raw text associated to id 0x71777777
```

e.g.

Show channel of files:  

```
walle-cli show /foo/bar/A.apk /foo/bar/bar/B.apk
```

Show raw channel info (may contains extras) of files:  

```
walle-cli show -r /foo/bar/A.apk /path/to/B.apk
```

#### gen  ####
```
walle-cli gen [-o out] [-f] [-d] -c <channel> [-e extras] <file>
      -c  channel(s)
        generate apk with specified channel(s), split multiple channels with ','
      -d  debug
        print debug log
      -e  extras
        generate apk with specified extras info (key value pairs, e.g thing=test,boom=1)
      -f  force
        force to overwrite existing channeled apk in output directory
      -h  help
        print help message of command `gen`
      -o  output
        output dir, generated channel apk(s) will store in here. default is input's dir
```
e.g.

Generate channel `babala` apk to outdir `/foo/bar/channel/` :  

```
walle-cli gen -o /foo/bar/channel/ -c babala /foo/bar/bar/B.apk
```

Generate channel `babala` apk to same dir of input :  

```
walle-cli gen -c babala /foo/bar/bar/B.apk
```

Generate channels `babala,balala,balaba` apks  to outdir `/foo/bar/channel/` :  

```
walle-cli gen -o /foo/bar/channel/ -c babala,balala,balaba /foo/bar/A.apk
```

Generate apk with channel `babala` and extras `a=1,b=true` to the same dir of input:  

```
walle-cli gen -c babala -e a=1,b=true /foo/bar/A.apk
```
