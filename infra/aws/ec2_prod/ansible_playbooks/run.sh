target_ip=$(terraform -chdir=../  output -raw my_host_01_public_ip)
# echo $target_ip



ansible-playbook -i inventory.yaml  main.yaml -e "target_ip=$target_ip ssh_private_key_path=/home/nam/.ssh/playlist_manager"
