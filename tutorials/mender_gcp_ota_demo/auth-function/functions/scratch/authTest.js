var {google} = require('googleapis')

let cb = function(err, cred, project) {
  console.log(err);
  console.log(project);
  console.log(cred);
}

google.auth.getApplicationDefault(cb);

setTimeout(function() {console.log('done')}, 5000);