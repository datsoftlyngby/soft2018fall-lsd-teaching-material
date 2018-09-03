# An Ubuntu Linux in a VM

In case you do not have setup a dual boot setup with Linux next to your usual operating system, you may want to install it in a virtual machine (VM).

In the following is a step by step guide to get you up and running.

  1. Download and install VirtualBox ([https://www.virtualbox.org/wiki/Downloads](https://www.virtualbox.org/wiki/Downloads))
  * Download and install Vagrant ([https://www.vagrantup.com/downloads.html](https://www.vagrantup.com/downloads.html))
  * Get/open a terminal.
    - If you have Windows as OS you have to install _GitBash_, see [https://git-scm.com/downloads](https://git-scm.com/downloads). With it comes the _GitBash_ terminal, that
  you can use instead of the Windows `cmd` tool.
    -  If you are on MacOS you simply just open a new console/terminal.
  * If you did not already generate an SSH keypair, generate one (e.g., via `ssh-keygen -t rsa`)
  * See if you can run `vagrant --version` in the terminal/GitBash
  * Run `vagrant init bento/ubuntu-16.04`, which creates a file called `Vagrantfile` into the current directory
  * Open the `Vagrantfile` with the editor of your choice and make it look like in the listing below.


Before instantiating the VM customize it to your liking. For example:

  * Give it more RAM than the 2GB (`vb.memory = "2048"`). A good number is to give it your host's amount of RAM minus 1GB. You give the RAM size in MB, i.e., as multiples of 1024.  
  * Give it more CPU cores. Change `vb.cpus = "1"` to an amount that better fits your computer. For example, take the host cores minus one.
  * The IP address `192.168.33.10` is the one under which you can reach your VM. If you prefer another IP then adapt it accordingly.

```ruby
# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "bento/ubuntu-16.04"
  config.vm.network "private_network", ip: "192.168.33.10"
  config.vm.synced_folder "~/", "/host_home", type: "virtualbox"
  config.vm.provider "virtualbox" do |vb|
    vb.memory = "2048"
    vb.cpus = "1"
  end
  config.vm.provision "shell", privileged: false, inline: <<-SHELL
      sudo apt-get update
  
      sudo echo "LC_ALL=\"en_US.UTF-8\"" >> /etc/environment
      sudo locale-gen UTF-8
  
      sudo apt-get install -y git
      sudo apt-get install -y wget
      sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common
      curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
      sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
      sudo apt-get update
      sudo apt-get install -y docker-ce
      curl https://getmic.ro | sudo bash
      sudo mv ./micro /usr/local/bin
  
      echo "==================================================================="
      echo "=                             DONE                                ="
      echo "==================================================================="
      echo "To log onto the VM:"
      echo "$ vagrant ssh"  
    SHELL
end
```

  * After saving that file, start up the VM, which will take a bit on the first start up as it has to download the Ubuntu image and a bit of other software.
  ```bash
  $ vagrant up
  ```
  **OBS**: Do not type the `$` sign into your terminal. It is only there to indicate that it is a shell command.

  * To log onto the virtual machine (VM) execute
  ```bash
  $ vagrant ssh
  ```
  * Now, you should be logged onto the VM and you should see a Bash prompt similar to
  ```bash
  vagrant@vagrant:~$
  ```
  * In case you are done working on your virtual machine, you can leave it by issuing the exit command. Subsequently, you can put the virtual machine to "sleep" (just like closing the lid of your notebook) by running vagrant suspend on your host machine.
  ```bash
  vagrant@vagrant:~$ exit
  $ vagrant suspend
  ```
  * In case you want to discard this VM just run vagrant destroy from within the directory containing the `Vagrantfile`
