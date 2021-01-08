FROM ubuntu:18.04

RUN apt-get update --fix-missing
RUN apt-get install -y \
    vim \
    openssh-server

RUN mkdir /var/run/sshd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

WORKDIR /root
RUN mkdir .ssh
COPY integration/containers/ssh/id_rsa .ssh/id_rsa
COPY integration/containers/ssh/id_rsa.pub .ssh/id_rsa.pub
COPY integration/containers/ssh/id_rsa.pub .ssh/authorized_keys

RUN chmod 700 ~/.ssh
RUN chmod 644 ~/.ssh/authorized_keys
RUN chmod 600 ~/.ssh/id_rsa
RUN chmod 644 ~/.ssh/id_rsa.pub

EXPOSE 22

RUN echo "test file ssh" >> /root/int-test
RUN chmod 777 /root/int-test

CMD ["/usr/sbin/sshd", "-D"]
