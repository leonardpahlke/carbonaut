# Development Setup

This folder contains a test environment which can be used to test Carbonaut and plugins.
The test environment is deployed to [Equinix](https://www.equinix.com/) infrastructure which offers access to bare metal instances.

## Getting started:

1. To get started make sure to install [opentofu](https://opentofu.org/).
2. Next, you need to create an account in [Equinix](https://www.equinix.com/)
3. After that you naviagte to your console, create a new project and generate an API key ([docs](https://deploy.equinix.com/developers/docs/metal/identity-access-management/users/#api-keys)). A project is required to host the project, a key to access Equinix and start machines.
4. The API key needs to be exported `export METAL_AUTH_TOKEN=XXxYOURxTOKENxHERExXX`. This can be added optionally to the `~/.zshrc` file or `~/.bashrc` file so its getting loaded by default.
5. Next, we need to reference a local public SSH key so its possible to access the VM from your local machine. YOu can create a new SSH key if you prefer or use an existing one ([docs](https://deploy.equinix.com/developers/docs/metal/identity-access-management/ssh-keys/)).
   1. The Makefile uses `id_equinix_carbonaut_ed25519` as default keyname.
6. If a key was created or you choose to reuse an existing key, we are ready to get the setup deployed. Run the script
