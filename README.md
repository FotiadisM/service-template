# Mock microservice

## Getting Started

This project uses [devbox](https://github.com/jetify-com/devbox) to manage its development environment.

Install devbox:

```sh
curl -fsSL https://get.jetpack.io/devbox | bash
```

Start the devbox shell (this can be automated using [direnv](https://github.com/direnv/direnv) and `devbox generate direnv`):

```sh
devbox shell
```

## TODO

- [x] Make sure return bodies for panics, errors is all okay
- [ ] Improve configuration of packages
- [ ] Testing infrastructure
  - [x] Spawn test DB with Testcontainers
  - [x] Every test gets a fresh DB
  - [ ] Create more test utility func in `test/`.
  - [ ] Integrate snapshot testing
- [ ] Test for all packages
- [ ] Create github actions for validating code gen. (Need to solve: How do we handle depandabot's PRs?)
- [ ] Create efficient Dockerfile
- [ ] Finilaize project
  - [ ] Create new service
  - [ ] Update README
  - [ ] Update documentation on packages.

## Optional TODO

- [ ] Integrate some kubernetes development tool
