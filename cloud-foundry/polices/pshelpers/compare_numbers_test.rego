package pshelpers

test_limit_equals {
  limit = 20
  set = [10, 12, 15, 20]
  false == isLimitExceeded(limit, set)
}

test_limit_is_exceeded {
  limit = 20
  set = [10, 12, 15, 21]
  true == isLimitExceeded(limit, set)
}

test_limit_not_exceeded {
  limit = 20
  set = [10, 12, 15]
  false == isLimitExceeded(limit, set)
}

test_multiple_limits_exceeded {
  limit = 20
  set = [21, 22, 1]
  true == isLimitExceeded(limit, set) 
}
