SUMMARY = "Google Cloud Platform Python Samples"
LICENSE = "Apache-2.0"
LIC_FILES_CHKSUM = "file://LICENSE;md5=86d3f3a95c324c9479bd8986968f4327"

FILESEXTRAPATHS_prepend := "${THISDIR}/files:"
SRC_URI = " \
     git://github.com/GoogleCloudPlatform/python-docs-samples;branch=master \
     file://start-mqtt-example.sh \
"
SRCREV = "47a39ccedf3cfdaa7825269800af7bf1294cc79c"

S = "${WORKDIR}/git"
B = "${WORKDIR}/build"

PACKAGES += "${PN}-mqtt-example"

FILES_${PN}-mqtt-example = " \
    /opt/gcp${bindir}/cloudiot_mqtt_example.py \
    /opt/gcp${bindir}/start-mqtt-example.sh \
"

do_install() {
    install -m 0700 -d ${D}/opt/gcp${bindir}
    install -m 0700 ${S}/iot/api-client/mqtt_example/cloudiot_mqtt_example.py ${D}/opt/gcp${bindir}
    install -m 0700 ${WORKDIR}/start-mqtt-example.sh ${D}/opt/gcp${bindir}
}

RDEPENDS_${PN} += "bash python gcp-config"

inherit deploy

do_deploy() {
    if [ -z "${PROJECT_ID}" ]; then
       echo "Error. PROJECT_ID bitbake/shell variable unset." >&2
       exit 1
    fi
    if [ -z "${REGION_ID}" ]; then
       echo "Error. REGION_ID bitbake/shell variable unset." >&2
       exit 1
    fi
    if [ -z "${REGISTRY_ID}" ]; then
       echo "Error. REGISTRY_ID bitbake/shell variable unset." >&2
       exit 1
    fi
    
    install -d ${DEPLOYDIR}/persist/gcp
    install -m 0700 ${S}/iot/api-client/mqtt_example/resources/roots.pem ${DEPLOYDIR}/persist/gcp
}
addtask do_deploy after do_install before do_package
