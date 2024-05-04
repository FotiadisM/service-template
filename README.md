# Mock microservice

## TODO

- [x] Use Atlas migrations instead of a single file.
  - [ ] Debug formating, doesn't apply to later `atlas migrate diff`.
  - [ ] Debug workflow, do I have to restart postgres every time I apply and want to generate new migrations?
- [x] Rename sqlc `queries/` into repository (maybe?). Also allow for different repositories?
- [ ] Create github actions for validating code gen. (Need to solve: How do we handle depandabot's PRs?)
- [ ] Create auth gRPC middleware.
- [ ] Update `otelgrpc/` with the latest semantic conventions.
- [ ] Update or remove `otel/`.
- [ ] Improve configuration (including server and logging), also allow for `.env.local` maybe(?).
- [ ] Improve testing infrastracture.
- [ ] Change service name.
- [ ] Rename repository. (Create also template?)

## Nice to Have

- [ ] Update documentation on packages.
