$script = <<EOF
set -e

sudo apt-get install build-essential git-core scons -y

if [ -d bloomd ]; then
  cd bloomd && git pull && cd ..
else
  git clone https://github.com/armon/bloomd
fi

cd bloomd
scons
mv bloomd /usr/local/bin/bloomd
cd ..
echo > /home/vagrant/bloomd.conf <<EOE
[bloomd]
data_dir=/home/vagrant
flush_interval=60000
port=8673
workers=1
log_level=debug"
EOE

ps elf | grep -i bloomd | awk '{print "kill -9 "$2}' |sh

bloomd -f /home/vagrant/bloomd.conf &
echo "Done."
EOF

Vagrant.configure("2") do |config|
  config.vm.box = "precise64"
  config.vm.provision :shell, :inline => $script
  config.vm.network :bridged, :bridge => "en0: Wi-Fi (AirPort)"
end
