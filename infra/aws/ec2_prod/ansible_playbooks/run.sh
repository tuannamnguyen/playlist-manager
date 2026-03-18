instance_id=$(terraform -chdir=../  output -raw instance_id)
ansible_ssm_bucket_name=$(terraform -chdir=../  output -raw ansible_ssm_bucket_name)
# echo $target_ip



ansible-playbook -i inventory.yaml  main.yaml -e "instance_id=$instance_id ansible_ssm_bucket_name=$ansible_ssm_bucket_name"
