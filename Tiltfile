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

def deploy_dependencies():
    k8s_yaml("./.tilt/k8s/redis.yaml")
    k8s_resource("redis", labels="dependencies")

    k8s_yaml("./.tilt/k8s/postgres.yaml")
    k8s_resource("postgres", objects=["postgres-init:configmap"], port_forwards="5432", labels="dependencies")

    k8s_yaml("./.tilt/k8s/postgres-dev.yaml")
    k8s_resource("postgres-dev", objects=["postgres-dev-init:configmap"], port_forwards="5433:5432", labels="dependencies")

def deploy_service_and_helpers():
    local_resource(
      "service-template-compile",
      "task build",
      deps=["./api/docs/", "./api/gen/", "./cmd/", "./internal/"],
      labels="service-template"
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
    k8s_resource("service-template", port_forwards="8080", labels="service-template")

    # Helpers, aka SwaggerUI and grpcUI
    k8s_yaml("./.tilt/k8s/grpcui/grpcui.yaml")
    k8s_resource("grpcui", port_forwards="8082:8080", labels="service-template")
    k8s_yaml("./.tilt/k8s/swagger/swagger.yaml")
    k8s_resource("swagger", port_forwards="8085:8080", labels="service-template")


def deploy_observability():
    k8s_yaml("./.tilt/k8s/observability/jaeger.yaml")
    k8s_resource("jaeger", port_forwards=["16686"], labels="observability")
    k8s_yaml("./.tilt/k8s/observability/otel-collector.yaml")
    k8s_resource("otel-collector", objects=["otel-collector:configmap"], labels="observability")
    k8s_yaml("./.tilt/k8s/observability/prometheus.yaml")
    k8s_resource("prometheus", objects=["prometheus-config:configmap"], port_forwards="9090", labels="observability")
    k8s_yaml("./.tilt/k8s/observability/grafana.yaml")
    k8s_resource("grafana",  port_forwards="3000", labels="observability")

deploy_dependencies()

deploy_service_and_helpers()

deploy_observability()
