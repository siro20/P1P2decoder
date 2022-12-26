#!/bin/sh
if [ -d "/tmp" ]; then
  echo $(systemctl show p1p2gateway@*|grep Id|cut -d= -f2) > /tmp/.p1p2gateway_restart
fi
systemctl stop p1p2gateway@* || /bin/true
systemctl stop p1p2gateway.service || /bin/true