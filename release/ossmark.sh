#!/bin/bash

set -e

LOG() { echo "[      ]" "$@"; }
OK() { echo "[  OK  ]" "$@"; }

install_dir="$HOME"/.ossmark

uninstall() {
    LOG "Uninstall from '$install_dir'"

    if [ "$(systemctl list-unit-files | grep -wc ossmark-article)" -ne 0 ]; then
        if [ "$(systemctl is-active ossmark-article)" == "active" ]; then
            systemctl stop ossmark-article
        fi
        systemctl disable ossmark-article
        systemctl daemon-reload
    fi

    rm -f /usr/lib/systemd/system/ossmark-article.service
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
    cp service/ossmark-article.service /usr/lib/systemd/system/ossmark-article.service
    sed -i "s@TODO_OSSMARK_DIR@$install_dir@" /usr/lib/systemd/system/ossmark-article.service
    sed -i "s@TODO_OSSMARK_USER@$user@" /usr/lib/systemd/system/ossmark-article.service
    sed -i "s@TODO_OSSMARK_GROUP@$group@" /usr/lib/systemd/system/ossmark-article.service
    systemctl daemon-reload
    systemctl enable ossmark-article
    systemctl start ossmark-article

    # shellcheck disable=SC2016
    if [ "$(grep -Fwc 'export PATH="$HOME/.ossmark/bin:$PATH"' "$HOME"/.bashrc)" -eq 0 ]; then
        echo 'export PATH="$HOME/.ossmark/bin:$PATH"' >>"$HOME"/.bashrc
        # shellcheck source=/dev/null
        source "$HOME"/.bashrc
    fi

    OK "Install success"
}

update() {
    LOG "Update to '$install_dir'"

    LOG "Update binaries"
    cp bin/ossmark-sync "$install_dir"/bin/ossmark-sync

    LOG "Update service"
    systemctl stop ossmark-article
    cp bin/ossmark-article "$install_dir"/bin/ossmark-article
    systemctl start ossmark-article

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
