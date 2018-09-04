
<h1>GCP IoT Solutions</h1>
<p>Mender.io over-the-air (OTA) software updates for embedded Linux with Google
Cloud IoT Core
</p>
<h2><strong>Tutorial: Over-the-air software updates for embedded Linux with
Mender.io on Google Cloud Platform</strong></h2>
<p>
This tutorial demonstrates how to successfully deploy over-the-air (OTA)
software update solution for embedded Linux devices using Mender on Google Cloud
Platform.
</p>
<h3>Background</h3>
<h5>Mender Overview</h5>
<p>
Mender is an open source remote update manager for embedded Linux devices. The
aim of the project is to help secure connected devices by providing a robust and
easy software update process.
</p>
<p>
Some of the key features of Mender include
</p><ul>
<li>OTA update server and client
<li>Full system image update
<li>Symmetric A/B image update client
<li>Bootloader support: U-Boot and GRUB
<li>Volume formats: MBR and UEFI partitions
<li>Update commit and roll-back
<li>Build system: Yocto Project (meta-mender)
<li>Remote features: deployment server, build artifact management, device
management console</li></ul>
<p>
More information on Mender can be found <a href="https://mender.io/">here</a>
including <a href="https://mender.io/what-is-mender">What is Mender and How it
works</a>
</p>
<h5>Mender Components</h5>
<ul>
<li><strong>Mender management server</strong> - Mender Management Server, which
is the central point for deploying updates to a population of devices. Among
other things, it monitors the current software version that is installed on each
device and schedules the rollout of new releases.
<li><strong>Mender build system - </strong>Software build system generates a new
version of software for a device. The software build system is a standard
component, such as the Yocto Project. It creates build artifacts in the format
required by the target device. There will be different build artifacts for each
type of device being managed.
<li><strong>Mender Client -</strong> Each device runs a copy of the Mender
update client, which polls the Management Server from time to time to report its
status and to discover if there is a software update waiting. If there is, the
update client downloads and installs it.</li></ul>
<h3>Mender on Google Cloud Platform (GCP)</h3>
<p>
This section provides a high level architecture overview of Mender on GCP and
detailed instructions for deploying the Mender Management Server, Building
Mender artifacts and client configuration including integration with GCP IoT
Core.
</p>
<h4>Mender on GCP - High Level Architecture Diagram:</h4>
<p>
The following architecture diagram provides a high level overview of the various
components on GCP to enable OTA updates with Mender and Google Cloud IOT Core
</p>
<p>

