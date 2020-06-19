# Terror

## Keeping program terror in check, reveal where go code exploded 💥

Go stacktrace sometimes are hard read and figure out, especially when it is very nested.

Sometimes you just want to know just the basic error message and their location. This helps you to do that.

### Install

`go get github.com/ninja-software/terror`

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

### Note

Always use `terror.New()` or it will not trace. E.g.

```go
return terror.New(terror.ErrBadContext, "")  // good
```

not

```go
return terror.ErrBadContext  // bad
```

Blank string like `terror.New(err, "")`, will default to use err.Error() string.

### Output

![console_output](terror.png)

```

```
