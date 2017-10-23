#!/bin/bash -e


FIRST_START_DONE="/var/lib/openldap/slapd-first-start-done"

if [ ! -e "$FIRST_START_DONE" ]; then

	CONF=/etc/openldap/slapd.conf
	LDAP_ADMIN_PASSWORD_ENCRYPTED=$(slappasswd -s $LDAP_ADMIN_PASSWORD)

	sed -i 's|^include		/etc/openldap/schema/core.schema|&\
include     /etc/openldap/schema/cosine.schema\
include     /etc/openldap/schema/dyngroup.schema\
include     /etc/openldap/schema/inetorgperson.schema\
include     /etc/openldap/schema/staffio.schema\
include     /etc/openldap/schema/misc.schema\
include     /etc/openldap/schema/nis.schema\
|g' $CONF

	sed -i 's|^# rootdn can always read and write EVERYTHING!|&\n\
access to *\
    by self write\
    by users read\
    by anonymous auth\
    |g' $CONF

	sed -i "s|^suffix	.*|suffix 		\"${LDAP_BASE_DN}\"|g" $CONF
	sed -i "s|^rootdn	.*|rootdn 		\"cn=${LDAP_ADMIN_NAME},${LDAP_BASE_DN}\"|g" $CONF
	sed -i "s|^rootpw	.*|rootpw 		\"${LDAP_ADMIN_PASSWORD_ENCRYPTED}\"|g" $CONF

	echo 'index   uid         eq' >> $CONF
	echo 'index   cn          eq' >> $CONF
	echo 'index   sn          eq' >> $CONF
	echo 'index   mail        eq' >> $CONF

	touch $FIRST_START_DONE
fi

slapd -h "ldap://$HOSTNAME ldaps://$HOSTNAME ldapi:///" -u ldap -d ${LDAP_LOG_LEVEL}
