$script = <<EOF
set -e

sudo apt-get install build-essential git-core scons -y

if [ -d bloomd ]; then
  cd bloomd && git pull && cd ..
else
  git clone https://github.com/armon/bloomd
fi

pushd bloomd
scons
mv bloomd /usr/local/bin/bloomd
popd
echo -e "[bloomd]\ndata_dir=/home/vagrant\nflush_interval=60000\nport=8673\nworkers=1\nlog_level=debug" > bloomd.conf

ps elf | grep -i bloomd | awk '{print "kill -9 "$2}' |sh

bloomd -f /home/vagrant/bloomd.conf &
echo "Done."
EOF

Vagrant.configure("2") do |config|
  config.vm.box = "hashicorp/precise64"
  config.vm.provision :shell, :inline => $script
  config.vm.network :public_network, :bridge => "en0: Wi-Fi (AirPort)"
end
