# Run

You can run the app with port and/or test

if you run with test it will not set the folder and output to the console

```bash
go run main.go -t
```

you can also change the port

```bash
go run main.go -p 8501
```

# Test

```bash
go test ./...
```

# Brief Description
This is a simple app that counts the number of connections to a port and sends it to a folder for prometheus to read

I am using the [go-bugfixes](https://github.com/bugfixes/go-bugfixes) logging library for this app, which is my own so that it had things like timestamps and colors

I am using the [gopsutil](https://github.com/shirou/gopsutil) library to get the connections, that way I had a testable library that could listen on IPv4 and IPv6

I am using the [cobra](https://github.com/spf13/cobra) library to make the CLI to make it easier to have CLI functionality