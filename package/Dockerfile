FROM docker:19.03.8
ARG VERSION
ARG RELEASES
RUN wget -O /usr/local/bin/dapper https://${RELEASES}/dapper/${VERSION}/dapper-$(uname -s)-$(uname -m) && \
	chmod +x /usr/local/bin/dapper
