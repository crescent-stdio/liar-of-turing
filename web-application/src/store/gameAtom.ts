import { atomWithReset } from "jotai/utils";
import { MAX_PLAYER } from "./gameStore";
import { playerListAtom, userAtom } from "./chatAtom";
import { atom, useAtomValue } from "jotai";

export const maxPlayerAtom = atom<number>(MAX_PLAYER);

export const isGameStartedAtom = atom<boolean>((get) => {
  const playerList = get(playerListAtom);
  const maxPlayer = get(maxPlayerAtom);
  return playerList.length === maxPlayer;
});

export const isUserJoinGameAtom = atom<boolean>((get) => {
  const isGameStarted = get(isGameStartedAtom);
  const user = get(userAtom);
  return isGameStarted && user.player_type === "player";
});
