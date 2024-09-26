# TestLogger

## Logger Usage

```go
func Example_test(t *testing.T) {
	logger := testlog.NewTestLogger(logruslog.NewLogrusLogger())
    ctx := logctx.CtxWithLogger(context.Background(), logger)
	
	err := RunTestFunction(ctx)

	if err != nil {
		logger.ShowStoredLogs()
    }
}
```

