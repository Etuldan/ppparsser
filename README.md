# ppparsser

## License

[GPLv3](https://www.gnu.org/licenses/quick-guide-gplv3.en.html), see LICENSE file

## Compiling

`go build -o ./ppparsser`

## Building on Docker

`docker build -t ppparsser .`

## Deploying on Docker (compose)

```yml
services:
    ppparsser:
        build: .
        container_name: ppparsser
        restart: always
        ports:
            - 8080:8080
```

## Example

- [info-communes.fr](https://info-communes.fr)