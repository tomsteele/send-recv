# send-recv

## Install

Download latest release [here](https://github.com/tomsteele/send-recv/releases/latest).

## Usage
```
Usage of ./send-recv:
  -c int
      Amount of worker goroutines to spawn. (default 100)
  -hex
      Payload provided to -send is hex encoded.
  -i string
      Input file containing newline separated list of hosts. (default: stdin)
  -o string
      Output file to write the results to. (default: stdout)
  -p int
      Port to connect to.
  -recv int
      How much data to recv before closing the connection. (default 1024)
  -send string
      Payload to send after connection. (default "id\n")
  -timeout int
      Timeout in millseconds to wait for connection. (default 500)
```
