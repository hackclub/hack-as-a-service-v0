const WebSocket = require("ws");

const token = "d4dff855f2f7481bde747dd736458924";

// const url = "http://localhost:5000/api/apps/1/stats";
const url = "http://localhost:5000/api/teams/3/expenses";

const ws = new WebSocket(url, {
  headers: {
    Cookie: `token=${token}`,
  },
});

ws.on("open", function () {
  console.log("Connected to WebSocket");
});

ws.on("message", function (data) {
  console.log(`Data: ${data}`);
});
