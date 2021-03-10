
#!/bin/sh
#
#chkconfig: 2345 80 90
# Simple Redis init.d script conceived to work on Linux systems
# as it does use of the /proc filesystem.
# redis 启动停止脚本,该redis需要密码登录，如没有密码，去掉stop函数里的 -a 

REDISPORT=6379
EXEC=./bin/redis-server
CLIEXEC=./bin/redis-cli
 
PIDFILE=./db/redis_${REDISPORT}.pid
CONF=" ./etc/redis.conf"
 
case "$1" in
    start)
        if [ -f $PIDFILE ]
        then
                echo "$PIDFILE exists, process is already running or crashed"
        else
                echo "Starting Redis server..."
                $EXEC $CONF &
        fi
        ;;
    stop)
        if [ ! -f $PIDFILE ]
        then
                echo "$PIDFILE does not exist, process is not running"
        else
                PID=$(cat $PIDFILE)
                echo "Stopping ..."
                #$CLIEXEC -a "2018" -p $REDISPORT shutdown
                $CLIEXEC   -p $REDISPORT shutdown
                while [ -x /proc/${PID} ]
                do
                    echo "Waiting for Redis to shutdown ..."
                    sleep 1
                done
                echo "Redis stopped"
        fi
        ;;
    *)
        echo "Please use start or stop as first argument"
        ;;
esac
