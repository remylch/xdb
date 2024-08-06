export class XdbClient {
  private static instance: XdbClient | null = null;
  private connectionPromise: Promise<void> | null = null;
  private ws: WebSocket | null = null;

  static getInstance(serverAddress: string): XdbClient {
    if (!XdbClient.instance) {
      XdbClient.instance = new XdbClient(serverAddress);
    }
    return XdbClient.instance;
  }

  constructor(readonly serverAddress: string) {}

  connect() {
    if (!this.connectionPromise) {
      return new Promise((resolve, reject) => {
        this.ws = new WebSocket(`ws://${this.serverAddress}`);

        this.ws.onopen = () => {
          console.log("Connected to server");
          resolve("resolved");
        };

        this.ws.onmessage = (data: unknown) => {
          console.log("Received:", data);
        };

        this.ws.onerror = (error: Event) => {
          console.error("WebSocket error:", error);
          reject(error);
        };

        this.ws.onclose = () => {
          console.log("Disconnected from server");
        };
      });
    }
    return this.connectionPromise;
  }

  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN;
  }

  sendMessage(collection: string, data: unknown, operation: "read" | "write") {
    if (!this.isConnected()) {
      throw new Error("WebSocket is not connected");
    }
    const message = JSON.stringify({ collection, data, operation });
    this.ws!.send(message);
  }

  close() {
    if (this.ws) {
      this.ws.close();
    }
    this.connectionPromise = null;
  }
}
