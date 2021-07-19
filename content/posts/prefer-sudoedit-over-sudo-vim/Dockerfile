FROM ubuntu:latest

RUN export EDITOR=vim
RUN apt-get update -y && apt-get install sudo -y
RUN useradd -d /home/alice -m -s /bin/bash alice
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
USER alice
CMD ["bash"]
