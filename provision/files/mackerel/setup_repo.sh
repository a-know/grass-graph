#!/bin/sh

set -e
sudo -k

# yum repo setup part is taken from setup-yum.sh!
sudo sh <<'SCRIPT'
  set -x

  # import GPG key
  gpgkey_path=`mktemp`
  curl -fsS -o $gpgkey_path https://mackerel.io/file/cert/GPG-KEY-mackerel-v2
  rpm --import $gpgkey_path
  rm $gpgkey_path

  # add config for mackerel yum repos
  cat >/etc/yum.repos.d/mackerel.repo <<'EOF';
[mackerel]
name=mackerel-agent
baseurl=http://yum.mackerel.io/v2/$basearch
gpgcheck=1
EOF
SCRIPT

echo '*************************************'
echo ''
echo '     Done!'
echo ''
echo '*************************************'
