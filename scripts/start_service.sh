#!/bin/sh

systemctl daemon-reload

if [ -f "/tmp/.p1p2gateway_restart" ]; then
  systemctl restart $(cat /tmp/.p1p2gateway_restart)
  systemctl enable $(cat /tmp/.p1p2gateway_restart)
  rm /tmp/.p1p2gateway_restart
fi
systemctl restart p1p2gateway.service
systemctl enable p1p2gateway.service