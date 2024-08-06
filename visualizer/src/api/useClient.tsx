import { useEffect, useState } from "react";
import { XdbClient } from "./xdb-client";

export const useClient = () => {
    const [client, setClient] = useState<XdbClient | undefined>(undefined);
    const [isConnected, setIsConnected] = useState(false);

    useEffect(() => {
        const xdbClient = XdbClient.getInstance("localhost:6789");
        setClient(xdbClient);

        xdbClient.connect().then(() => setIsConnected(true)).catch(console.error);

        return () => {
            xdbClient.close();
        };
    }, [])

    return {
        serverAddress: client?.serverAddress,
        sendMessage: client?.sendMessage,
        isConnected,
    };
}