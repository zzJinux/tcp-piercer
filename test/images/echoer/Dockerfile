FROM alpine
ENV LISTEN_PORT=8080
ENV RUN_COMMAND='xargs -n1 echo ECHO: '
RUN apk add nmap-ncat
CMD ncat -vk -l ${LISTEN_PORT} -c \
    'echo $NCAT_REMOTE_ADDR $NCAT_REMOTE_PORT $NCAT_LOCAL_ADDR $NCAT_LOCAL_PORT; $RUN_COMMAND'
