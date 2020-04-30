# Terror

## Keeping program terror in check, reveal where go code exploded

Go stacktrace sometimes are hard read and figure out, especially when it is very nested.

Sometimes you just want to know just the basic error message and their location. This helps you to do that.

### Usage

caller

```go
resp, err := http.Get("http://example.com/")
if err != nil {
    return nil, terror.New(err, "get website")
}
```

recover

```go
terror.Echo(err)
```

### Output

![console_output](terror.png)
