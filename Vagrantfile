# -*- mode: ruby -*-
# vi: set ft=ruby :

# Install vagrant-disksize to allow resizing the vagrant box disk.
unless Vagrant.has_plugin?("vagrant-disksize")
    raise  Vagrant::Errors::VagrantError.new, "vagrant-disksize plugin is missing. Please install it using 'vagrant plugin install vagrant-disksize' and rerun 'vagrant up'"
end

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.
  config.vm.box = "archlinux/archlinux"
  config.vm.box_check_update = true

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # Accessing "localhost:443" will access port 443 on the guest machine.
  # NOTE: This will enable public access to the opened port
  config.vm.network "forwarded_port", guest: 443, host: 443
  config.disksize.size = '20GB'

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine and only allow access
  # via 127.0.0.1 to disable public access
  # config.vm.network "forwarded_port", guest: 80, host: 8080, host_ip: "127.0.0.1"

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Note: use `vagrant reload` after changing virtualbox parameters.
  config.vm.provider "virtualbox" do |v|
    v.cpus = 4
    v.memory = 4096
  end

  config.vm.provision "shell", inline: <<-SHELL
    pacman --noconfirm -Sy base-devel docker docker-compose git go mkcert nano nodejs npm postgresql-libs vim
    systemctl start docker.service
    systemctl enable docker.service
    usermod -aG docker vagrant
    docker run --rm hello-world
  SHELL

  config.vm.provision "shell", privileged: false, inline: "/bin/bash --login /vagrant/utils/vagrant-bootstrap"
end
