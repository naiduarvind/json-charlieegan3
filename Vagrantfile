VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"

  config.ssh.forward_agent = true
  config.vm.provision "file", source: "~/.gitconfig", destination: "~/.gitconfig"

  config.vm.provision :ansible do |ansible|
    ansible.playbook = "playbook.yml"
  end
end
