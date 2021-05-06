#!/bin/bash

source "$PLUGIN_AVAILABLE_PATH/letsencrypt/functions"

set +e

letsencrypt_is_active $2 > /dev/null

case $? in

"0")
	# Enabled
	echo "true"
	;;

"1")
	# Something went wrong
	echo "App not found"
	exit 1
	;;

"2")
	# Not enabled
	echo "false"
	;;

*)
	# Wheeeeee
	echo "An unknown error occured."
	exit 1
	;;
esac