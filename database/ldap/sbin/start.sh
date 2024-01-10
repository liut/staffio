#!/bin/bash -e


FIRST_START_DONE="/var/lib/openldap/slapd-first-start-done"

CONF=/etc/openldap/slapd.conf

if [ ! -e "$FIRST_START_DONE" ]; then

new_includes="
include		/etc/openldap/schema/core.schema
include		/etc/openldap/schema/cosine.schema
include		/etc/openldap/schema/dyngroup.schema
include		/etc/openldap/schema/inetorgperson.schema
include		/etc/openldap/schema/staffio.schema
include	 	/etc/openldap/schema/nis.schema
"
for include in $(echo "$new_includes"); do
    if ! grep -q "$include" ${CONF}; then
        echo $include
        last_include_line=$(grep -n "^include" ${CONF} | tail -1 | cut -d':' -f1)

        sed -i "${last_include_line}a include	 	${include}" ${CONF}
    fi
done

	LDAP_ADMIN_PASSWORD_ENCRYPTED=$(slappasswd -s $LDAP_ADMIN_PASSWORD)

	sed -i 's|^# rootdn can always read and write EVERYTHING!|&\n\
access to *\
    by self write\
    by users read\
    by anonymous auth\
    |g' $CONF

	sed -i "s|^suffix	.*|suffix 		\"${LDAP_BASE_DN}\"|g" $CONF
	sed -i "s|^rootdn	.*|rootdn 		\"cn=${LDAP_ADMIN_NAME},${LDAP_BASE_DN}\"|g" $CONF
	sed -i "s|^rootpw	.*|rootpw 		\"${LDAP_ADMIN_PASSWORD_ENCRYPTED}\"|g" $CONF

	sed -i 's|^index	objectClass	eq|&\
index   uid         eq\
index   cn          eq\
index   sn          eq\
index   mail        eq\
index   mobile      eq\
index   entryCSN      eq\
index   entryUUID      eq\
|g' $CONF

	echo '' >> $CONF

	echo 'loglevel 0x40 0x100' >> $CONF

	echo '' >> $CONF

	ETC_HOSTS=$(cat /etc/hosts | sed "/$HOSTNAME/d")
	echo "0.0.0.0 $HOSTNAME" > /etc/hosts
	echo "$ETC_HOSTS" >> /etc/hosts

	cat $CONF

	[ -e /run/openldap ] || mkdir /run/openldap && chown ldap:ldap /run/openldap

	touch $FIRST_START_DONE
fi

echo "slapd starting on $HOSTNAME"
exec /usr/sbin/slapd -h "ldap://$HOSTNAME ldaps://$HOSTNAME" -u ldap -g ldap -f ${CONF} -d ${LDAP_LOG_LEVEL}
