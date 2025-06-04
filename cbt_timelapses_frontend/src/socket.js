// "undefined" means the URL will be computed from the `window.location` object

const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
const URL = protocol + '//' + window.location.hostname + ":" + window.location.port +'/ws';

export const socket = new WebSocket(URL);