<img src="images/Mender-on0.png" width="" alt="alt_text" title="image_tooltip">
</p>
<h4>Mender Management Server Deployment Options :</h4>
<p>
There are several options for successfully setting up Mender services with
Google Cloud Platform (GCP), this tutorial will use a minimally configured
Mender Management Production Server to test the end to end workflow:
</p><ul>
<li><a
href="https://docs.mender.io/getting-started/create-a-test-environment">Mender
Management Demo Server</a> - For quickly testing the Mender server, Mender
provides a  pre-built demo version that does not take into account
production-grade issues like security and scalability
<li><a
href="https://docs.mender.io/administration/production-installation">Mender
Management Production Server</a> - Mender Server for production environments,
and includes security and reliability aspects of Mender production
installations.
<li><a href="https://mender.io/signup">Hosted Mender Service</a>  - Hosted
Mender is a secure management service so you don't have to spend time
maintaining security, certificates, uptime, upgrades and compatibility of the
Mender server. Simply point your Mender clients to the Hosted Mender
service.</li></ul>
<h3>Mender on GCP Setup Instructions (Tutorial)</h3>
<h4>Mender Management Server on Google Cloud </h4>
<p>
Mender Management server requirements from Mender are outlined <a
href="https://docs.mender.io/getting-started/requirements">here</a> and we will
be using the base instructions as documented for setting up a production
environment and deploy on Google Cloud Platform, however this is minimally
configured and not suited for actual production use.
</p>
<h6>Before you begin</h6>
<p>
Please review the billable components of the GCP for this tutorial and
pre-requisites below.
</p>
<h6>Costs</h6>
<p>
This tutorial uses billable components of GCP, including:
</p><ul>
<li>Cloud IoT Core
<li>Cloud PubSub
<li>Google Compute Engine
<li>Google Cloud Storage
<li>Cloud Functions for Firebase
<li>Stackdriver Logging</li></ul>
<h6>Pre-requisites:</h6>
<p>
This tutorial assumes you already have a Cloud Platform account set up and have
completed the <a href="https://cloud.google.com/iot/docs/quickstart">IoT Core
quickstart</a>. You need to associate Firebase to your cloud project. Visit the
Firebase Console and choose to add a project. You can then choose to add
Firebase to an existing Cloud Project.
</p><ul>
<li>Access to <a href="https://cloud.google.com/">Google Cloud Platform
(GCP)</a>
<li>Install the <a href="https://cloud.google.com/sdk/downloads">Google Cloud
SDK </a>if you prefer to use local terminal over Google Cloud Shell
<li>Demo environment variables from Google Cloud Shell or local terminal (Please
use the below variables for every new session)
<li>A public IP address assigned and ports 443 and 9000 publicly
accessible.</li></ul>
<h6>Setup the "cloud api shell" environment (you will use several different
shell environments)</h6>
<ul>
<li>If you do not already have a development environment set up with <a
href="https://cloud.google.com/sdk/install">gcloud</a> and <a
href="https://firebase.google.com/docs/cli/">Firebase</a> tools, it is
recommended that you use <a href="https://cloud.google.com/shell/docs/">Cloud
Shell</a> for command line instructions that don't involve SSH to devices on
your local network.


<pre class="prettyprint">gcloud auth login</pre>
</li></ul>



<pre
class="prettyprint">gcloud services enable compute.googleapis.com</pre>



<pre
class="prettyprint">gcloud compute firewall-rules create mender-ota-443 --allow tcp:443
gcloud compute firewall-rules create mender-ota-9000 --allow tcp:9000</pre>



<pre
class="prettyprint">export FULL_PROJECT=$(gcloud config list project --format "value(core.project)")
export PROJECT="$(echo $FULL_PROJECT | cut -f2 -d ':')"
export CLOUD_REGION='us-central1'
# Create 2 Cloud Storage buckets you will use for updates and storage
gsutil mb -l $CLOUD_REGION gs://$PROJECT-mender-server
gsutil mb -l $CLOUD_REGION gs://$PROJECT-mender-builds</pre>
<h6>Installing Mender Management Server on GCP</h6>
<ul>
<li>Step 1: Create Google Cloud Compute Engine and runs a <a
href="https://cloud.google.com/compute/docs/startupscript">startup script</a> to
install various dependencies including Docker, as well as installing and
starting the <a
href="https://docs.mender.io/administration/production-installation">Mender
Server</a>.</li></ul>



<pre
class="prettyprint">gcloud beta compute --project $PROJECT instances create "mender-ota-demo" --zone "us-central1-c" --machine-type "n1-standard-2" --subnet "default" --maintenance-policy "MIGRATE" --scopes "https://www.googleapis.com/auth/cloud-platform" --metadata=startup-script-url=https://raw.githubusercontent.com/Kcr19/community/master/tutorials/mender_gcp_ota_demo/server/mender_server_install.sh --min-cpu-platform "Automatic" --tags "https-server" --image "ubuntu-1604-xenial-v20180814" --image-project "ubuntu-os-cloud" --boot-disk-size "10" --boot-disk-type "pd-standard" --boot-disk-device-name "mender-ota-demo"</pre>
<p>
Note: Please note the startup script will take roughly 3-5 minutes to completely
install all the prerequisites including Docker CE, Docker compose and Mender
Server.
</p><ul>
<li>Step 2 : Navigate to the Mender UI by clicking on the external IP address of
"mender-ota-demo" which can be found from the <a
href="https://console.cloud.google.com/compute">GCP console → Compute
Engine</a>. If you are using Chrome as your web browser you might get a
certificate warning and you will need to click "advanced" and "proceed". In an
actual production environment, you would provision this server with a trusted
certificate.</li></ul>
<p>

