sudo systemctl stop nginx
sudo cp ./isolet.conf /etc/nginx/sites-available/isolet.conf
sudo rm /etc/nginx/sites-enabled/isolet.conf
sudo ln -s /etc/nginx/sites-available/isolet.conf /etc/nginx/sites-enabled/isolet.conf
sudo systemctl start nginx

sudo nginx -t