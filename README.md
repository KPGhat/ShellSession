# Shell Session

## Usage

```shell
gsh -port 2333
```

### Basic session manage

```shell
# list all connection
session -l

# interact with a session
session -i [id]

# execute command for all session
session -a [cmd]
```

### Context manage

```shell
# Create Context
context -c

# List All Context
context -l

# Enter The Context By id
context -i [id]
```

```shell
# manage command
add [id, all]
del [id, all]
list
sh [cmd]
exit
```

### Log manage

```shell
# open/close log info
log on
log off
```

### Clear not alive session

```shell
# clear not alive session
clear
# execute the echo command to check alive
# if result is error, delete the session
clear -a
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