<img src="images/Mender-on1.png" width="" alt="alt_text" title="image_tooltip">
</p><ul>
<li>Once you are the Mender UI login screen please login using credentials
created in the above step which should take the Mender Dashboard.  <ul>
 <li>Username - <a href="mailto:mender@example.com">mender@example.com</a>
 <li>Password - mender_gcp_ota</li> </ul>
</li> </ul>
<p>

<img src="images/Mender-on2.png" width="" alt="alt_text" title="image_tooltip">
</p>
<p>
Congrats you just finished creating the Mender Server on Google Cloud Platform.
</p>
<h5>Hosted Mender Service</h5>
<p>
The above steps are for self-managed Open Source Mender Management Server on
GCP,Mender also provides fully managed <a href="https://mender.io/signup">Hosted
Mender service</a> .
</p>
<p>
The next section describes how to build a Yocto Project image for a raspberry
Pi3 device.
</p>
<h4>Mender Build Server on GCP</h4>
<h5>Build a Mender Yocto project OS image for Raspberry Pi3 device</h5>
<p>
These steps outline how to build a Yocto Project image for a Raspberry Pi3
device.
</p>
<p>
The <a href="https://www.yoctoproject.org/">Yocto Project</a> is an open source
collaboration project that helps developers create custom Linux-based systems
for embedded products, regardless of the hardware architecture.
</p>
<p>
The build output includes a file that can be flashed to the device storage
during initial provisioning, it has suffix ."sdimg. Additionally copies of the
same image with ".img" and ".bmap" suffix are uploaded to GCS bucket.
</p>
<p>
Yocto image builds generally take a while to complete. Generally instances with
high CPU cores and memory along with faster disks such as SSD will help speed up
the overall build process. Additionally, there are pre-built images for quick
testing available in the GCS bucket which you can download and proceed directly
to "Working with pre-built Mender Yocto Images" section
</p>
<p>
Below are the instructions to build a custom Mender Yocto image for Raspberry
Pi3 device. This image will have a number of requirements needed to communicate
with IoT Core built in.
</p>
<p>
Use the <em>"cloud api shell" environment you used earlier.</em>
</p><ol>
<li>Create a GCE instance for Mender Yocto Project OS builds:


<pre
class="prettyprint">gcloud beta compute --project $PROJECT instances create "mender-ota-build" --zone "us-central1-c" --machine-type "n1-standard-16" --subnet "default" --maintenance-policy "MIGRATE" --scopes "https://www.googleapis.com/auth/cloud-platform" --min-cpu-platform "Automatic" --tags "https-server" --image "ubuntu-1604-xenial-v20180405" --image-project "ubuntu-os-cloud" --boot-disk-size "150" --boot-disk-type=pd-ssd --boot-disk-device-name "mender-ota-build"</pre>
</li></ol>
<ol>
<li>SSH into the image and install the necessary updates required for the Yocto
Project Builds


<pre
class="prettyprint">gcloud compute --project $PROJECT ssh --zone "us-central1-c" "mender-ota-build"</pre>
</li></ol>
<ol>
<li>Install the Mender Yocto custom image build including dependencies by
downloading script from the below github repo and executing on the build server.
<p>
    This step initially builds a custom image for Raspberry Pi device as well as
mender artifact update which can be used to test the OTA feature of Mender. All
images after the completion of the build are automatically uploaded into GCS
bucket.
</p>
<p>
    Note that you are now switching to the "build server shell" environment, and
not in the "cloud api shell" environment.
</p>


<pre
class="prettyprint">export GCP_IOT_MENDER_DEMO_HOST_IP_ADDRESS=$(gcloud compute instances describe mender-ota-demo --format="value(networkInterfaces.accessConfigs[0].natIP)")</pre>
</li></ol>



