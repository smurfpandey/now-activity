#!/bin/bash
# Script to download and setup now-activity service
# https://github.com/smurfpandey/now-activity

APP_NAME=now-activity
WORK_DIR=/var/www/now-activity
NGINX_CONF_NAME=now.smurfpandey.me.conf
SYSTEMD_SERVICE_NAME=now-activity.service

# Create directory if it does not exists
mkdir -p $WORK_DIR

# download latest release from Github
wget -c https://github.com/smurfpandey/$APP_NAME/releases/latest/download/$APP_NAME.tar.gz -O - | tar xvz -C $WORK_DIR/

# set correct permissions for binary
chmod u+x $WORK_DIR/$APP_NAME

# setup nginx virtual host
ln -s $WORK_DIR/$NGINX_CONF_NAME /etc/nginx/sites-available/$NGINX_CONF_NAME
ln -s /etc/nginx/sites-available/$NGINX_CONF_NAME /etc/nginx/sites-enabled/$NGINX_CONF_NAME
service nginx reload

# setup systemd service
cp $WORK_DIR/$SYSTEMD_SERVICE_NAME /lib/systemd/system/$SYSTEMD_SERVICE_NAME
systemctl start $APP_NAME
systemctl status $APP_NAME
systemctl enable $APP_NAME