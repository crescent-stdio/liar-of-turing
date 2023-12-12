import { User } from "@/types/playerTypes";
import crypto from "crypto";

const setIsLiars = (players: User[], liars: number[]) => {
  const newPlayers = [...players];
  liars.forEach((liar) => {
    newPlayers[liar].role = "liar";
  });
  return newPlayers;
};

export const getUserUUID = (): string => {
  if (typeof window !== "undefined") {
    const item = localStorage.getItem("key");
    const userUUID = localStorage.getItem("turing_uuid");

    if (!userUUID || userUUID.length === 0) {
      const item = crypto.randomBytes(16).toString("hex");

      localStorage.setItem("turing_uuid", item);
      return item;
    }
    return userUUID;
  }
  return "";
};
