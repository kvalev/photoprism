variable "TAG" {
    default = "latest"
}
variable "DOCKER_REPO" {
    default = "photoprism/photoprism"
}

group "default" {
    targets = ["arm", "amd64"]
}

group "arm" {
    targets = ["armv7", "arm64"]
}

group "uffizzi" {
    targets = ["amd64", "demo-uffizzi"]
}

target "armv7" {
    dockerfile = "docker/photoprism/armv7/Dockerfile"
    tags = ["${DOCKER_REPO}:${TAG}-armv7"]
    platforms = ["linux/arm/v7"]
}

target "arm64" {
    dockerfile = "docker/photoprism/bookworm/Dockerfile"
    tags = ["${DOCKER_REPO}:${TAG}-arm64"]
    platforms = ["linux/arm64"]
}

target "amd64" {
    dockerfile = "docker/photoprism/bookworm/Dockerfile"
    tags = ["${DOCKER_REPO}:${TAG}-amd64"]
    platforms = ["linux/amd64"]
}

target "demo" {
    dockerfile = "docker/demo/Dockerfile"
    tags = ["${DOCKER_REPO}:${TAG}-demo"]
}

target "demo-uffizzi" {
    inherits = ["demo"]
    contexts = {
        "photoprism/photoprism:preview-ce" = "target:amd64"
    }
}
