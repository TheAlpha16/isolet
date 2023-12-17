FROM alpine:latest

RUN apk -yq update
RUN apk add --no-cache openssh
RUN echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config
RUN ssh-keygen -A

COPY ./run.sh /root
RUN chmod +x /root/run.sh

EXPOSE 22
CMD ["/root/run.sh"]