<pre
class="prettyprint">wget https://raw.githubusercontent.com/Kcr19/community/master/tutorials/mender_gcp_ota_demo/image/mender-gcp-build.sh</pre>



<pre class="prettyprint">chmod +x ./mender-gcp-build.sh</pre>



<pre class="prettyprint">. ./mender-gcp-build.sh</pre>
<p>
<strong>Note:</strong> <em>The build process will generally take anywhere from
45 minutes to 60 minutes and will create custom embedded Linux image with all
the necessary packages and dependencies to be able to connect to Google Cloud
IoT Core. The output of the build will be (2) files which will be uploaded into
Google Cloud Storage Bucket. </em>
</p><ol>
<li><em>gcp-mender-demo-image-raspberrypi3.sdimg - This will be the core image
files which will be used by the client to connect to GCP IoT core and Mender
Server. Copy of the same image file with ".bmap" and ".img" are also generated
and uploaded to GCS bucket.</em>
<li><em>gcp-mender-demo-image-raspberrypi3.mender - This will be the mender
artifact file which you will upload to the mender server and deploy on the
client as part of the OTA update process.</em></li></ol>
<ol>
<li>This completes the build process, the next step is to provision the build to
<a href="https://docs.mender.io/artifacts/provisioning-a-new-device">new
device</a> (Raspberry Pi3). The build image was copied automatically to the GCS
bucket which was created earlier. </li></ol>
<p>
    Download the newly built image to your local PC where you can write the
image to SD card as outlined in the next step. Note: this is done only for the
initial provisioning of the starter image to a new device. Updates from this
point on are managed by Mender.
</p>


<pre
class="prettyprint">gsutil cp gs//$PROJECT-mender-builds/gcp-mender-demo-image-raspberrypi3.sdimg  /&lt;local PC path where you want to write the image to></pre>
<ol>
<li>Provisioning a new device (Writing the image to Raspberry Pi3 device) <ul>
 <li>Insert the SD card into the SD card slot of your local PC where you have
the "gcp-mender-demo-image-raspberrypi3.sdimg" image downloaded. If you are
using a utility such as <a href="https://etcher.io/">Etcher</a> please use the
image with ".img" suffix - "gcp-mender-demo-image-raspberrypi3.img" which can be
downloaded from the same GCS bucket as above and skip the steps below.
 <li>Unmount the drive (instructions below for Mac)


<pre
class="prettyprint">df  -h  (Use this command to determine where the drive is mounted)</pre>
</li> </ul>
</li> </ol>



<pre class="prettyprint"># on OS X:
diskutil unmountDisk /dev/disk3 (assuming /dev/disk 3 is SD card)

# on Linux:
umount &lt;mount-path></pre>
 <ul>
 <li>Command to write the image to SD card and please adjust the local path to
your .sdimg file location. Depending on the image size it may take roughly 20
minutes so please be patient until the image is completely written to the SD
card


<pre
class="prettyprint">sudo dd if=/Users/&lt;local PC path where you have your image downloaded>/gcp-mender-demo-image-raspberrypi3.sdimg of=/dev/disk2 bs=1m && sudo sync  </pre>
</li> </ul>
<h5>{Optional} Working with pre-built Mender Yocto Images</h5>
<p>
This section outlines the steps  involved in configuring and working directly
with pre-built images. Mender Yocto images are available to download from the
GCS bucket. If you have already generated a Mender Yocto build image from
previous steps please proceed directly to the Mender Client Integration section
</p><ol>
<li>Download the pre-built Mender Yocto Image from GCS bucket to local PC with
access to terminal or console  <ol>
 <li><a
href="https://storage.googleapis.com/mender-gcp-ota-images/gcp-mender-demo-image-raspberrypi3.sdimg">gcp-mender-demo-image-raspberrypi3.sdimg</a>
- Base or Core image
 <li><a
