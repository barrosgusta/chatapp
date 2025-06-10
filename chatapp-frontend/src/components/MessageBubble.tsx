// src/components/MessageBubble.tsx
import { clsx } from "clsx";
import { useChatStore } from "../store/useChatStore";

type Props = {
  user: string;
  text: string;
  timestamp: string;
};

export default function MessageBubble({ user, text, timestamp }: Props) {
  const currentUser = useChatStore((s) => s.user)?.name;
  const isMine = currentUser === user;

  return (
    <div
      className={clsx("flex my-3", isMine ? "justify-end" : "justify-start")}
    >
      <div
        className={clsx(
          "max-w-[75%] rounded-xl px-3 py-2 shadow-lg text-sm flex flex-col",
          isMine ? "bg-blue-500 text-white" : "bg-gray-200 text-gray-900"
        )}
      >
        <span className="font-bold">{user}</span>
        <span className="mb-2 max-w-[75%]">{text}</span>
        <span
          className={clsx(
            "text-xs",
            isMine ? "text-blue-200" : "text-gray-400"
          )}
        >
          {new Date(timestamp).toLocaleString()}
        </span>
      </div>
    </div>
  );
}
