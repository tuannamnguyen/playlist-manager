# TuneBridge

## Description: To be updated

## Deployment guides

1. Local deployment:

    - Clone this repo
    - [Install mise](https://mise.jdx.dev/getting-started.html)
    - Run ```mise install``` to install needed tools for development
    - Set value for DOTENVX_PRIVATE_KEY environment variable: ```export DOTENVX_PRIVATE_KEY=...```
    - Run ```docker compose --profile test_minimal up``` to start services
    - Run the ```install_migrate.sh``` script and ```migrate_up.sh``` script to apply migration to Postgres DB

2. EC2 deployment:

    - Provision services in AWS with Terraform in ```infra/aws/ec2_prod``` folder
    - Run Ansible playbooks in ```infra/aws/ansible_playbooks``` to setup EC2 servers
