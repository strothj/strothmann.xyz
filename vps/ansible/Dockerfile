FROM ubuntu:trusty
RUN apt-get update && apt-get install -y \
  git \
  python \
  python-pip
RUN git clone git://github.com/ansible/ansible.git --recursive /app
