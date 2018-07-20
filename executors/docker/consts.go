package docker

import "time"

const DockerAPIVersion = "1.18"
const dockerLabelPrefix = "com.gitlab.gitlab-runner"

const prebuiltImageName = "gitlab/gitlab-runner-helper"
const prebuiltImageExtension = ".tar.xz"

const dockerCleanupTimeout = 5 * time.Minute

const dockerAttachRetries = 3

var dockerAttachRetryTime = time.Second
