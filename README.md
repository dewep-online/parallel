# parallel
Parallel command execution


# install

```go
 go install github.com/dewep-online/parallel@latest
```

# how work

```shell
Current Command: 
  parallel  [arg]  --shell=bash --exit --timeout=3

Flags:
  --shell      default shell (bash, sh, ... etc) (default: bash)
  --exit       stop parallel if error or exit 1 (default: false)
  --timeout    restart timeout in sec if error or exit 1 (default: 3)

```

# example

```shell
parallel "ping google.com" "ping yandex.ru" "date" "exit 1" --timeout=1 --exit
```