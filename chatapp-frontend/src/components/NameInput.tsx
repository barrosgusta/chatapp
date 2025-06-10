import { useState } from "react";
import { useChatStore } from "../store/useChatStore";
import { Input } from "./ui/input";
import { Button } from "./ui/button";

export default function NameInput() {
  const [name, setName] = useState("");
  const setUser = useChatStore((s) => s.setUser);
  const error = useChatStore((s) => s.error);
  const setError = useChatStore((s) => s.setError);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null); // clear any old errors
    if (name.trim()) {
      setUser({ name: name.trim() }); // Triggers ws connect and SET_NAME message
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-2">
      <label className="font-bold">Enter your display name</label>
      <Input
        value={name}
        onChange={(e) => setName(e.target.value)}
        maxLength={20}
        required
      />
      <Button type="submit" variant="outline">
        Join Chat
      </Button>
      {error && <div className="text-red-600 text-sm mt-1">{error}</div>}
    </form>
  );
}
