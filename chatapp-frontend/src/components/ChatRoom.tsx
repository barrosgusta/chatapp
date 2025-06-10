// src/components/ChatRoom.tsx
import { useEffect, useRef } from "react";
import UserList from "./UserList";
import MessageBubble from "./MessageBubble";
import MessageInput from "./MessageInput";
import TypingIndicator from "./TypingIndicator";
import { useChatStore } from "../store/useChatStore";
import { useChatWebSocket } from "../hooks/useChatWebSocket";
import { ScrollArea } from "./ui/scroll-area";

// Adjust to match your Vite/WebSocket backend config:
const WS_URL = import.meta.env.VITE_WS_URL || "ws://localhost:8080/ws";

export default function ChatRoom() {
  const messages = useChatStore((s) => s.messages);
  const user = useChatStore((s) => s.user);

  // Set up the WebSocket connection as soon as we have a user
  console.log("ChatRoom user TESTE:", user);
  const { send } = useChatWebSocket(WS_URL);

  // For typing event debouncing:
  const typingTimeout = useRef<NodeJS.Timeout | null>(null);

  // When a message is sent:
  function handleSend(msg: string) {
    send({ type: "MESSAGE", text: msg });
    send({ type: "TYPING_STOP" });
  }

  // When typing starts:
  function handleTypingStart() {
    send({ type: "TYPING_START" });
    if (typingTimeout.current) clearTimeout(typingTimeout.current);
    typingTimeout.current = setTimeout(() => {
      send({ type: "TYPING_STOP" });
    }, 2000); // Stops typing after 2 seconds idle
  }

  // When typing stops manually (input is cleared or sent):
  function handleTypingStop() {
    send({ type: "TYPING_STOP" });
    if (typingTimeout.current) clearTimeout(typingTimeout.current);
  }

  // Scroll to bottom on new message (optional polish)
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  return (
    <div className="flex flex-col h-[32rem]">
      <div hidden={!user} className="mb-2">
        <UserList />
      </div>
      <ScrollArea className="flex-1 overflow-y-auto mb-2 flex flex-col-reverse border rounded-lg">
        <div className="m-3">
          {messages
            .slice()
            .sort(
              (a, b) =>
                new Date(a.timestamp).getTime() -
                new Date(b.timestamp).getTime()
            )
            .map((msg) => (
              <MessageBubble key={msg.id} {...msg} />
            ))}
          <div ref={bottomRef} />
        </div>
      </ScrollArea>
      <TypingIndicator />
      <MessageInput
        onSend={handleSend}
        onTypingStart={handleTypingStart}
        onTypingStop={handleTypingStop}
        disabled={!user} // Disable input if no user
        placeholder={
          user ? "Type a message..." : "Please enter your name first"
        }
      />
    </div>
  );
}
