import { atomWithReset } from "jotai/utils";
import { MAX_PLAYER } from "./gameStore";
import { playerListAtom, userAtom } from "./chatAtom";
import { atom } from "jotai";

export const maxPlayerAtom = atomWithReset<number>(MAX_PLAYER);

export const isGameStartedAtom = atom<boolean>((get) => {
  const playerList = get(playerListAtom);
  return playerList.length === MAX_PLAYER;
});

export const isUserJoinGameAtom = atom<boolean>((get) => {
  const isGameStarted = get(isGameStartedAtom);
  const user = get(userAtom);
  return isGameStarted && user.player_type === "player";
});
