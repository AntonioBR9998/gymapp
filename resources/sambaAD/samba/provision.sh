#!/usr/bin/env bash
set -e

echo "Provisioning DC"

# SAMBA FILESYSTEM (non-root user)
echo "Setting up filesystem for Samba AD DC (non-root user)"
mkdir -p /var/lib/samba/private /var/lib/samba/sysvol /var/lib/samba/netlogon

# Set permissions for samba directories
chmod 755 /var/lib/samba
chmod 700 /var/lib/samba/private 2>/dev/null || true

# Try to enable ACLs if possible, but don't fail if we can't
if mountpoint -q /var/lib/samba; then
    echo "Volume already mounted, trying to enable ACLs"
    mount -o remount,user_xattr,acl /var/lib/samba 2>/dev/null || echo "Warning: Could not enable ACLs on mounted volume"
else
    echo "Using volume for /var/lib/samba"
fi

SAMBA_DB="/var/lib/samba/private/sam.ldb"

if [ ! -f "$SAMBA_DB" ]; then
    if [ -f /etc/samba/smb.conf ]; then
        echo "Deleting smb.conf"
        rm -f /etc/samba/smb.conf
    fi

    echo "Running domain provision"
    {
        samba-tool domain provision \
        --use-rfc2307 \
        --server-role=dc \
        --dns-backend=SAMBA_INTERNAL \
        --realm="${DOMAIN}" \
        --domain="${DOMAINNETBIOS}" \
        --adminpass="${ADMINPASS}" \
        --option="server role check:inhibit = yes"
    } || {
        echo "WARNING: provision failed, but continuing"
    }

    # Replace smb.conf with our custom one that excludes winbind
    if [ -f "/var/lib/samba/private/sam.ldb" ]; then
        echo "Replacing smb.conf with custom configuration"
        envsubst < /bootstrap/smb.conf.template > /etc/samba/smb.conf

        echo "Adding bootstrap LDIF"
        ldbmodify -H /var/lib/samba/private/sam.ldb /bootstrap/bootstrap.ldif 2>/dev/null || echo 'LDAP modification failed'
        echo "Setting test_user user password"
        samba-tool user setpassword test_user --newpassword='Aa@1234' 2>/dev/null || echo 'Password setting failed'

        echo "Assigning ACL permissions to test_user on OU=Users"

        OU_DN="OU=Users,DC=gymapp,DC=es"

        # Check if already applied
        if ! samba-tool dsacl show "$OU_DN" | grep -q "test_user"; then
            TRUSTEE_SID=$(samba-tool user show test_user --attributes=objectSid \
            | grep '^objectSid:' | awk '{print $2}')

            if [ -n "$TRUSTEE_SID" ]; then
                echo "Applying AD permissions for test_user"
                samba-tool dsacl set \
                --objectdn="CN=Users,DC=gymapp,DC=es" \
                --trusteedn="CN=test_user,OU=Users,DC=gymapp,DC=es" \
                --sddl="(A;CI;RPWPCCDCLOCRRCWDWO;;;${TRUSTEE_SID})"
            else
                echo "ERROR: Could not get SID for test_user" >&2
            fi
        else
            echo "ACL already present, skipping"
        fi
    fi
else
    echo "DC already provisioned"
    # Ensure we use our custom smb.conf
    if ! grep -q "server services.*-winbind" /etc/samba/smb.conf; then
        echo "Ensuring custom smb.conf is in place"
        envsubst < /bootstrap/smb.conf.template > /etc/samba/smb.conf
    fi
fi

echo "Starting DC"

# Ensure samba can bind to all interfaces
export SAMBA_INTERFACES="0.0.0.0/0"

# Fix negative TLS certificates
TLS_DIR="/var/lib/samba/private/tls"

echo "Ensuring Samba TLS certificates exist"

if [[ ! -f "$TLS_DIR/cert.pem" ]] || [[ ! -f "$TLS_DIR/key.pem" ]] || [[ ! -f "$TLS_DIR/ca.pem" ]]; then
    echo "Creating new TLS certificates"

    openssl genrsa -out "$TLS_DIR/key.pem" 2048

    SERIAL=$(printf "%016X" $RANDOM)

    openssl req -new -x509 -days 3650 \
        -key "$TLS_DIR/key.pem" \
        -out "$TLS_DIR/cert.pem" \
        -subj "/CN=${HOSTNAME}" \
        -set_serial 0x$SERIAL \
        -addext "basicConstraints=critical,CA:TRUE,pathlen:0" \
        -addext "keyUsage=critical,keyCertSign,digitalSignature,keyEncipherment"

    cp "$TLS_DIR/cert.pem" "$TLS_DIR/ca.pem"
fi

chmod 600 "$TLS_DIR"/*.pem
chown root:root "$TLS_DIR"/*.pem

# Start samba with proper debugging
echo "Starting Samba AD DC"
exec samba -i --no-process-group --debug-stdout
