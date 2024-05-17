
// "undefined" means the URL will be computed from the `window.location` object
const URL =
    process.env.NODE_ENV === "production" ? "ws://localhost:5000/ws" : "ws://localhost:5000/ws";

export const socket = new WebSocket(URL);