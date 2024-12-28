const WebSocket = require('ws'); // Import the WebSocket library

const ws = new WebSocket("ws://localhost:3000/ws");

ws.on('open', () => {
    console.log("WebSocket connection opened");
});

ws.on('message', (data) => {
    console.log("Message from server:", data.toString());
});

ws.on('close', () => {
    console.log("WebSocket connection closed");
});

ws.on('error', (error) => {
    console.error("WebSocket error:", error);
});