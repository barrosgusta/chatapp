import { useChatStore } from "./store/useChatStore";
import NameInput from "./components/NameInput";
import ChatRoom from "./components/ChatRoom";
import { useEffect } from "react";
import axios from "axios";
import type { Message } from "./types";
import { Card } from "./components/ui/card";
import { ThemeProvider } from "./components/ThemeProvider";

export default function App() {
  const user = useChatStore((s) => s.user);

  const messages = useChatStore((s) => s.messages);
  console.log("Messages in store:", messages);
  useEffect(() => {
    if (!messages || messages.length === 0) {
      console.log("Fetching messages from server...");
      axios
        .get(`${import.meta.env.VITE_CHAT_SERVICE_URL}/history`)
        .then((response) => {
          console.log("Messages fetched from server:", response.data);
          if (Array.isArray(response.data)) {
            const messages = response.data.map((msg) => ({
              id: msg.id,
              user: msg.user,
              text: msg.text,
              timestamp: msg.timestamp,
            }));
            console.log("Fetched messages:", messages);
            messages.forEach((msg: Message) => {
              useChatStore.getState().addMessage(msg);
            });
            console.log("Messages loaded:", messages);
          }
        })
        .catch((error) => {
          console.error("Failed to fetch messages:", error);
          useChatStore.getState().setError("Failed to load messages.");
        });
    }
  }, []);

  console.log("Current user:", user);
  return (
    <ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
      <div className="min-h-screen min-w-screen px-3 py-4 xl:px-80 xl:py-10 flex items-center justify-center">
        <Card className="w-full p-4">
          {!user && <NameInput />}
          <ChatRoom />
        </Card>
      </div>
    </ThemeProvider>
  );
}
