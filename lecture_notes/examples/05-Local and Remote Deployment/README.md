# Using Vagrant for provisioning of local and remote development environments

## What do I need?

You have to have Vagrant installed. Follow the instructions on https://www.vagrantup.com/intro/getting-started/install.html and download a binary for your system from https://www.vagrantup.com/downloads.html

In case you are working on Windows, the remainder of this document assumes that you are working from a Git Bash shell. Get it from here: https://git-for-windows.github.io

Furthermore, it is assumed that you have an SSH key-pair readily available and that they are stored in the `~/.ssh` directory and that the private key there is called `id_rsa`. That is, the private key is available as `~/.ssh/id_rsa` and the public key is `~/.ssh/id_rsa.pub`. It is further assumed that the public key is the one that you registered at DigitalOcean.

If you do not have already generated an SSH key-pair you can do so, e.g., via `ssh-keygen -t rsa`.


### For local development


On a local machine, this `Vagrantfile` will generate a VirtualBox VM. Consequently, you have to have VirtualBox (https://www.virtualbox.org/wiki/Downloads) installed.

On PCs, with an Intel CPU, the Intel VT-x feature has to be enabled via a BIOS or UEFI firmware setting. That is, to make VirtualBox work properly you have to restart your computer, enter the BIOS/UEFI and enable Intel VT-x.

Now you can open your shell (Git Bash on Windows) and create the VM:

```bash
$ vagrant up loc
```

If needed, you can log onto the VM with:

```bash
$ vagrant ssh loc
```




You can access the Tomcat web console via http://localhost:6080 and the MySQL database via `localhost:3306`.



### For remote development

If you do not already have the DO Vagrant Provider installed do this now (only once), via:

```bash
$ vagrant plugin install vagrant-digitalocean
```

For remote development, this `Vagrantfile` will generate the second smallest droplet (1gb in Frankfurt) on DigitalOcean (DO). That is, you need to be a registered DO user. Furthermore, you need to have registered your public SSH key there and you have to have an API token from DO to use their API. In case you do not have the latter get one by following the steps described in this tutorial: https://www.digitalocean.com/community/tutorials/how-to-use-the-digitalocean-api-v2.

Additionally, the it is assumed that you have two environment variables set on you host machine. One accessible via `$DIGITAL_OCEAN_TOKEN` and the other one via `$SSH_KEY_NAME`. If you do not have them set already, add the following two lines to your `~/.bash_profile`:

```bash
export DIGITAL_OCEAN_TOKEN="<YOUR_DO_TOKEN>"
export SSH_KEY_NAME="<NAME OF SSH KEY AT DO>"
```

Now, either reload your new environment settings via your `source ~/.bash_profile` or open a new terminal.


After all the above setup and configuration, you can now run:

```bash
$ vagrant up prod --provider=digital_ocean
```

which will instantiate a new droplet at DO. **OBS** from here on you are paying them until you destroy the droplet, see below.

### I am done...

In case you are done working and you want to delete both, the local VM and the remote droplet run the following command:

```bash
$ vagrant destroy
```

To destroy only the remote droplet run:

```bash
$ vagrant destroy prod
```

And vice-versa, to destroy only the local VM run:

```bash
$ vagrant destroy loc
```


You can suspend droplets and VMs with:

```bash
$ vagrant suspend
```

```bash
$ vagrant suspend loc
```

```bash
$ vagrant suspend prod
```


To get an overview of provisioned local VMs and remote droplets you can run `$ vagrant global-status`.