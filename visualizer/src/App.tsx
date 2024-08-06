import { useClient } from "./api/useClient";

function App() {
 const {isConnected, sendMessage, serverAddress} = useClient();

  return (
    <div>
      <p>Server Address: {serverAddress}</p>
      <p>Connection Status: {isConnected ? "Connected" : "Disconnected"}</p>
      {isConnected && (
        <button onClick={() => sendMessage?.("test", {}, "read")}>
          Send Test Message
        </button>
      )}
    </div>
  );
}

export default App;