href="https://storage.googleapis.com/mender-gcp-ota-images/gcp-mender-demo-image-raspberrypi3.mender">gcp-mender-demo-image-raspberrypi3.mender</a>
- Mender artifcat<ol>
<li>Update Raspberry Pi 3 Images with Google Cloud IoT Core settings</li></ol>
</li></ol>
</li></ol>
<p>
Settings related to Google Cloud such as the REGISTRY_ID, REGION_ID and
PROJECT_ID are stored in the binary images in the file
<em>/opt/gcp/etc/config-gcp.sh. </em>The values of these settings, as shown in
this document are sample values. You will need to edit this file, and the
instructions below, to match your settings.
</p><ul>
<li>Download the mender-artifact utility:</li></ul>



<pre
class="prettyprint">wget https://d1b0l86ne08fsf.cloudfront.net/mender-artifact/master/mender-artifact</pre>



<pre class="prettyprint">chmod +x ./mender-artifact</pre>
<ul>
<li>Copy the file out of the SDIMG or MENDER binary file:</li></ul>



<pre
class="prettyprint">./mender-artifact cat gcp-mender-demo-image-raspberrypi3.sdimg:/opt/gcp/etc/gcp-config.sh > ./gcp-config.sh</pre>
<ul>
<li>Now edit the file <em>./gcp-config.sh</em> in your editor of choice and
update the values to match your Google Cloud settings.</li></ul>
<ul>
<li>Update the file in the SDIMG or MENDER binary file:


<pre
class="prettyprint">cat ./gcp-config.sh | ./mender-artifact cp gcp-mender-demo-image-raspberrypi3.sdimg:/opt/gcp/etc/gcp-config.sh</pre>
</li></ul>
<p>
</p>
<p>
Next you will configure the Mender Client including to connect to Mender
Management Server and Google IOT core with the same private/public key pair.
</p>
<h4>Mender Client Integration - GCP IoT Core and Mender Management Server</h4>
<h5>Mender Client Configuration:</h5>
<p>
This section outlines the steps to connect the Mender Client (Raspberry Pi3
device) to Google Cloud IoT Core as well as Mender Server with the same
public/private key authentication and additionally will deploy an OTA update to
the device remotely.
</p>
<p>
Key components you will use in this section are:
</p><ul>
<li>Google Cloud IoT Core
<li>Google Cloud Functions/Firebase Functions
<li>Google Cloud/Stackdriver Logging
<li>Mender Server on Google Cloud
<li>Raspberry Pi3 (Device/Client)</li></ul>
<h6>Step 1: Create registry in Google Cloud IoT Core</h6>
<p>
Using the "cloud api shell" environment:
</p>


<pre class="prettyprint">export REGISTRY_ID=mender-demo</pre>
<p>
Create a Cloud IoT Core registry and Cloud Pub/Sub topic for this tutorial that
will be used for the mender client to authenticate and send telemetry data.
</p>
<p>
Create Pub/Sub topics for "telemetry events" as well as a topic for device
lifecycle events which you will use for device preauthorization with mender
server.
</p>


<pre
class="prettyprint">gcloud pubsub topics create mender-events</pre>



<pre
class="prettyprint">gcloud pubsub topics create registration-events</pre>
<p>
Create IoT Core Registry
</p>


<pre
class="prettyprint">gcloud iot registries create $REGISTRY_ID --region=$CLOUD_REGION --event-notification-config=subfolder="",topic=mender-events</pre>
<h6>Step 2: Stackdriver logging export of Cloud IoT Device</h6>
<p>
From the Cloud console go to <a
href="https://console.cloud.google.com/logs/viewer">stackdriver logging</a> and
click on "Exports" then create "export" and
</p>
<p>
select the drop-down menu at the end of the search bar, and choose "Convert to
advanced filter" as shown in the below image
</p>
<p>

<img src="images/Mender-on3.png" width="" alt="alt_text" title="image_tooltip">
</p>
<p>
In the advanced filter text search field please enter the below filter and click
"Submit Filter".
</p>


