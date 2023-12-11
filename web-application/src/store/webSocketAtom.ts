import { atom } from "jotai";
import { Message } from "@/types/playerTypes";

export const isConnectedAtom = atom<boolean>(false);
// export const messageLogAtom = atom<Message[]>([]);
export const socketAtom = atom<WebSocket | null>(null);
