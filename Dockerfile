
FROM --platform=linux/amd64 debian:stable-slim

ADD paperhat /usr/bin/paperhat
COPY site .en[v] /

CMD ["paperhat"]
