# go-diff

Example:

```go
hunks := diff.Diff(
  []string{"aaa", "bbb", "aaa", "ccc", "xxx"},
  []string{"aaa", "ccc", "ddd", "ccc", "zzz"},
)
```
