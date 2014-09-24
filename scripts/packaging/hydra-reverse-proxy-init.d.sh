#!/bin/bash

# Hydra Reverse Proxy - Startup script for Hydra Reverse Proxy

# chkconfig: 35 99 15
# description: Reverse proxy service to applications balanced by hydra 
# processname: hydra-reverse-proxy
# Default-Start: 2 3 4 5
# Default-Stop: 0 1 6
# config: 
# pidfile: /var/run/hydra-reverse-proxy.pid

DISTRO_INFO=$(cat /proc/version)

if [[ $(echo $DISTRO_INFO | grep 'Debian\|Ubuntu') == "" ]]; then
	. /etc/rc.d/init.d/functions
fi

APP_NAME=hydra-reverse-proxy
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin
DAEMON=/usr/local/hydra-reverse-proxy
DAEMON_ARGS=""
RUNDIR=/usr/local
PID_DIR=/var/run
PID_NAME=$APP_NAME.pid
PID_FILE=$PID_DIR/$PID_NAME
LOCK_FILE=/var/lock/subsys/${APP_NAME}
USER=root
GROUP=root

rh_status() {
    status $PID_DIR/$APP_NAME $DAEMON
    RETVAL=$?
    return $RETVAL
}

case "$1" in
start)
  if [ -f $PID_FILE ]
  then
    echo Already running with PID `cat $PID_FILE`
  else
    cd $RUNDIR
    rm -rf conf log snapshot
    if [[ $(echo $DISTRO_INFO | grep 'Debian\|Ubuntu') != "" ]]; then
      if start-stop-daemon --start --pidfile $PID_FILE --chdir $RUNDIR --background --make-pidfile --exec $DAEMON -- $DAEMON_ARGS &> /var/log/${APP_NAME}.log
      then
        echo ok
      else
        echo start failed
      fi
    else  
      $DAEMON $DAEMON_ARGS &> /var/log/${APP_NAME}.log &
      RETVAL=$?
      if [ $RETVAL -eq 0 ]
      then
        echo [OK]
        PID=$!
        touch $LOCK_FILE
        echo $PID > $PID_FILE
      else
        echo [ERROR]
      fi
    fi
  fi
  ;;
stop)
  if [ -f $PID_FILE ]
  then
    kill -9 `cat $PID_FILE`
    rm -f $PID_FILE
  else
    echo $PID_FILE not found
  fi
  ;;
restart)
  ${0} stop
  ${0} start
  ;;
status)
  rh_status
  ;;
*)
  echo "Usage: /etc/init.d/$NAME {start|stop|restart}"
  exit 1
  ;;
esac

exit 0
