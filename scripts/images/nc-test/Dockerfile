FROM alpine:latest

RUN apk add netcat-openbsd

RUN echo -e '#!/bin/sh\n'\
    'while true; do echo "$FLAG $USERNAME $CTF_NAME. Done" | nc -l -p 53581; done' > /nc-start.sh

RUN chmod +x /nc-start.sh

EXPOSE 53581
CMD ["/nc-start.sh"]
