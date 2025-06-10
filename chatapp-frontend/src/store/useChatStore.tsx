import { create } from "zustand";

type User = { name: string };
type Message = { id: string; user: string; text: string; timestamp: string };

interface ChatState {
  user: User | null;
  users: User[];
  messages: Message[];
  typingUsers: string[];
  error: string | null;
  setUser: (user: User | null) => void;
  setUsers: (users: User[]) => void;
  addMessage: (msg: Message) => void;
  setTypingUsers: (names: string[]) => void;
  setError: (err: string | null) => void;
  reset: () => void;
}

export const useChatStore = create<ChatState>((set) => ({
  user: null,
  users: [],
  messages: [],
  typingUsers: [],
  error: null,
  setUser: (user) => set({ user }),
  setUsers: (users) => set({ users }),
  addMessage: (msg) =>
    set((state) => {
      // Prevent duplicate messages by id
      if (state.messages.some((m) => m.id === msg.id)) {
        return state;
      }
      return { messages: [...state.messages, msg] };
    }),
  setTypingUsers: (names) => set({ typingUsers: names }),
  setError: (err) => set({ error: err }),
  reset: () =>
    set({
      user: null,
      users: [],
      messages: [],
      typingUsers: [],
      error: null,
    }),
}));
