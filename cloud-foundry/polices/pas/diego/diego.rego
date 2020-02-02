package main

find_vm_type(dataObj, service) = vmType {
    find_VMType := [type |
        dataObj.resources[i].identifier == service
        type := dataObj.resources[i]
    ]
    vmType := {
        "resource": find_VMType,
        "present": count(find_VMType) > 0
    }
}
# deny_if_memory_over_allocated[msg] {
#     false == true

#     msg = sprintf("%v", [vmType])
# }

