#!/bin/bash

set -e

LOG() { echo "[      ]" "$@"; }
OK() { echo "[  OK  ]" "$@"; }

install_dir="$HOME"/.ossmark

uninstall() {
    LOG "Uninstall from '$install_dir'"

    if [ "$(systemctl list-unit-files | grep -wc ossmark)" -ne 0 ]; then
        if [ "$(systemctl is-active ossmark)" == "active" ]; then
            systemctl stop ossmark
        fi
        systemctl disable ossmark
        systemctl daemon-reload
    fi

    rm -f /usr/lib/systemd/system/ossmark.service
    rm -rf "$install_dir"

    OK "Uninstall success"
}

install() {
    LOG "Install to '$install_dir'"

    uninstall

    user=$(id -un)
    group=$(id -gn)

    mkdir -p "$install_dir"
    cp -r {bin,conf} "$install_dir"
    cp service/ossmark.service /usr/lib/systemd/system/ossmark.service
    sed -i "s@TODO_OSSMARK_DIR@$install_dir@" /usr/lib/systemd/system/ossmark.service
    sed -i "s@TODO_OSSMARK_USER@$user@" /usr/lib/systemd/system/ossmark.service
    sed -i "s@TODO_OSSMARK_GROUP@$group@" /usr/lib/systemd/system/ossmark.service
    systemctl daemon-reload
    systemctl enable ossmark
    systemctl start ossmark

    # shellcheck disable=SC2016
    if [ "$(grep -Fwc 'export PATH="$HOME/.ossmark/bin:$PATH"' "$HOME"/.bashrc)" -eq 0 ]; then
        echo 'export PATH="$HOME/.ossmark/bin:$PATH' >>"$HOME"/.bashrc
        # shellcheck source=/dev/null
        source "$HOME"/.bashrc
    fi

    OK "Install success"
}

update() {
    LOG "Update to '$install_dir'"

    systemctl stop ossmark
    cp bin/ossmark "$install_dir"/bin/ossmark
    systemctl start ossmark

    OK "Update success"
}

usage() {
    LOG "Usage: bash ossmark.sh  [install | uninstall | update]"
    exit 1
}

case "$1" in
"install")
    install
    ;;
"uninstall")
    uninstall
    ;;
"update")
    update
    ;;
*)
    usage
    ;;
esac
