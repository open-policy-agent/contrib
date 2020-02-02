package pshelpers

path_value_match(searchObject, searchPath, searchValue) = exists { 
  matches := [ foundPath |
    walk(searchObject, [searchPath, searchValue])
    foundPath = searchPath
  ]
  exists = count(matches) > 0 
}
