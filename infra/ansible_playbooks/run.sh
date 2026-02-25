my_host_1_ip=$(terraform output -raw my_host_01_public_ip)
ansible myhosts -m ping -i inventory.yaml -e "target_ip=$my_host_1_ip ssh_private_key_path=/home/nam/.ssh/playlist_manager"
