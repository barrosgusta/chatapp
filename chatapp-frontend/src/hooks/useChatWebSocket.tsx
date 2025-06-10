// src/hooks/useChatWebSocket.ts
import { useEffect, useRef } from "react";
import { useChatStore } from "../store/useChatStore";

type Outgoing =
  | { type: "SET_NAME"; name: string }
  | { type: "MESSAGE"; text: string }
  | { type: "TYPING_START" }
  | { type: "TYPING_STOP" };

export function useChatWebSocket(wsUrl: string) {
  const wsRef = useRef<WebSocket | null>(null);

  // Zustand actions/state
  const {
    setUser,
    setUsers,
    addMessage,
    setTypingUsers,
    setError,
    // reset,
    user,
  } = useChatStore();

  useEffect(() => {
    // Only connect if user exists
    if (!user) {
      // Optional: cleanup any previous ws on logout
      wsRef.current?.close();
      return;
    }

    console.log("Connecting WebSocket for user:", user.name);
    // Create new WebSocket connection
    wsRef.current = new WebSocket(wsUrl);
    console.log("WebSocket created:", wsRef.current);
    wsRef.current.onopen = () => {
      console.log("WebSocket opened, sending SET_NAME:", user.name);
      wsRef.current?.send(
        JSON.stringify({ type: "SET_NAME", name: user.name })
      );
    };

    wsRef.current.onmessage = (event) => {
      console.log("WebSocket message received:", event.data);
      let data;
      try {
        data = JSON.parse(event.data);
      } catch (e) {
        console.error("Malformed WebSocket message:", event.data);
        return;
      }

      switch (data.type) {
        case "NAME_ACCEPTED":
          setError(null);
          break;
        case "NAME_REJECTED":
          setUser(null); // Return to name input
          setError(data.reason || "Name rejected. Try another.");
          break;
        case "USER_LIST":
          setUsers((data.users || []).map((name: string) => ({ name })));
          break;
        case "MESSAGE":
          if (data.message) addMessage(data.message);
          break;
        case "TYPING":
          console.log("Typing event received:", data);
          setTypingUsers(data.typing || []);
          break;
        default:
          console.warn("Unknown WS message type:", data);
          break;
      }
    };

    wsRef.current.onerror = (err) => {
      console.error("WebSocket error:", err);
      setError("WebSocket error. Check your connection.");
      // reset();
    };

    wsRef.current.onclose = () => {
      console.warn("WebSocket closed");
      setTypingUsers([]);
      setUsers([]);
      // setError("Disconnected from server.");
      // reset();
    };

    // Cleanup on unmount or user logout/change
    return () => {
      wsRef.current?.close();
    };

    // Only rerun effect when username or URL changes
    // eslint-disable-next-line
  }, [user?.name, wsUrl]);

  // Send message over WebSocket if ready
  const send = (msg: Outgoing) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
    } else {
      console.warn(
        "Tried to send WebSocket message before connection was open",
        msg
      );
    }
  };

  return { send };
}
