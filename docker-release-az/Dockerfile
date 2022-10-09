FROM ubuntu:22.04

# ssh
ENV SSH_PASSWD "root:Docker!"
RUN apt-get update \
        && apt-get install -y --no-install-recommends dialog \
        && apt-get update \
	&& apt-get install -y --no-install-recommends openssh-server \
	&& echo "$SSH_PASSWD" | chpasswd
COPY assets/sshd_config /etc/ssh/

# App
COPY usr/ /usr/
RUN chmod 755 /usr/bin/http-redirect-geo

# Init script
COPY assets/init.sh /
RUN chmod u+x /init.sh

EXPOSE 2222 80

CMD /init.sh
