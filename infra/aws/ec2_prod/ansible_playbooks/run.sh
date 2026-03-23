instance_id=$(terraform -chdir=../  output -raw instance_id)
ansible_ssm_bucket_name=$(terraform -chdir=../  output -raw ansible_ssm_bucket_name)
echo $instance_id
echo $ansible_ssm_bucket_name



ansible-playbook -v -i inventory.yaml  main.yaml -e "instance_id=$instance_id ansible_ssm_bucket_name=$ansible_ssm_bucket_name"
