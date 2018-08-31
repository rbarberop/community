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
    /opt/gcp${sysconfdir}/roots.pem \
"

do_install() {
    install -m 0700 -d ${D}/opt/gcp${bindir}
    install -m 0700 ${S}/iot/api-client/mqtt_example/cloudiot_mqtt_example.py ${D}/opt/gcp${bindir}
    install -m 0700 ${WORKDIR}/start-mqtt-example.sh ${D}/opt/gcp${bindir}
    install -m 0700 -d ${D}/opt/gcp${sysconfdir}
    install -m 0700 ${S}/iot/api-client/mqtt_example/resources/roots.pem ${D}/opt/gcp${sysconfdir}
}

RDEPENDS_${PN} += "bash python gcp-config"
