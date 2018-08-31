SUMMARY = "Google Cloud Platform Configuration"
LICENSE = "Apache-2.0"
LIC_FILES_CHKSUM = "file://LICENSE;md5=0ea4e253cc22ddc22117b9796e5ce5b7"

FILESEXTRAPATHS_prepend := "${THISDIR}/files:"
SRC_URI = "file://gcp-config.sh file://LICENSE"

FILES_${PN} = "/opt/gcp${sysconfdir}/gcp-config.sh"

S = "${WORKDIR}"

do_install() {
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
    
    install -m 0700 -d ${D}/opt/gcp${sysconfdir}
    install -m 0700 ${WORKDIR}/gcp-config.sh ${D}/opt/gcp${sysconfdir}

    sed -i -e 's,@PROJECT_ID@,${PROJECT_ID},g' \
           -e 's,@REGION_ID@,${REGION_ID},g' \
           -e 's,@REGISTRY_ID@,${REGISTRY_ID},g' \
           ${D}/opt/gcp${sysconfdir}/gcp-config.sh
}

RDEPENDS_${PN} += "bash"
