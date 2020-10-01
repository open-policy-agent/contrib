package kubernetes.validating.flux

image := "fluxcd/flux"

deny[msg] {
    # Ensure only applies to flux images
    containerImage := input.spec.template.spec.containers[i].image
    contains(containerImage, image)
    # Check if git poll interval arg is present
    args := input.spec.template.spec.containers[_].args
    contains(args[_], "--git-poll-interval")
    # If single digit is present, deny
    regex.match("--git-poll-interval=[0-9]m", args[_])
    msg := "--git-poll-interval must be at least 10m"
}

deny[msg] {
    # Ensure only applies to flux images
    containerImage := input.spec.template.spec.containers[i].image
    contains(containerImage, image)
    # Check if sync interval arg is present
    args := input.spec.template.spec.containers[_].args
    contains(args[_], "--sync-interval")
    # If single digit is present, deny
    regex.match("--sync-interval=[0-9]m", args[_])
    msg := "--sync-interval must be at least 10m"
}
