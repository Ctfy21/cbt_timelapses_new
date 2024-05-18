
// "undefined" means the URL will be computed from the `window.location` object
const URL =
    process.env.NODE_ENV === "production" ? "ws://192.168.42.119:5000/ws" : "ws://192.168.42.119:5000/ws";

export const socket = new WebSocket(URL);