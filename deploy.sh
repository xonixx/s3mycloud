#!/usr/bin/env bash

APP=s3mycloud
USER=apps1
SERV=prod.cmlteam.com

echo
echo "BUILD..."
echo

./mvnw clean package -DskipTests

echo
echo "DEPLOY..."
echo

scp $APP.conf target/$APP.jar $USER@$SERV:~/

echo
echo "RESTART..."
echo

ssh $USER@$SERV "
if [ ! -f /etc/init.d/$APP ]
then
    sudo ln -s /home/$USER/$APP.jar /etc/init.d/$APP
    sudo update-rc.d $APP defaults 99
fi
sudo /etc/init.d/$APP restart
sleep 20
tail -n 100 /var/log/$APP.log
"