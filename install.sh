#! /bin/bash

installCapitalFlows() {
    mkdir -p /etc/flows
    curl -Lo- https://github.com/sunshineplan/capital-flows/releases/download/v1.0/release.tar.gz | tar zxC /etc/flows
    cd /etc/flows
    chmod +x flows
}

configCapitalFlows() {
    read -p 'Please enter metadata server: ' server
    read -p 'Please enter VerifyHeader header: ' header
    read -p 'Please enter VerifyHeader value: ' value
    read -p 'Please enter update URL: ' update
    sed "s,\$server,$server," /etc/flows/config.ini.default > /etc/flows/config.ini
    sed -i "s/\$header/$header/" /etc/flows/config.ini
    sed -i "s/\$value/$value/" /etc/flows/config.ini
    sed -i "s,\$update,$update," /etc/flows/config.ini
    ./flows install || exit 1
    service flows start
}

main() {
    installCapitalFlows
    configCapitalFlows
}

main
