#!/bin/bash

source "$PLUGIN_AVAILABLE_PATH/letsencrypt/functions"

set +e

letsencrypt_is_active $2 > /dev/null

case $? in

"0")
	# Enabled
	echo "true"
	;;

*)
	# Wheeeeee
	echo "false"
	;;
esac
