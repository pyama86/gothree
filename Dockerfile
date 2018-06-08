FROM ubuntu:latest
RUN apt-get update -qqy && apt-get install -qqy wget logrotate unzip vim rsyslog
ADD pkg/linux_amd64/gothree /usr/local/bin/gothree
ADD misc/rsyslog /etc/logrotate.d/
CMD service rsyslog start && logrotate -f /etc/logrotate.conf &&  /bin/bash
