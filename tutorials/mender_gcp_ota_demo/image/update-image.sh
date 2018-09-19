#! /bin/bash
wget https://d1b0l86ne08fsf.cloudfront.net/mender-artifact/master/mender-artifact
chmod +x mender-artifact
wget https://storage.googleapis.com/mender-gcp-ota-images/gcp-mender-demo-image-raspberrypi3.sdimg
sudo apt install -y parted
cat >gcp-config.sh <<'EOF'
#!/bin/sh
#
export PROJECT_ID=${PROJECT}
export REGION_ID=us-central1
export REGISTRY_ID=mender-demo
...
EOF
cat ./gcp-config.sh | ./mender-artifact cp gcp-mender-demo-image-raspberrypi3.sdimg:/data/gcp/gcp-config.sh
gsutil cp gcp-mender-demo-image-raspberrypi3.sdimg gs://$PROJECT-mender-builds/updated-demo-image-raspberrypi3.img
rm gcp-mender-demo-image-raspberrypi3.sdimg
echo "Retrieve image with gsutil cp gs://${PROJECT}-mender-builds/updated-demo-image-raspberrypi3.img ."