<pre class="prettyprint">resource.type="cloudiot_device"
(protoPayload.methodName="google.cloud.iot.v1.DeviceManager.CreateDevice" OR
protoPayload.methodName="google.cloud.iot.v1.DeviceManager.UpdateDevice")</pre>
<p>

<img src="images/Mender-on4.png" width="" alt="alt_text" title="image_tooltip">
</p>
<p>
Under "Edit Export" section provide a name for the sink, select sink service as
"Cloud Pub/Sub" and Sink Destination as "registration-events" as shown below
</p>
<p>

<img src="images/Mender-on5.png" width="" alt="alt_text" title="image_tooltip">
</p>
<h6>Step 3: Deploy Firebase Functions to call Mender Preauthorization API </h6>
<p>
Note: be sure you associated Firebase with your cloud project as noted in
"Before you begin"
</p>
<p>
Deploy Firebase Functions to subscribe to Pub/Sub topic "registration-events"
which you created in the last step to <a
href="https://docs.mender.io/server-integration/preauthorizing-devices">preauthorize</a>
IoT Core Devices with the Mender Server every time a new device is created in
IoT Core
</p>
<p>
Using an existing "cloud api shell" environment clone the source repository
which contains the Firebase functions code
</p>


<pre
class="prettyprint">git clone https://github.com/Kcr19/community.git</pre>
<table>
  <tr>
  </tr>
</table>



<pre
class="prettyprint">cd community/tutorials/mender_gcp_ota_demo/auth-function/functions</pre>



<pre class="prettyprint">firebase login</pre>



<pre class="prettyprint">firebase use --add $PROJECT</pre>
<p>
Let's set the environment variables for the functions. Please replace the IP
address for mender.url with the external IP address of your mender server
</p>


<pre
class="prettyprint">export GCP_IOT_MENDER_DEMO_HOST_IP_ADDRESS=$(gcloud compute instances describe mender-ota-demo --project $PROJECT --format="value(networkInterfaces.accessConfigs[0].natIP)")
firebase functions:config:set mender.url=https://$GCP_IOT_MENDER_DEMO_HOST_IP_ADDRESS
firebase functions:config:set mender.username=mender@example.com
firebase functions:config:set mender.pw=mender_gcp_ota</pre>



<pre class="prettyprint">npm install</pre>
<table>
  <tr>
  </tr>
</table>



<pre class="prettyprint">firebase deploy --only functions</pre>
<h6>Step 4: Connect to Mender Client to extract the public key and create device
in IoT Core</h6>
<p>
Let's bring up the Raspberry Pi device and extract the public key so you can
create device in IoT Core Registry and the same private/public key pair will be
used to authorize the device in Mender Server as well.
</p>
<p>
On your local PC open terminal or console to perform the following commands.
This needs to be a shell which has access to the Raspberry Pi on your local
network, so can not be Cloud Shell. Find and add the IP address of your
Raspberry Pi device below. To locate the IP address of your Raspberry Pi device
you can invoke "nmap" command for host discovery as shown below by replacing the
subnet range with one that matches your own local network.
</p>


<pre class="prettyprint">sudo nmap -sn 192.168.1.0/24</pre>



<pre
class="prettyprint">export DEVICE_IP=&lt;your raspberry pi ip address></pre>
<p>
We use a random ID generated by the OS on first boot as our device ID in IoT
Core, this can be adapted to any potential HW based identifier such as a board
serial number, MAC address, or crypto key id:
</p>


<pre
class="prettyprint">export DEVICE_ID=$(ssh root@$DEVICE_IP /usr/share/mender/identity/mender-device-identity| head -n 1 | cut -d '=' -f 2)</pre>
<p>
Note: You will be prompted several times for the root password which is
"<strong>mender_gcp_ota</strong>"
</p>
<p>
Extract the public key from the mender-agent.pem file on the device.
</p>


<pre
class="prettyprint">ssh root@$DEVICE_IP openssl rsa -in /var/lib/mender/mender-agent.pem -pubout  -out /var/lib/mender/rsa_public.pem</pre>



