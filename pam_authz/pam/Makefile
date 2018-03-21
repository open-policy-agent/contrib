MODULE := pam_authz
CLEAN_OBJECTS := -rm -f *.o *.gch

all: module

module: http flag display authz log json pull sysinfo
	gcc -fPIC -fno-stack-protector -c pam.c -o ${MODULE}.o
	gcc -shared -lcurl -ljansson -Wl,-soname,${MODULE}.so -o ${MODULE}.so ${MODULE}.o $(foreach var,$?,$(var).o)

	${CLEAN_OBJECTS}

http:
	gcc -fPIC -c $@.h $@.c -lcurl

flag:
	gcc -fPIC -c $@.h $@.c 

display:
	gcc -fPIC -c $@.h $@.c

authz:
	gcc -fPIC -c $@.h $@.c	

log:
	gcc -fPIC -c $@.h $@.c

json:
	gcc -fPIC -c $@.h $@.c

pull:
	gcc -fPIC -c $@.h $@.c

sysinfo:
	gcc -fPIC -c $@.h $@.c

clean:
	-rm -f ${MODULE}.so
	${CLEAN_OBJECTS}

.PHONY: module clean
