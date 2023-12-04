import { Player } from "@/types/playerTypes";

const setIsLiars = (players: Player[], liars: number[]) => {
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
      const item = crypto.randomUUID();
      localStorage.setItem("turing_uuid", item);
      return item;
    }
    return userUUID;
  }
  return "";
};
