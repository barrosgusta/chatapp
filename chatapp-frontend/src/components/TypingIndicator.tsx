// src/components/TypingIndicator.tsx
import { useChatStore } from "../store/useChatStore";
export default function TypingIndicator() {
  const typingUsers = useChatStore((s) => s.typingUsers);
  const current = useChatStore((s) => s.user)?.name;
  const othersTyping = typingUsers.filter((name) => name !== current);
  if (!othersTyping.length) return null;
  return (
    <div className="text-xs text-gray-500 mt-1 mb-2">
      {othersTyping.join(", ")} {othersTyping.length === 1 ? "is" : "are"}{" "}
      typing…
    </div>
  );
}
