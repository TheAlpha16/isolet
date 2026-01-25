#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Detected macOS - using brew services"
    
    sudo brew services stop nginx
    
    NGINX_CONF_DIR=$(brew --prefix)/etc/nginx
    sudo mkdir -p "$NGINX_CONF_DIR/servers"
    sudo cp ./isolet.conf "$NGINX_CONF_DIR/servers/isolet.conf"
    
    if sudo nginx -t; then
        echo "Configuration test passed"
        sudo brew services start nginx
        
        if pgrep -q nginx; then
            echo "✓ nginx is running successfully"
        else
            echo "✗ nginx process not found. Checking service status..."
            sudo brew services info nginx
            echo ""
            echo "Check error log: tail -f /opt/homebrew/var/log/nginx/error.log"
            exit 1
        fi
    else
        echo "✗ Configuration test failed"
        exit 1
    fi
    
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "Detected Linux - using systemctl"
    
    sudo systemctl stop nginx
    
    sudo cp ./isolet.conf /etc/nginx/sites-available/isolet.conf
    sudo rm -f /etc/nginx/sites-enabled/isolet.conf
    sudo ln -s /etc/nginx/sites-available/isolet.conf /etc/nginx/sites-enabled/isolet.conf
    
    sudo systemctl start nginx
    
    sudo nginx -t
    
else
    echo "Unsupported operating system: $OSTYPE"
    exit 1
fi