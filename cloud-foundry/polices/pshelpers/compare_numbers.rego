package pshelpers

isLimitExceeded(limit, set) = exceeded { 
  exceedsLimit := [
    aboveLimit |
    set[_] > limit
    aboveLimit := true
  ]
  exceeded = count(exceedsLimit) > 0 
}
