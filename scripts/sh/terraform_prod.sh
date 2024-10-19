#!/bin/bash

cd ../../cmd/api/
dotenv_key=$(npx dotenv-vault keys production)
cd - > /dev/null

read -p "Enter the image tag (e.g., sha-123456): " image_tag
read -p "Enter the DB root password: " db_root_password

echo "Do you want to apply or destroy the resources?"
select action in "Apply" "Destroy"; do
    case $action in
        Apply)
            terraform -chdir=../../terraform/prod/ apply \
                -var image_tag=$image_tag \
                -var dotenv_key=$dotenv_key \
                -var db_root_password=$db_root_password
            break
            ;;
        Destroy)
            terraform -chdir=../../terraform/prod/ destroy \
                -var image_tag=$image_tag \
                -var dotenv_key=$dotenv_key \
                -var db_root_password=$db_root_password
            break
            ;;
        *)
            echo "Invalid option. Please choose Apply or Destroy."
            ;;
    esac
done
