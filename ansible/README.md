# Ansible Integration for ChatApp

This directory contains Ansible playbooks and configuration for automating deployment and management tasks for the ChatApp project.

## Overview

- Use Ansible to automate post-provisioning, configuration, and deployment steps after infrastructure is created with Terraform.
- Ansible can also be used to run Helm and kubectl commands on your EKS cluster.

## Structure

```
ansible/
  ansible.cfg           # Ansible configuration
  inventories/
    production          # Inventory file for production hosts
  playbooks/
    deploy.yaml         # Main deployment playbook
  group_vars/           # Group variables (optional)
  host_vars/            # Host variables (optional)
  roles/                # Custom roles (optional)
```

## Usage

1. **Configure Inventory**

   - Edit `inventories/production` and add the public IP or DNS of your deployment host (bastion or jump box with access to EKS).

2. **Set Up SSH Access**

   - Ensure your SSH key is available and the `ansible_user` matches your host's user (e.g., `ubuntu`).

3. **Run the Playbook**
   ```fish
   ansible-playbook -i inventories/production playbooks/deploy.yaml
   ```

## Playbook Tasks

- Fetch Terraform outputs
- Update kubeconfig for EKS
- Deploy services using Helm
- Apply Kubernetes manifests

## Requirements

- Ansible installed on your local machine or CI runner
- AWS CLI, kubectl, and Helm installed on the target host
- AWS credentials available as environment variables or via Ansible Vault

## References

- [Ansible Documentation](https://docs.ansible.com/)
- [AWS EKS Guide](https://docs.aws.amazon.com/eks/latest/userguide/what-is-eks.html)
- [Helm Documentation](https://helm.sh/docs/)
