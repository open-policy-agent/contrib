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

