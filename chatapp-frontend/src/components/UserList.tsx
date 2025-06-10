// src/components/UserList.tsx
import { useChatStore } from "../store/useChatStore";
export default function UserList() {
  const users = useChatStore((s) => s.users);
  const current = useChatStore((s) => s.user)?.name;
  return (
    <div className="border-b pb-2 mb-2 flex flex-wrap gap-2">
      <span className="font-bold mr-2">Online users:</span>
      {users.map((u) => (
        <span
          key={u.name}
          className={u.name === current ? "font-bold text-blue-600" : ""}
        >
          {u.name}
        </span>
      ))}
    </div>
  );
}
