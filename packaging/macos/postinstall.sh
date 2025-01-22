#!/usr/bin/env bash

# Make sure the script is being run as root
if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root" >&2
    exit 1
fi

# Check if brew is installed
if ! command -v brew &>/dev/null; then
    echo "Homebrew is not installed" >&2
    exit 1
fi

# Check if dnsmasq is installed
if ! brew list dnsmasq &>/dev/null; then
    echo "Dnsmasq is not installed" >&2
    exit 1
fi

# Check if nginx is installed
if ! brew list nginx &>/dev/null; then
    echo "Nginx is not installed" >&2
    exit 1
fi

# Check if the SUDO_USER environment variable is set
if [ -z "$SUDO_USER" ]; then
    echo "The SUDO_USER environment variable is not set" >&2
    exit 1
fi

# Get the home directory of the user that ran the script
USER_HOME="/Users/$SUDO_USER"

# Make sure the user's home directory exists
if [ ! -d "$USER_HOME" ]; then
    echo "The user's home directory does not exist" >&2
    exit 1
fi

SPINUP_SHARE_DIR="/usr/local/share/spinup"

#####################
### General setup ###
#####################

# Remove the old symlink if it exists and create a new one
rm -f /usr/local/bin/spinup
ln -s $SPINUP_SHARE_DIR/bin/spinup /usr/local/bin/spinup

# Set the correct permissions on the spinup directory
chown -R root:wheel $SPINUP_SHARE_DIR

###############
### Dnsmasq ###
###############

if [ ! -d /etc/dnsmasq.d ]; then
    # Create the dnsmasq.d directory if it doesn't exist
    mkdir /etc/dnsmasq.d
fi

# Check if the spinup.conf already exists in the dnsmasq.d directory
if [ ! -f /etc/dnsmasq.d/spinup.conf ]; then
    # Link the dnsmasq configuration file for spinup to the dnsmasq.d directory if it doesn't exist
    ln -s $SPINUP_SHARE_DIR/config/dnsmasq.conf /etc/dnsmasq.d/spinup.conf

    # Restart the dnsmasq service
    brew services restart dnsmasq
fi

#############
### Nginx ###
#############

USER_SPINUP_NGINX_DIR="$USER_HOME/.config/spinup/nginx"
SPINUP_NGINX_DIR="/usr/local/share/spinup/config/nginx"
BREW_NGINX_DIR="/opt/homebrew/etc/nginx"

# Create the user's spinup nginx directory if it doesn't exist
mkdir -p $USER_SPINUP_NGINX_DIR

if [ ! -d $SPINUP_NGINX_DIR ]; then
    # Link the user's spinup nginx directory to the spinup config directory if it doesn't exist
    ln -s $USER_SPINUP_NGINX_DIR $SPINUP_NGINX_DIR
fi

# Add the user's spinup nginx directory to the nginx configuration
if ! grep -q "include ${SPINUP_NGINX_DIR}/\*.conf;" $BREW_NGINX_DIR/nginx.conf; then
    # Make a backup of the nginx configuration file
    cp $BREW_NGINX_DIR/nginx.conf $BREW_NGINX_DIR/nginx.conf.bak

    # Add the include directive to the bottom of the http section in the nginx configuration file
    sed -i '/http {/!b; :a; N; /}/!ba; s/}/\tinclude '"${SPINUP_NGINX_DIR//\//\\/}"'\/\*.conf;\n}/' $BREW_NGINX_DIR/nginx.conf
fi

# Restart the nginx service
brew services restart nginx

###############
### Sudoers ###
###############

# Add an include for the sudoers.d directory to the sudoers file if it doesn't already exist
if ! grep -q "^@includedir /etc/sudoers.d[[:space:]]*$" /etc/sudoers; then
    # Make a backup of the sudoers file
    cp /etc/sudoers /etc/sudoers.bak

    # Add the include directive to the sudoers file
    echo "@includedir /etc/sudoers.d" | visudo -c -f - &>/dev/null && echo "@includedir /etc/sudoers.d" | EDITOR='tee -a' visudo &>/dev/null
fi

# Move the spinup sudoers file to the sudoers.d directory
mv $SPINUP_SHARE_DIR/config/sudoers /etc/sudoers.d/spinup
chmod 440 /etc/sudoers.d/spinup
chown root:wheel /etc/sudoers.d/spinup

# Set the correct permissions on the user's spinup .config directory
chown -R $SUDO_USER:staff $USER_HOME/.config/spinup
