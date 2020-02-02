package main

deny_if_extensions_do_not_match[msg] {
    pasExtensions := [ key |
        input[_]["resource-config"]["router"]["additional_vm_extensions"][i]
        key := input[_]["resource-config"]["router"]["additional_vm_extensions"][i]
    ]
    cfExtensions := [ key | 
        input[_]["resource-config"]["router"]["additional_vm_extensions"][i]
        key := input[_]["resource-config"]["router"]["additional_vm_extensions"][i]
    ]
    pasExtensions != cfExtensions
    msg = sprintf("Expected vm extension to match instead got and %v for PAS and %v", [pasExtensions, cfExtensions])
}

#this currently doesn't work because it cannot 