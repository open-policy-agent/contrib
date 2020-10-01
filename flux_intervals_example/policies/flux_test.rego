package kubernetes.validating.flux

test_ten_minutes_allowed {
    count(deny) == 0 with input as {"spec":{"template":{"spec":{"containers":[{"image":"fluxcd/flux:1.20.2","args":["--git-poll-interval=10m","--sync-interval=10m"]}]}}}}
}

test_one_hour_allowed {
    count(deny) == 0 with input as {"spec":{"template":{"spec":{"containers":[{"image":"fluxcd/flux:1.20.2","args":["--git-poll-interval=1h","--sync-interval=1h"]}]}}}}
}

test_five_minutes_denied{
    deny with input as {"spec":{"template":{"spec":{"containers":[{"image":"fluxcd/flux:1.20.2","args":["--git-poll-interval=5m","--sync-interval=5m"]}]}}}}
}

test_five_minutes_denied_git_poll {
    count(deny) == 1 with input as {"spec":{"template":{"spec":{"containers":[{"image":"fluxcd/flux:1.20.2","args":["--git-poll-interval=5m","--sync-interval=10m"]}]}}}}
}

test_five_minutes_denied_sync_interval {
    count(deny) == 1 with input as {"spec":{"template":{"spec":{"containers":[{"image":"fluxcd/flux:1.20.2","args":["--git-poll-interval=10m","--sync-interval=5m"]}]}}}}
}
