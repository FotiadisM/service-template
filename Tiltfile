settings = {
    "allowed_contexts": ["k3d-fotiadism-local"],
    "enabled_resources": [
        "redis",
        "postgres",
        "postgres-dev",
        "service-template",
        "service-template-compile"
    ]
}
config.set_enabled_resources(settings.get("enabled_resources"))

# Extensions
load("ext://restart_process", "docker_build_with_restart")

allow_k8s_contexts(settings.get("allowed_contexts"))

k8s_yaml("./.tilt/k8s/redis.yaml")
k8s_yaml("./.tilt/k8s/postgres.yaml")
k8s_resource("postgres", port_forwards="5432")
k8s_yaml("./.tilt/k8s/postgres-dev.yaml")
k8s_resource("postgres-dev", port_forwards="5433:5432")

local_resource(
  "service-template-compile",
  "task build",
  deps=["./api/docs/", "./api/gen/", "./cmd/", "./internal/"]
)

docker_build_with_restart(
    "fotiadism/service-template",
    ".",
    dockerfile="./.tilt/Dockerfile",
    entrypoint="/app/bin/app",
    only=["./bin"],
    live_update=[sync("./bin", "/app/bin")]
)

k8s_yaml("./.tilt/k8s/service-template.yaml")
k8s_resource("service-template", port_forwards="8080")

k8s_yaml("./.tilt/k8s/swagger/swagger.yaml")
k8s_resource("swagger", port_forwards="8085:8080")

k8s_yaml("./.tilt/k8s/otel/jaeger.yaml")
k8s_resource("jaeger", port_forwards="16686")

k8s_yaml("./.tilt/k8s/otel/otel-collector.yaml")
k8s_resource("otel-collector", objects=["otel-collector:configmap"])

k8s_yaml("./.tilt/k8s/otel/prometheus.yaml")
k8s_resource("prometheus", objects=["prometheus-config:configmap"], port_forwards="9090")
