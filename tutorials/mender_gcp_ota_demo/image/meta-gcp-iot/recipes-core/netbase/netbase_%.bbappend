inherit deploy

do_deploy() {
   install -d ${DEPLOYDIR}/persist
   mv ${D}${sysconfdir}/hosts ${DEPLOYDIR}/persist
   ln -s /data/hosts ${D}${sysconfdir}/hosts
}
addtask do_deploy after do_install before do_package
