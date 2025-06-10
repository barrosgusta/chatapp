// src/components/MessageInput.tsx
import { useState } from "react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";

type Props = {
  onSend: (msg: string) => void;
  onTypingStart?: () => void;
  onTypingStop?: () => void;
  disabled?: boolean;
  placeholder?: string;
};

export default function MessageInput({
  onSend,
  onTypingStart,
  onTypingStop,
  disabled,
  placeholder,
}: Props) {
  const [msg, setMsg] = useState("");

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setMsg(e.target.value);
    if (e.target.value && onTypingStart) onTypingStart();
    if (!e.target.value && onTypingStop) onTypingStop();
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (msg.trim()) {
      onSend(msg.trim());
      setMsg("");
      if (onTypingStop) onTypingStop();
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex gap-2 mt-2">
      <Input
        // className="flex-1 border rounded px-2 py-1"
        value={msg}
        onChange={handleChange}
        placeholder={placeholder || "Type a message…"}
        maxLength={256}
        autoFocus
        disabled={disabled}
      />
      <Button type="submit" disabled={disabled || !msg.trim()}>
        Send
      </Button>
    </form>
  );
}
