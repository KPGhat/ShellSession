# Shell Session

## Usage

```shell
gsh -port 2333
```

```shell
# list all connection
session -l

# interact with a session
session -i [id]

# manager the session
session -m

# manage command
  add [id, all]
  del [id, all]
  sh [cmd]
  exit

# execute command for all session
session -a [cmd]
```

## Task List

- [x] interact with shell
- [x] Upload file
  - [ ] Upload file in chunks
- [x] Execute command for all shell
- [ ] Download file
- [ ] Compatible with all operating systems
  - [x] Linux
  - [ ] Windows
  - [ ] Darwin
- [ ] Setup proxy or port forward