FROM lchaia/monerod

USER root
RUN apt-get install -y curl

USER monero
RUN mkdir /home/monero/wallets

RUN echo "curl $@" >> '/home/monero/script.sh' \
    && chmod +x /home/monero/script.sh