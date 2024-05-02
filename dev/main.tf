terraform {
  required_providers {
    equinix = {
      source = "equinix/equinix"
      version = "~> 1.36.0"
    }
  }
}

# --- RESOURCES ---

resource "equinix_metal_project" "dev_project" {
  name = "carbonaut-dev"
}

resource "equinix_metal_project_ssh_key" "dev_key" {
  name       = "dev"
  public_key = var.public_key
  project_id = equinix_metal_project.dev_project.id
}

resource "equinix_metal_device" "development_server" {
  hostname         = "tf.coreos2"
  plan             = "c3.small.x86"
  metro            = "fr"
  operating_system = "debian_12"
  project_ssh_key_ids = [equinix_metal_project_ssh_key.dev_key.id]
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.dev_project.id
}

# --- VARIABLES ---

variable "public_key" {
  description = "Public SSH key for project access"
  type        = string
  sensitive   = true
}

# --- EXPORTS ---

output "device_public_ip" {
  value       = equinix_metal_device.development_server.access_public_ipv4
  description = "The public IP address of the web server."
}