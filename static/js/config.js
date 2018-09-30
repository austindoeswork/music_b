let urlparts = window.location.href.split('/')[2].split(':')

const URL = urlparts[0];
const PORT = urlparts[1] || null;
const FULL_URL = PORT ? URL + ':' + PORT : URL
const PROTOCOL = window.location.protocol.split(':')[0];
