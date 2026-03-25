FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    samba \
    openssl \
    ldb-tools \
    krb5-user \
    smbclient \
    acl \
    attr \
    libpam-modules \
    libpam-modules-bin \
    libpam-runtime \
    gettext-base \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /bootstrap /var/lib/samba && \
    mkdir -p /var/lib/samba/private/tls

# Remove winbindd to prevent it from starting
RUN rm -f /usr/sbin/winbindd || true

COPY samba/bootstrap.ldif /bootstrap/bootstrap.ldif
COPY samba/provision.sh /bootstrap/provision.sh

RUN chmod +x /bootstrap/provision.sh

EXPOSE 53/udp 88 135 137-138/udp 139 389/udp 445 464 636/udp 3268-3269

ENTRYPOINT ["/bootstrap/provision.sh"]
