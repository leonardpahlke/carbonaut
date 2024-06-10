terraform {
  required_providers {
    equinix = {
      source = "equinix/equinix"
      version = "~> 1.36.0"
    }
  }
}

# --- VARIABLES ---

variable "num_projects" {
  type = number
  default = 1
  description = "Number of Equinix Metal projects to create"
}

variable "vm_count_per_project" {
  type = number
  default = 1
  description = "Number of VMs to create per project"
}

variable "public_key" {
  description = "Public SSH key for project access"
  type        = string
  sensitive   = true
}

# --- RESOURCES ---

resource "equinix_metal_project" "project" {
  for_each = { for i in range(var.num_projects) : i => format("project-%d", i) }
  name     = each.value
}

resource "equinix_metal_project_ssh_key" "project_key" {
  for_each = { for i in range(var.num_projects) : i => format("key-%d", i) }
  name     = each.value
  public_key = var.public_key
  project_id = equinix_metal_project.project[each.key].id
}

locals {
  device_map = flatten([
    for project_id in keys(equinix_metal_project.project) : [
      for vm_index in range(var.vm_count_per_project) : {
        key = format("%s-%d", project_id, vm_index)
        value = {
          hostname        = format("tf-vm-%s-%d", project_id, vm_index)
          operating_system = "debian_12"
          plan            = "c3.small.x86"
          project_id      = equinix_metal_project.project[project_id].id
          project_index   = project_id
        }
      }
    ]
  ])
}

resource "equinix_metal_device" "vm" {
  for_each = { for item in local.device_map : item.key => item.value }
  hostname        = each.value.hostname
  operating_system = each.value.operating_system
  plan            = each.value.plan
  project_id      = each.value.project_id
  billing_cycle    = "hourly"
  metro            = "fr"
  project_ssh_key_ids = [equinix_metal_project_ssh_key.project_key[each.value.project_index].id]
}

# --- EXPORTS ---

output "vm_public_ip" {
  value = { 
    for key, vm in equinix_metal_device.vm : 
    key => vm.access_public_ipv4 
  }
}
