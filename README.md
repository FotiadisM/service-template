# Mock microservice

## Getting Started

This project uses [devbox](https://github.com/jetify-com/devbox) to manage its development environment.

Install devbox:

```sh
curl -fsSL https://get.jetpack.io/devbox | bash
```

Start the devbox shell:

```sh
devbox shell
```

## TODO

- [ ] Make sure return bodies for panics, errors is all okay
- [ ] Improve configuration of packages
- [ ] Testing infrastructure
  - [x] Spawn test DB with Testcontainers
  - [ ] Create test utility func in `test/` for creating server, applying migraiton, fresh DB for every test etc.
- [ ] Create github actions for validating code gen. (Need to solve: How do we handle depandabot's PRs?)
- [ ] Finilaize project
  - [ ] Create new service
  - [ ] Update README
- [ ] Update documentation on packages.
