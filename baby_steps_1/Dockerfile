FROM ubuntu:24.04
WORKDIR /challenge
RUN apt-get update \
  && apt-get install -y gcc bash 
COPY HelloWorld .
RUN useradd --create-home stratizen
USER stratizen
EXPOSE 48000
