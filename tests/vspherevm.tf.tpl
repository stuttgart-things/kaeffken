module "{{ .name }}" {
  source                 = "github.com/stuttgart-things/vsphere-vm?ref=v1.7.5-2.7.0-1"
  vm_count               = {{ .vmCount }}
  vsphere_vm_name        = "{{ .name }}"
  vm_memory              = {{ .vmMemory }}
}
