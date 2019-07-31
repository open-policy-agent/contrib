package main

const installationHelp = `
****************************| Installation Guide |**********************************

iptables is pre-installed in most linux distributions.But if it doesn't then you 
can use following instruction to install it.

Ubuntu/Debian:
--------------

Install iptables using apt-get package manager using the following command.
	
	$ sudo apt-get install iptables

Alpine:
-------

Install iptables using apk package manager using the following command.

	$ sudo apk add iptables

CentOS/RHEL 7:
--------------

Install iptables service using yum package manager using the following command.

	$ sudo yum install iptables-services

	After installing enable iptables service and start using below commands.

	$ sudo systemctl enable iptables
	$ sudo systemctl start iptables

Now check the iptables service status using below command.

	$ sudo systemctl status iptables
`