<pre
class="prettyprint">scp root@$DEVICE_IP:/var/lib/mender/rsa_public.pem ./rsa_public.pem</pre>
<p>
Now create an IoT Core Device with the public key (rsa_public.pem) which you
extracted in the last step (Please make sure you are in the same directory where
you have extracted the "rsa_public.pem" file). Run the following command from
the same local console or terminal where you have ssh access to the device. You
may need to set your project in gcloud first.
</p>


<pre class="prettyprint">export REGISTRY_ID=mender-demo
export CLOUD_REGION=us-central1 # or change to an alternate region;
export PROJECT=$(gcloud config list project --format "value(core.project)")</pre>



<pre
class="prettyprint">gcloud iot devices create $DEVICE_ID --region=$CLOUD_REGION --project $PROJECT --registry=$REGISTRY_ID --public-key path=rsa_public.pem,type=RSA_PEM</pre>
<p>
Once the device is created in IoT Core, the Firebase function deployed earlier
will make REST API call to the Mender Server to preauthorize the device with the
same public key credentials used to create the device in IoT Core. Once the
preauthorization is complete the function will push a config update to Cloud IoT
Core which will configure the mender client on the device with the specific IP
address of the mender server.
</p>
<h6>Step 5: Verify the Mender Client "heartbeat" in Mender Server and Cloud IoT
Core</h6>
<p>
You can confirm the same from the Google Cloud Console as below under latest
activity.
</p>
<p>

<img src="images/Mender-on6.png" width="" alt="alt_text" title="image_tooltip">
</p>
<p>
Now open the Mender Management Server and make sure you are able to see the
device authorized and able to communicate.
</p>
<p>
Login to the Mender Server which you created part of the earlier steps - "Mender
Management Server on Google Cloud" and click on "Devices" to make sure you can
see the Raspberry Pi3 device as shown below.
</p>
<p>

<img src="images/Mender-on7.png" width="" alt="alt_text" title="image_tooltip">
</p>
<p>
This confirms the device has successfully connected to IoT core and Mender
Server with the same private/public key.
</p>
<h6>Step 6: OTA software update to Mender Client</h6>
<p>
As part of the last step let's perform Over-the-Air (OTA) update by deploying a
mender artifact from Mender Server to client.
</p>
<p>
First lets download the mender artifact
"gcp-mender-demo-image-raspberrypi3.mender" part of the Build step from the GCS
bucket and lets upload to the Mender Server under artifacts as shown below
</p>
<p>

<img src="images/Mender-on8.png" width="" alt="alt_text" title="image_tooltip">
</p>
<p>
Next you need to create a deployment and select the device which you want to
deploy the artifact. From Mender Server click "Create Deployment" and select the
target artifact as "release-2" and
</p>
<p>
group "all device" and click create deployment. Since you only have one device
currently which is "Raspberry Pi3". For various classes and types of devices
that Mender supports you can create groups and apply the target artifacts
accordingly (Eg: Raspberry Pi3, Beaglebone etc)
</p>
<p>

<img src="images/Mender-on9.png" width="" alt="alt_text" title="image_tooltip">
</p>
<p>
Deployment completion may take a moment as the update agent checks in
periodically. Progress can be monitored from the Mender Server Dashboard by
clicking on in progress deployments
</p>
<p>

<img src="images/Mender-on10.png" width="" alt="alt_text" title="image_tooltip">
</p>
<p>
Once the deployment is finished you should be able to see the deployment
successful from the Mender Dashboard and the new release of the software update
should be deployed on the mender client which can be confirmed by logging into
the device and running "mender -show-artifact" should output "release-2".
</p>
<p>
This completes the tutorial where you have successfully deployed Mender OTA
solution on Google Cloud Platform including building Mender Yocto custom
embedded OS image for Raspberry Pi device and integrated with Google Cloud IoT
Core solution.
</p>
<h4>Cleanup:</h4>
<p>
Since this tutorial uses multiple GCP components please ensure to delete the
cloud resources once you are done with the tutorial.
</p>
