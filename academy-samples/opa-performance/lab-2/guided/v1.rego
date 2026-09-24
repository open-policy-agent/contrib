package lab
import future.keywords.in

B_abs := [ absolute | some number in input.B; absolute := abs(number) ]

common[b] {
  some a in [ absolute | some number in input.A; absolute := abs(number) ]
  some b in B_abs
  b < 50
  a == b
}


# do not edit below
# testing code

test_case_1 {
  common == {0, 1, 2}
    with input as {"A": [0, -1, 2], "B": [0, 1, -2]}
}

test_case_2 {
  common == {1, 2, 4}
    with input as {"A": [4, 3, 2, 1, 2], "B": [0, 1, -2, 4]}
}

test_case_3 {
  common == set()  # empty set
    with input as {"A": [4, 3, -7], "B": [0, 1, -2]}
}

test_all {
  test_case_1
  test_case_2
  test_case_3
}