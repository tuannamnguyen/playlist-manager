#!/bin/bash

home_ip_address="$(curl -4 icanhazip.com)/32"

read -p "Enter Porkbun API key: " api_key
read -p "Enter Porkbun API secret key: " secret_key

echo "Do you want to apply or destroy the resources?"
select action in "Apply" "Destroy"; do
    case $action in
        Apply)
            terraform -chdir=../../infra/aws/ec2_prod/ apply \
                -var home_ip_address=$home_ip_address \
                -var porkbun_api_key=$api_key \
                -var porkbun_secret_key=$secret_key
            break
            ;;
        Destroy)
            terraform -chdir=../../infra/aws/ec2_prod/ destroy \
                -var home_ip_address=$home_ip_address \
                -var porkbun_api_key=$api_key \
                -var porkbun_secret_key=$secret_key
            break
            ;;
        *)
            echo "Invalid option. Please choose Apply or Destroy."
            ;;
    esac
done
